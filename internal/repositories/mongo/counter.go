package mongo

import (
	"errors"
	"fmt"

	"github.com/omareloui/odinls/internal/application/core/counter"
	"github.com/omareloui/odinls/internal/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const amountToIncrement = 1

func (r *repository) CreateCounter(cntr *counter.Counter) error {
	ctx, cancel := r.newCtx()
	defer cancel()

	doc, err := r.bu.MarshalBsonD(cntr, r.bu.WithObjectID("merchant"))
	if err != nil {
		return err
	}

	res, err := r.countersColl.InsertOne(ctx, doc)

	if err == nil {
		cntr.ID = res.InsertedID.(primitive.ObjectID).Hex()
	}

	if ok := mongo.IsDuplicateKeyError(err); ok {
		if se := mongo.ServerError(nil); errors.As(err, &se) {
			if se.HasErrorMessage(" merchant: ") {
				return counter.ErrAlreadyExistingCounter
			}
		}
	}

	return err
}

func (r *repository) GetCounterByID(id string) (*counter.Counter, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	filter := bson.M{"_id": objId}

	cntr := &counter.Counter{}

	err = r.countersColl.FindOne(ctx, filter).Decode(cntr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, counter.ErrCounterNotFound
		}
		return nil, err
	}

	return cntr, nil
}

func (r *repository) GetCounterByMerchantID(merchantId string) (*counter.Counter, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	mrId, err := primitive.ObjectIDFromHex(merchantId)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	filter := bson.M{"merchant": mrId}

	cntr := &counter.Counter{}

	err = r.countersColl.FindOne(ctx, filter).Decode(cntr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, counter.ErrCounterNotFound
		}
		return nil, err
	}

	return cntr, nil
}

func (r *repository) AddOneToProduct(merchantId, category string) (uint8, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	merId, err := primitive.ObjectIDFromHex(merchantId)
	if err != nil {
		return 0, errs.ErrInvalidID
	}

	cntr, err := r.GetCounterByMerchantID(merchantId)
	if err != nil {
		return 0, err
	}

	filter := bson.M{"merchant": merId}
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

func (r *repository) AddOneToOrder(merchantId string) (uint, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	merId, err := primitive.ObjectIDFromHex(merchantId)
	if err != nil {
		return 0, errs.ErrInvalidID
	}

	cntr, err := r.GetCounterByMerchantID(merchantId)
	if err != nil {
		return 0, err
	}

	filter := bson.D{{Key: "merchant", Value: merId}}
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
