package mongo

import (
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

func (r *repository) CreateUser(user *user.User) error {
	ctx, cancel := r.newCtx()
	defer cancel()

	res, err := r.usersColl.InsertOne(
		ctx,
		bson.M{
			"name":       bson.M{"first": user.Name.First, "last": user.Name.Last},
			"username":   user.Username,
			"email":      user.Email,
			"password":   user.Password,
			"phone":      user.Phone,
			"role":       user.Role,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	)

	if err == nil {
		user.ID = res.InsertedID.(primitive.ObjectID).Hex()
	}

	return err
}
