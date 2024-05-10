package db

import (
	"context"
	"fmt"
	"time"

	"github.com/omareloui/odinls/internal/application/core/domain"
	"github.com/omareloui/odinls/internal/misc/app_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	BCryptCost int = 10
)

type Adapter struct {
	client   *mongo.Client
	db       *mongo.Database
	usersCol *mongo.Collection
}

func NewAdapter(uri string) (*Adapter, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(fmt.Errorf("error trying to connect to the database \"%s\": %s", uri, err))
	}
	db := client.Database("ODINLS_DEV")
	usersCol := db.Collection("users")
	return &Adapter{client: client, db: db, usersCol: usersCol}, nil
}

func (a *Adapter) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	filter := bson.D{{Key: "email", Value: email}}
	result := &domain.User{}
	if err := a.usersCol.FindOne(ctx, filter).Decode(result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, app_errors.NewEntityNotFound("user", fmt.Sprintf("email: %s", email))
		}
		return nil, err
	}
	return result, nil
}

func (a *Adapter) CreateUser(ctx context.Context, dto domain.Register) (*domain.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(dto.Password), BCryptCost)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	usr := domain.User{
		Name:      dto.Name,
		Email:     dto.Email,
		Password:  string(hashed),
		CreatedAt: now,
		UpdatedAt: now,
	}
	result, err := a.usersCol.InsertOne(ctx, usr)
	if err != nil {
		return nil, err
	}
	usr.ID = result.InsertedID.(primitive.ObjectID)
	return &usr, nil
}
