package mongo

import (
	"fmt"

	"github.com/omareloui/odinls/internal/application/core/counter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const amountToIncrement = 1

func (r *repository) CreateCounter(cntr *counter.Counter) error {
	ctx, cancel := r.newCtx()
	defer cancel()

	doc, err := r.bu.MarshalBsonD(cntr)
	if err != nil {
		return err
	}

	res, err := r.countersColl.InsertOne(ctx, doc)

	if err == nil {
		cntr.ID = res.InsertedID.(primitive.ObjectID).Hex()
	}

	if ok := mongo.IsDuplicateKeyError(err); ok {
		return counter.ErrAlreadyExistingCounter
	}

	return err
}

func (r *repository) GetCounter() (*counter.Counter, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	cntr := &counter.Counter{}

	err := r.countersColl.FindOne(ctx, bson.M{}).Decode(cntr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, counter.ErrCounterNotFound
		}
		return nil, err
	}

	return cntr, nil
}

func (r *repository) AddOneToProduct(category string) (uint8, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	cntr, err := r.GetCounter()
	if err != nil {
		return 0, err
	}

	filter := bson.M{}
	update := bson.M{
		"$inc": bson.M{fmt.Sprintf("products_codes.%s", category): amountToIncrement},
	}

	updated, err := r.countersColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	if updated.ModifiedCount == 0 {
		return 0, counter.ErrCounterNotFound
	}

	return cntr.ProductsCodes[category] + amountToIncrement, nil
}

func (r *repository) AddOneToOrder() (uint, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	cntr, err := r.GetCounter()
	if err != nil {
		return 0, err
	}

	filter := bson.M{}
	update := bson.M{
		"$inc": bson.M{"orders_number": amountToIncrement},
	}

	updated, err := r.countersColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	if updated.ModifiedCount == 0 {
		return 0, counter.ErrCounterNotFound
	}

	return cntr.OrdersNumber + amountToIncrement, nil
}
