package mongo

import (
	"context"
	"log"
	"time"

	"github.com/omareloui/odinls/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	r "github.com/omareloui/odinls/internal/repositories"
)

const (
	merchantsCollectionName = "merchants"
	usersCollectionName     = "users"
	rolesCollectionName     = "roles"
	clientsCollectionName   = "clients"
	countersCollectionName  = "counters"
	productsCollectionName  = "products"
	ordersCollectionName    = "orders"
)

type repository struct {
	client  *mongo.Client
	timeout time.Duration
	db      *mongo.Database

	merchantsColl *mongo.Collection
	usersColl     *mongo.Collection
	rolesColl     *mongo.Collection
	clientsColl   *mongo.Collection
	countersColl  *mongo.Collection
	productsColl  *mongo.Collection
	ordersColl    *mongo.Collection
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

func createIndex(coll *mongo.Collection, idxModel mongo.IndexModel) {
	_, err := coll.Indexes().CreateOne(context.Background(), idxModel)
	if err != nil {
		log.Fatalf("error creating the for the %s model\n", coll.Name())
	}
}

func NewRepository(mongoURL, dbName string, mongoTimeout int) (r.Repository, error) {
	repo := &repository{timeout: time.Duration(mongoTimeout) * time.Second}
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, err
	}
	repo.client = client
	repo.db = client.Database(dbName)

	repo.usersColl = repo.db.Collection(usersCollectionName)
	createIndex(repo.usersColl, mongo.IndexModel{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)})
	createIndex(repo.usersColl, mongo.IndexModel{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)})

	repo.merchantsColl = repo.db.Collection(merchantsCollectionName)

	repo.rolesColl = repo.db.Collection(rolesCollectionName)
	createIndex(repo.rolesColl, mongo.IndexModel{Keys: bson.D{{Key: "name", Value: 1}}, Options: options.Index().SetUnique(true)})

	repo.clientsColl = repo.db.Collection(clientsCollectionName)
	createIndex(repo.clientsColl, mongo.IndexModel{Keys: bson.D{{Key: "name", Value: 1}, {Key: "merchant", Value: 1}}, Options: options.Index().SetUnique(true)})

	repo.countersColl = repo.db.Collection(countersCollectionName)
	createIndex(repo.countersColl, mongo.IndexModel{Keys: bson.D{{Key: "merchant", Value: 1}}, Options: options.Index().SetUnique(true)})

	repo.productsColl = repo.db.Collection(productsCollectionName)
	createIndex(repo.productsColl, mongo.IndexModel{Keys: bson.D{{Key: "merchant", Value: 1}}})
	createIndex(repo.productsColl, mongo.IndexModel{Keys: bson.D{{Key: "variants._id", Value: 1}}, Options: options.Index().SetUnique(true)})

	repo.ordersColl = repo.db.Collection(ordersCollectionName)
	createIndex(repo.ordersColl, mongo.IndexModel{Keys: bson.D{{Key: "ref", Value: 1}}, Options: options.Index().SetUnique(true)})
	createIndex(repo.ordersColl, mongo.IndexModel{Keys: bson.D{{Key: "merchant", Value: 1}}})
	createIndex(repo.ordersColl, mongo.IndexModel{Keys: bson.D{{Key: "client", Value: 1}}})
	createIndex(repo.ordersColl, mongo.IndexModel{Keys: bson.D{{Key: "items._id", Value: 1}}, Options: options.Index().SetUnique(true)})

	return repo, nil
}
