package mongo

import (
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *repository) GetProducts(options ...product.RetrieveOptsFunc) ([]product.Product, error) {
	opts := product.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	return PopulateAggregation[product.Product](ctx, r.productsColl, bson.A{}, r.productPotsToPopulateOpts(opts)...)
}

func (r *repository) GetProductByID(id string, options ...product.RetrieveOptsFunc) (*product.Product, error) {
	opts := product.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	return PopulateAggregationByID[product.Product](ctx, r.productsColl, id, r.productPotsToPopulateOpts(opts)...)
}

func (r *repository) GetProductByVariantID(id string, options ...product.RetrieveOptsFunc) (*product.Product, error) {
	opts := product.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	docs, err := PopulateAggregation[product.Product](ctx, r.productsColl,
		bson.A{
			bson.M{"$match": bson.M{"variants._id": objID}},
		},
		r.productPotsToPopulateOpts(opts)...)
	if err != nil {
		return nil, err
	}

	return &docs[0], nil
}

func (r *repository) GetProductByIDAndVariantID(id string, variantId string, options ...product.RetrieveOptsFunc) (*product.Product, error) {
	opts := product.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	varObjID, err := primitive.ObjectIDFromHex(variantId)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	docs, err := PopulateAggregation[product.Product](ctx, r.productsColl,
		bson.A{
			bson.M{
				"$match": bson.M{
					"_id":          objID,
					"variants._id": varObjID,
				},
			},
		},
		r.productPotsToPopulateOpts(opts)...)
	if err != nil {
		return nil, err
	}

	return &docs[0], nil
}

func (r *repository) CreateProduct(prod *product.Product, options ...product.RetrieveOptsFunc) (*product.Product, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	doc, err := InsertStruct(ctx, r.productsColl, prod)
	if err != nil {
		return nil, err
	}

	return r.GetProductByID(doc.ID, options...)
}

func (r *repository) UpdateProductByID(id string, prod *product.Product, options ...product.RetrieveOptsFunc) (*product.Product, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	doc, err := UpdateStructByID(ctx, r.productsColl, id, prod)
	if err != nil {
		return nil, err
	}

	return r.GetProductByID(doc.ID, options...)
}

func (r *repository) productPotsToPopulateOpts(opts *product.RetrieveOpts) []populateOpts {
	return []populateOpts{{
		include:      opts.PopulateUsedMaterial,
		from:         usersCollectionName,
		foreignField: "_id",
		localField:   "variant.$.material_usage.$.material_id",
		as:           "variant.$.material_usage.$.populated_material",
	}}
}
