package mongo

import (
	"github.com/omareloui/odinls/internal/application/core/order"
	"github.com/omareloui/odinls/internal/repositories/mongo/bsonutils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) GetOrders(options ...order.RetrieveOptsFunc) ([]order.Order, error) {
	opts := order.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	return PopulateAggregation[order.Order](ctx, r.ordersColl, bson.A{}, r.orderOptsToPopulateOpts(opts)...)
}

func (r *repository) GetOrderByID(id string, options ...order.RetrieveOptsFunc) (*order.Order, error) {
	opts := order.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	return PopulateAggregationByID[order.Order](ctx, r.ordersColl, id, r.orderOptsToPopulateOpts(opts)...)
}

func (r *repository) CreateOrder(ord *order.Order, options ...order.RetrieveOptsFunc) (*order.Order, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	res, err := InsertStruct(ctx, r.ordersColl, ord,
		bsonutils.WithObjectID("client"),
		bsonutils.WithObjectID("craftsmen"),
		bsonutils.WithObjectID("items._id"),
		bsonutils.WithObjectID("items.craftsman"),
		bsonutils.WithObjectID("items.snapshot.product"),
		bsonutils.WithObjectID("items.snapshot.variant_id"),
	)
	if err != nil {
		return nil, err
	}

	return r.GetOrderByID(res.ID, options...)
}

func (r *repository) UpdateOrderByID(id string, ord *order.Order, options ...order.RetrieveOptsFunc) (*order.Order, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	_, err := UpdateStructByID(ctx, r.ordersColl, id, ord,
		bsonutils.WithObjectID("client"),
		bsonutils.WithObjectID("craftsmen"),
		bsonutils.WithObjectID("items._id"),
		bsonutils.WithObjectID("items.craftsman"),
		bsonutils.WithObjectID("items.snapshot.product"),
		bsonutils.WithObjectID("items.snapshot.variant_id"),
	)
	if err != nil {
		return nil, err
	}
	return r.GetOrderByID(id, options...)
}

func (r *repository) orderOptsToPopulateOpts(opts *order.RetrieveOpts) []populateOpts {
	return []populateOpts{
		{
			include:      opts.PopulateClient,
			from:         clientsCollectionName,
			foreignField: "_id",
			localField:   "client",
			as:           "populated_client",
		},
		{
			include:      opts.PopulateCraftsmen,
			from:         usersCollectionName,
			foreignField: "_id",
			localField:   "items.$.craftsman",
			as:           "items.$.populated_craftsman",
		},
		{
			include:      opts.PopulateItemProducts,
			from:         productsCollectionName,
			foreignField: "_id",
			localField:   "items.$.snapshot.product",
			as:           "items.$.populated_product",
		},
	}
}
