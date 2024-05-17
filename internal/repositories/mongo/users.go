package mongo

import (
	"errors"

	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) FindUser(id string) (*user.User, error) {
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

func (r *repository) FindUserByEmailOrUsernameFromUser(usr *user.User) (*user.User, error) {
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

func (r *repository) FindUserByEmailOrUsername(emailOrUsername string) (*user.User, error) {
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

	return u, nil
}

func (r *repository) CreateUser(usr *user.User) error {
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

	return err
}
