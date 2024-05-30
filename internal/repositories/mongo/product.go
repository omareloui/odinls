package mongo

import (
	"time"

	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetProducts(options ...product.RetrieveOptsFunc) ([]product.Product, error) {
	opts := product.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	pipeline := bson.A{}

	if opts.PopulateMerchant {
		pipeline = append(pipeline, bson.M{
			"$lookup": bson.M{
				"from":         merchantsCollectionName,
				"localField":   "merchant",
				"foreignField": "_id",
				"as":           "populatedMerchant",
			},
		}, bson.M{"$unwind": "$populatedMerchant"})
	}
	if opts.PopulateCraftsman {
		pipeline = append(pipeline, bson.M{
			"$lookup": bson.M{
				"from":         usersCollectionName,
				"localField":   "craftsman",
				"foreignField": "_id",
				"as":           "populatedCraftsman",
			},
		}, bson.M{"$unwind": "$populatedCraftsman"})
	}

	cur, err := r.productsColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	prods := new([]product.Product)
	if err := cur.All(ctx, prods); err != nil {
		return nil, err
	}

	return *prods, nil
}

func (r *repository) GetProductByID(id string, options ...product.RetrieveOptsFunc) (*product.Product, error) {
	opts := product.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	filter := bson.M{"_id": objId}

	prod := &product.Product{}

	err = r.productsColl.FindOne(ctx, filter).Decode(prod)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, product.ErrProductNotFound
		}
		return nil, err
	}

	if opts.PopulateCraftsman {
		r.populateCraftsmanForProduct(prod)
	}
	if opts.PopulateMerchant {
		r.populateMerchantForProduct(prod)
	}

	return prod, nil
}

func (r *repository) GetProductsByMerchantID(merchantId string, options ...product.RetrieveOptsFunc) ([]product.Product, error) {
	opts := product.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	merId, err := primitive.ObjectIDFromHex(merchantId)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	pipeline := bson.A{bson.M{"$match": bson.M{"merchant": merId}}}

	if opts.PopulateMerchant {
		pipeline = append(pipeline, bson.M{
			"$lookup": bson.M{
				"from":         merchantsCollectionName,
				"localField":   "merchant",
				"foreignField": "_id",
				"as":           "populatedMerchant",
			},
		}, bson.M{"$unwind": "$populatedMerchant"})
	}
	if opts.PopulateCraftsman {
		pipeline = append(pipeline, bson.M{
			"$lookup": bson.M{
				"from":         usersCollectionName,
				"localField":   "craftsman",
				"foreignField": "_id",
				"as":           "populatedCraftsman",
			},
		}, bson.M{"$unwind": "$populatedCraftsman"})
	}

	cur, err := r.productsColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	prods := new([]product.Product)
	if err := cur.All(ctx, prods); err != nil {
		return nil, err
	}

	return *prods, nil
}

func (r *repository) CreateProduct(prod *product.Product, options ...product.RetrieveOptsFunc) error {
	opts := product.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	merId, err := primitive.ObjectIDFromHex(prod.MerchantID)
	if err != nil {
		return errs.ErrInvalidID
	}
	crafId, err := primitive.ObjectIDFromHex(prod.CraftsmanID)
	if err != nil {
		return errs.ErrInvalidID
	}

	now := time.Now()

	doc := bson.M{
		"merchant":  merId,
		"craftsman": crafId,
		"number":    prod.Number,
		"name":      prod.Name,
		"category":  prod.Category,

		"variants": prod.Variants,

		"created_at": now,
		"updated_at": now,
	}

	if prod.Description != "" {
		doc["description"] = prod.Description
	}

	res, err := r.productsColl.InsertOne(ctx, doc)

	if err == nil {
		prod.ID = res.InsertedID.(primitive.ObjectID).Hex()
		if opts.PopulateCraftsman {
			r.populateCraftsmanForProduct(prod)
		}
		if opts.PopulateMerchant {
			r.populateMerchantForProduct(prod)
		}
	}

	return err
}

func (r *repository) UpdateProductByID(id string, prod *product.Product, options ...product.RetrieveOptsFunc) error {
	opts := product.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errs.ErrInvalidID
	}

	var merId primitive.ObjectID
	var crafId primitive.ObjectID

	if prod.MerchantID != "" {
		merId, err = primitive.ObjectIDFromHex(prod.MerchantID)
		if err != nil {
			return errs.ErrInvalidID
		}
	}

	if prod.CraftsmanID != "" {
		crafId, err = primitive.ObjectIDFromHex(prod.CraftsmanID)
		if err != nil {
			return errs.ErrInvalidID
		}
	}

	filter := bson.M{"_id": objId}

	doc := bson.M{
		"name":       prod.Name,
		"number":     prod.Number,
		"category":   prod.Category,
		"variants":   prod.Variants,
		"updated_at": time.Now(),
	}

	if !crafId.IsZero() {
		doc["craftsman"] = crafId
	}
	if !merId.IsZero() {
		doc["merchant"] = merId
	}
	if prod.Description != "" {
		doc["description"] = prod.Description
	}

	update := bson.M{"$set": doc}

	res := r.productsColl.FindOneAndUpdate(ctx, filter, update)

	err = res.Err()
	if err == nil {
		if opts.PopulateCraftsman {
			r.populateCraftsmanForProduct(prod)
		}
		if opts.PopulateMerchant {
			r.populateMerchantForProduct(prod)
		}
	}

	return err
}

func (r *repository) populateCraftsmanForProduct(prod *product.Product) {
	craftsman, err := r.FindUser(prod.CraftsmanID)
	if err == nil {
		prod.Craftsman = craftsman
	}
}

func (r *repository) populateMerchantForProduct(prod *product.Product) {
	merchant, err := r.FindMerchant(prod.MerchantID)
	if err == nil {
		prod.Merchant = merchant
	}
}
