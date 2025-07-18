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

	return prod, nil
}

func (r *repository) GetProductByVariantID(id string, options ...product.RetrieveOptsFunc) (*product.Product, error) {
	opts := product.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	filter := bson.M{"variants._id": objId}

	prod := &product.Product{}

	err = r.productsColl.FindOne(ctx, filter).Decode(prod)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, product.ErrVariantNotFound
		}
		return nil, err
	}

	if opts.PopulateCraftsman {
		r.populateCraftsmanForProduct(prod)
	}

	return prod, nil
}

func (r *repository) GetProductByIDAndVariantID(id string, variantId string, options ...product.RetrieveOptsFunc) (*product.Product, error) {
	opts := product.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	varObjId, err := primitive.ObjectIDFromHex(variantId)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	filter := bson.M{"_id": objId, "variants._id": varObjId}

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

	return prod, nil
}

func (r *repository) CreateProduct(prod *product.Product, options ...product.RetrieveOptsFunc) error {
	opts := product.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	doc, err := mapProductToMongoDoc(prod)
	if err != nil {
		return err
	}

	res, err := r.productsColl.InsertOne(ctx, doc)

	if err == nil {
		prod.ID = res.InsertedID.(primitive.ObjectID).Hex()
		if opts.PopulateCraftsman {
			r.populateCraftsmanForProduct(prod)
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

	doc, err := mapProductToMongoDoc(prod)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objId}
	update := bson.M{"$set": doc}

	res := r.productsColl.FindOneAndUpdate(ctx, filter, update)

	err = res.Err()
	if err == nil {
	}

	return err
}

func mapProductToMongoDoc(prod *product.Product) (bson.M, error) {
	var crafId primitive.ObjectID

	var err error

	variants := make(bson.A, len(prod.Variants))

	now := time.Now()

	doc := bson.M{
		"number":     prod.Number,
		"name":       prod.Name,
		"category":   prod.Category,
		"variants":   variants,
		"updated_at": now,
	}

	if prod.CreatedAt.IsZero() {
		doc["created_at"] = now
	}
	if !crafId.IsZero() {
		doc["craftsman"] = crafId
	}
	if prod.Description != "" {
		doc["description"] = prod.Description
	}

	for i, variant := range prod.Variants {
		doc["variants"].(bson.A)[i] = bson.M{
			"suffix":          variant.Suffix,
			"name":            variant.Name,
			"materials_cost":  variant.MaterialsCost,
			"price":           variant.Price,
			"wholesale_price": variant.WholesalePrice,
			"time_to_craft":   variant.TimeToCraft,
			"product_ref":     variant.ProductRef,
		}
		if variant.ID == "" {
			doc["variants"].(bson.A)[i].(bson.M)["_id"] = primitive.NewObjectID()
		} else {
			doc["variants"].(bson.A)[i].(bson.M)["_id"], err = primitive.ObjectIDFromHex(variant.ID)
			if err != nil {
				return nil, errs.ErrInvalidID
			}
		}
		if variant.Description != "" {
			doc["variants"].(bson.A)[i].(bson.M)["description"] = variant.Description
		}
	}

	return doc, nil
}
