package mongo

import (
	"time"

	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/repositories/mongo/bsonutils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) GetUsers() ([]user.User, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	return GetAll[user.User](ctx, r.usersColl)
}

func (r *repository) GetUsersByIDs(ids []string) ([]user.User, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	return GetByIDs[user.User](ctx, r.usersColl, ids)
}

func (r *repository) GetUser(id string) (*user.User, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	return GetByID[user.User](ctx, r.usersColl, id)
}

func (r *repository) GetUserByEmailOrUsernameFromUser(usr *user.User) (*user.User, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	filter := bson.M{
		"$or": bson.A{
			bson.M{"email": usr.Email},
			bson.M{"username": usr.Username},
		},
	}

	u, err := GetOne[user.User](ctx, r.usersColl, filter)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *repository) GetUserByEmailOrUsername(emailOrUsername string) (*user.User, error) {
	u := &user.User{Email: emailOrUsername, Username: emailOrUsername}
	return r.GetUserByEmailOrUsernameFromUser(u)
}

func (r *repository) CreateUser(u *user.User) (*user.User, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	return InsertStruct(ctx, r.usersColl, u)
}

func (r *repository) UpdateUserByID(id string, u *user.User) (*user.User, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	return UpdateStructByID(ctx, r.usersColl, id, u,
		bsonutils.WithFieldToRemove("password"))
}

func (r *repository) UnsetCraftsmanByID(id string) (*user.User, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	update := bson.M{
		"$unset": bson.M{"craftsman": ""},
		"$set":   bson.M{"updated_at": time.Now()},
	}

	return UpdateByID[user.User](ctx, r.usersColl, id, update)
}
