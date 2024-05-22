package mongo

import (
	"errors"
	"time"

	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetUsers(options ...user.RetrieveOptsFunc) ([]user.User, error) {
	opts := user.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	var cursor *mongo.Cursor
	var err error

	filter := bson.M{}

	if !opts.PopulateRole {
		cursor, err = r.usersColl.Find(ctx, filter)
	} else {
		cursor, err = r.usersColl.Aggregate(ctx, bson.A{
			bson.M{
				"$lookup": bson.M{
					"from":         rolesCollectionName,
					"localField":   "role",
					"foreignField": "_id",
					"as":           "populatedRole",
				},
			},
			bson.M{"$unwind": "$populatedRole"},
		})
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	usrs := new([]user.User)
	if err := cursor.All(ctx, usrs); err != nil {
		return nil, err
	}

	return *usrs, nil
}

func (r *repository) FindUser(id string, options ...user.RetrieveOptsFunc) (*user.User, error) {
	opts := user.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	u := &user.User{}
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrInvalidID
	}
	filter := bson.M{"_id": objId}
	err = r.usersColl.FindOne(ctx, filter).Decode(u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	if opts.PopulateRole {
		r.populateRoleForUser(u)
	}

	return u, nil
}

func (r *repository) FindUserByEmailOrUsernameFromUser(usr *user.User, options ...user.RetrieveOptsFunc) (*user.User, error) {
	opts := user.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	u := &user.User{}
	filter := bson.M{
		"$or": bson.A{
			bson.M{"email": usr.Email},
			bson.M{"username": usr.Username},
		},
	}

	err := r.usersColl.FindOne(ctx, filter).Decode(u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	if opts.PopulateRole {
		r.populateRoleForUser(u)
	}

	return u, nil
}

func (r *repository) FindUserByEmailOrUsername(emailOrUsername string, options ...user.RetrieveOptsFunc) (*user.User, error) {
	opts := user.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	u := &user.User{}
	filter := bson.M{
		"$or": bson.A{
			bson.M{"email": emailOrUsername},
			bson.M{"username": emailOrUsername},
		},
	}

	err := r.usersColl.FindOne(ctx, filter).Decode(u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	if opts.PopulateRole {
		r.populateRoleForUser(u)
	}

	return u, nil
}

func (r *repository) CreateUser(usr *user.User, options ...user.RetrieveOptsFunc) error {
	opts := user.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	// TODO(security): make sure to prevent to create multiple emails with +
	// eg. "contact@omareloui.com" is the same as "contact+whatever@omareloui.com"

	res, err := r.usersColl.InsertOne(ctx, usr)

	if err == nil {
		usr.ID = res.InsertedID.(primitive.ObjectID).Hex()
	}

	if ok := mongo.IsDuplicateKeyError(err); ok {
		if se := mongo.ServerError(nil); errors.As(err, &se) {
			if se.HasErrorMessage("{ username: ") {
				return user.ErrUsernameAlreadyExists
			}
			if se.HasErrorMessage("{ email: ") {
				return user.ErrEmailAlreadyExists
			}
		}
	}

	if opts.PopulateRole {
		r.populateRoleForUser(usr)
	}

	return err
}

func (r *repository) UpdateUserByID(id string, usr *user.User, options ...user.RetrieveOptsFunc) error {
	opts := user.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errs.ErrInvalidID
	}

	roleId, err := primitive.ObjectIDFromHex(usr.RoleID)
	if err != nil {
		return errs.ErrInvalidID
	}

	filter := bson.D{{Key: "_id", Value: objId}}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.M{
				"name":       usr.Name,
				"email":      usr.Email,
				"username":   usr.Username,
				"role":       roleId,
				"updated_at": time.Now(),
			},
		},
	}

	updated, err := r.usersColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if updated.ModifiedCount == 0 {
		return user.ErrUserNotFound
	}

	if opts.PopulateRole {
		r.populateRoleForUser(usr)
	}

	return nil
}

func (r *repository) populateRoleForUser(u *user.User) {
	role, err := r.FindRole(u.RoleID)
	if err == nil {
		u.Role = role
	}
}
