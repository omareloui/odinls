package mongo

import (
	"errors"
	"fmt"
	"time"

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

	now := time.Now()

	if cntr.MerchantID == "" {
		return errs.ErrInvalidID
	}

	merId, err := primitive.ObjectIDFromHex(cntr.MerchantID)
	if err != nil {
		return errs.ErrInvalidID
	}

	document := bson.M{
		"merchant":       merId,
		"orders_number":  cntr.OrdersNumber,
		"products_codes": cntr.ProductsCodes,
		"created_at":     now,
		"updated_at":     now,
	}

	res, err := r.countersColl.InsertOne(ctx, document)

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

	filter := bson.D{{Key: "merchant", Value: merId}}
	update := bson.M{
		"$inc": bson.M{fmt.Sprintf("products_codes.%s", category): amountToIncrement},
	}

	updated, err := r.usersColl.UpdateOne(ctx, filter, update)
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

	updated, err := r.usersColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	if updated.ModifiedCount == 0 {
		return 0, counter.ErrCounterNotFound
	}

	return cntr.OrdersNumber + amountToIncrement, nil
}

// func (r *repository) GetCounterByID(id string, options ...counter.RetrieveOptsFunc) (*counter.Counter, error) {
// 	opts := counter.ParseRetrieveOpts(options...)

// 	ctx, cancel := r.newCtx()
// 	defer cancel()

// 	objId, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, errs.ErrInvalidID
// 	}

// 	filter := bson.M{"_id": objId}

// 	cli := &counter.Counter{}

// 	err = r.countersColl.FindOne(ctx, filter).Decode(cli)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, counter.ErrCounterNotFound
// 		}
// 		return nil, err
// 	}

// 	if opts.PopulateMerchant {
// 		r.populateMerchantForCounter(cli)
// 	}

// 	return cli, nil
// }

// func (r *repository) CreateCounter(cli *counter.Counter, options ...counter.RetrieveOptsFunc) error {
// 	opts := counter.ParseRetrieveOpts(options...)

// 	mrId, err := primitive.ObjectIDFromHex(cli.MerchantID)
// 	if err != nil {
// 		return errs.ErrInvalidID
// 	}

// 	ctx, cancel := r.newCtx()
// 	defer cancel()

// 	res, err := r.countersColl.InsertOne(ctx, bson.M{
// 		"merchant":             mrId,
// 		"name":                 cli.Name,
// 		"notes":                cli.Notes,
// 		"contact_info":         cli.ContactInfo,
// 		"wholesale_as_default": cli.WholesaleAsDefault,
// 		"created_at":           cli.CreatedAt,
// 		"updated_at":           cli.UpdatedAt,
// 	})

// 	if err == nil {
// 		cli.ID = res.InsertedID.(primitive.ObjectID).Hex()
// 		if opts.PopulateMerchant {
// 			r.populateMerchantForCounter(cli)
// 		}
// 	}

// 	if ok := mongo.IsDuplicateKeyError(err); ok {
// 		if se := mongo.ServerError(nil); errors.As(err, &se) {
// 			if se.HasErrorMessage(" name: ") && se.HasErrorMessage(" merchant: ") {
// 				return counter.ErrCounterExistsForMerchant
// 			}
// 		}
// 	}

// 	return err
// }

// func (r *repository) UpdateCounterByID(id string, cli *counter.Counter, options ...counter.RetrieveOptsFunc) error {
// 	opts := counter.ParseRetrieveOpts(options...)

// 	ctx, cancel := r.newCtx()
// 	defer cancel()

// 	fmt.Println(id)
// 	objId, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return errs.ErrInvalidID
// 	}

// 	filter := bson.M{"_id": objId}

// 	res := r.countersColl.FindOneAndUpdate(ctx, filter, bson.M{
// 		"$set": bson.M{
// 			"name":                 cli.Name,
// 			"notes":                cli.Notes,
// 			"contact_info":         cli.ContactInfo,
// 			"wholesale_as_default": cli.WholesaleAsDefault,
// 			"created_at":           cli.CreatedAt,
// 			"updated_at":           cli.UpdatedAt,
// 		},
// 	})

// 	err = res.Err()
// 	if err == nil {
// 		cli.ID = id
// 		if opts.PopulateMerchant {
// 			r.populateMerchantForCounter(cli)
// 		}
// 	}

// 	if ok := mongo.IsDuplicateKeyError(err); ok {
// 		if se := mongo.ServerError(nil); errors.As(err, &se) {
// 			if se.HasErrorMessage(" name: ") && se.HasErrorMessage(" merchant: ") {
// 				return counter.ErrCounterExistsForMerchant
// 			}
// 		}
// 	}

// 	return err
// }

// func (r *repository) populateMerchantForCounter(cli *counter.Counter) {
// 	merchant, err := r.FindMerchant(cli.MerchantID)
// 	if err == nil {
// 		cli.Merchant = merchant
// 	}
// }
