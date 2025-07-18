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
	ctx, cancel := r.newCtx()
	defer cancel()

	var cursor *mongo.Cursor
	var err error

	filter := bson.M{}

	cursor, err = r.usersColl.Find(ctx, filter)
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

func (r *repository) FindUsersByIDs(ids []string) ([]user.User, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	objIds := make([]primitive.ObjectID, len(ids))

	for i, id := range ids {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errs.ErrInvalidID
		}

		objIds[i] = objId
	}

	filter := bson.M{"_id": bson.M{"$in": objIds}}

	cursor, err := r.usersColl.Find(ctx, filter)
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

	return u, nil
}

func (r *repository) FindUserByEmailOrUsernameFromUser(usr *user.User, options ...user.RetrieveOptsFunc) (*user.User, error) {
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
	return u, nil
}

func (r *repository) FindUserByEmailOrUsername(emailOrUsername string, options ...user.RetrieveOptsFunc) (*user.User, error) {
	u := &user.User{Email: emailOrUsername, Username: emailOrUsername}
	return r.FindUserByEmailOrUsernameFromUser(u, options...)
}

func (r *repository) CreateUser(u *user.User, options ...user.RetrieveOptsFunc) error {
	ctx, cancel := r.newCtx()
	defer cancel()

	// TODO(security): make sure to prevent to create multiple emails with +
	// eg. "contact@omareloui.com" is the same as "contact+whatever@omareloui.com"

	doc, err := r.bu.MarshalBsonD(u, r.bu.WithStringfied("role"))
	if err != nil {
		return err
	}

	res, err := r.usersColl.InsertOne(ctx, doc)

	if err == nil {
		u.ID = res.InsertedID.(primitive.ObjectID).Hex()
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

	return err
}

func (r *repository) UpdateUserByID(id string, u *user.User, options ...user.RetrieveOptsFunc) error {
	ctx, cancel := r.newCtx()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errs.ErrInvalidID
	}

	var doc primitive.D

	doc, err = r.bu.MarshalBsonD(u,
		r.bu.WithFieldToRemove("_id"),
		r.bu.WithStringfied("role"),
		r.bu.WithFieldToRemove("password"),
		r.bu.WithFieldToRemove("created_at"),
		r.bu.WithUpdatedAt(),
	)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objId}}
	update := bson.M{"$set": doc}

	updated, err := r.usersColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if updated.ModifiedCount == 0 {
		return user.ErrUserNotFound
	}

	return nil
}

func (r *repository) UnsetCraftsmanByID(id string) error {
	ctx, cancel := r.newCtx()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errs.ErrInvalidID
	}

	filter := bson.D{{Key: "_id", Value: objId}}
	update := bson.M{
		"$unset": bson.M{"craftsman": ""},
		"$set":   bson.M{"updated_at": time.Now()},
	}

	updated, err := r.usersColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if updated.ModifiedCount == 0 {
		return user.ErrUserNotFound
	}

	return nil
}
