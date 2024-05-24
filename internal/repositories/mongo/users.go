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

	if !opts.PopulateRole && !opts.PopulateMerchant {
		cursor, err = r.usersColl.Find(ctx, filter)
	} else {
		pipeline := []interface{}{}

		if opts.PopulateRole {
			pipeline = append(pipeline, bson.M{
				"$lookup": bson.M{
					"from":         rolesCollectionName,
					"localField":   "role",
					"foreignField": "_id",
					"as":           "populatedRole",
				},
			},
				bson.M{"$unwind": "$populatedRole"},
			)
		}
		if opts.PopulateMerchant {
			pipeline = append(pipeline, bson.M{
				"$lookup": bson.M{
					"from":         merchantsCollectionName,
					"localField":   "merchant",
					"foreignField": "_id",
					"as":           "populatedMerchant",
				},
			},
				bson.M{"$unwind": "$populatedMerchant"},
			)
		}

		cursor, err = r.usersColl.Aggregate(ctx, pipeline)
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
	if opts.PopulateMerchant {
		r.populateMerchantForUser(u)
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
	if opts.PopulateMerchant {
		r.populateMerchantForUser(u)
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
	if opts.PopulateMerchant {
		r.populateMerchantForUser(u)
	}

	return u, nil
}

func (r *repository) CreateUser(u *user.User, options ...user.RetrieveOptsFunc) error {
	opts := user.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	// TODO(security): make sure to prevent to create multiple emails with +
	// eg. "contact@omareloui.com" is the same as "contact+whatever@omareloui.com"

	res, err := r.usersColl.InsertOne(ctx, u)

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

	if opts.PopulateRole {
		r.populateRoleForUser(u)
	}
	if opts.PopulateMerchant {
		r.populateMerchantForUser(u)
	}

	return err
}

func (r *repository) UpdateUserByID(id string, u *user.User, options ...user.RetrieveOptsFunc) error {
	opts := user.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errs.ErrInvalidID
	}

	roleId, err := primitive.ObjectIDFromHex(u.RoleID)
	if err != nil {
		return errs.ErrInvalidID
	}

	filter := bson.D{{Key: "_id", Value: objId}}
	update := bson.M{
		"$set": bson.M{
			"name":       u.Name,
			"email":      u.Email,
			"username":   u.Username,
			"role":       roleId,
			"updated_at": time.Now(),
		},
	}

	if u.Craftsman != nil {
		merId, err := primitive.ObjectIDFromHex(u.Craftsman.MerchantID)
		if err != nil {
			return errs.ErrInvalidID
		}

		(update["$set"]).(bson.M)["craftsman"] = map[string]interface{}{
			"merchant":    merId,
			"hourly_rate": u.Craftsman.HourlyRate,
		}
	}

	updated, err := r.usersColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if updated.ModifiedCount == 0 {
		return user.ErrUserNotFound
	}

	if opts.PopulateRole {
		r.populateRoleForUser(u)
	}
	if opts.PopulateMerchant {
		r.populateMerchantForUser(u)
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

func (r *repository) populateMerchantForUser(u *user.User) {
	if u.Craftsman == nil || u.Craftsman.MerchantID == "" {
		return
	}

	merchant, err := r.FindMerchant(u.RoleID)
	if err == nil {
		u.Craftsman.Merchant = merchant
	}
}

func (r *repository) populateRoleForUser(u *user.User) {
	role, err := r.FindRole(u.RoleID)
	if err == nil {
		u.Role = role
	}
}
