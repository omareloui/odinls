package mongo

import (
	"context"
	"time"

	"github.com/omareloui/odinls/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	r "github.com/omareloui/odinls/internal/repositories"
)

const (
	merchantsCollectionName = "merchants"
	usersCollectionName     = "users"
)

type repository struct {
	client        *mongo.Client
	timeout       time.Duration
	db            *mongo.Database
	merchantsColl *mongo.Collection
	usersColl     *mongo.Collection
}

func (r *repository) newCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), r.timeout)
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL).SetAuth(config.GetMongoCred()))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, err
}

func NewRepository(mongoURL, dbName string, mongoTimeout int) (r.Repository, error) {
	repo := &repository{timeout: time.Duration(mongoTimeout) * time.Second}
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, err
	}
	repo.client = client
	repo.db = client.Database(dbName)
	repo.merchantsColl = repo.db.Collection(merchantsCollectionName)
	repo.usersColl = repo.db.Collection(usersCollectionName)
	return repo, nil
}
