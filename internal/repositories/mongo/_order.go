package mongo

import (
	"github.com/omareloui/odinls/internal/application/core/order"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/repositories/mongo/bsonutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetOrders(options ...order.RetrieveOptsFunc) ([]order.Order, error) {
	opts := order.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	pipeline := bson.A{}
	buildPipelineForOrdersFromOpts(&pipeline, opts)

	cur, err := r.ordersColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	ords := new([]order.Order)
	if err := cur.All(ctx, ords); err != nil {
		return nil, err
	}

	return *ords, nil
}

func (r *repository) GetOrderByID(id string, options ...order.RetrieveOptsFunc) (*order.Order, error) {
	opts := order.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	filter := bson.M{"_id": objId}

	ord := &order.Order{}

	err = r.ordersColl.FindOne(ctx, filter).Decode(ord)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, order.ErrOrderNotFound
		}
		return nil, err
	}

	r.populateOrder(ord, opts)
	return ord, nil
}

func (r *repository) CreateOrder(ord *order.Order, options ...order.RetrieveOptsFunc) error {
	opts := order.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	doc, err := mapOrderToMongoDoc(r.bu, ord)
	if err != nil {
		return err
	}

	res, err := r.ordersColl.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	ord.ID = res.InsertedID.(primitive.ObjectID).Hex()
	r.populateOrder(ord, opts)

	return nil
}

func (r *repository) UpdateOrderByID(id string, ord *order.Order, options ...order.RetrieveOptsFunc) error {
	opts := order.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errs.ErrInvalidID
	}

	doc, err := mapOrderToMongoDoc(r.bu, ord)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objId}
	update := bson.M{"$set": doc}

	res := r.ordersColl.FindOneAndUpdate(ctx, filter, update)

	err = res.Err()
	if err != nil {
		return err
	}

	r.populateOrder(ord, opts)
	return nil
}

func (r *repository) populateClientForOrder(ord *order.Order) {
	client, err := r.GetClientByID(ord.ClientID)
	if err == nil {
		ord.Client = client
	}
}

func (r *repository) populateCraftsmenForOrder(ord *order.Order) {
	craftsmen, err := r.FindUsersByIDs(ord.CraftsmenIDs)
	if err == nil {
		ord.Craftsmen = craftsmen
	}
}

func (r *repository) populateOrder(ord *order.Order, opts *order.RetrieveOpts) {
	if opts.PopulateClient {
		r.populateClientForOrder(ord)
	}
	if opts.PopulateCraftsmen {
		r.populateCraftsmenForOrder(ord)
	}
}

func buildPipelineForOrdersFromOpts(pipeline *bson.A, opts *order.RetrieveOpts) {
	if opts.PopulateCraftsmen {
		*pipeline = append(*pipeline, bson.M{
			"$lookup": bson.M{
				"from":         usersCollectionName,
				"localField":   "craftsmen",
				"foreignField": "_id",
				"as":           "populated_craftsmen",
			},
		}, bson.M{"$unwind": bson.M{"path": "$populated_craftsmen", "preserveNullAndEmptyArrays": true}})
	}
	if opts.PopulateClient {
		*pipeline = append(*pipeline, bson.M{
			"$lookup": bson.M{
				"from":         clientsCollectionName,
				"localField":   "client",
				"foreignField": "_id",
				"as":           "populated_client",
			},
		}, bson.M{"$unwind": bson.M{"path": "$populated_client", "preserveNullAndEmptyArrays": true}})
	}
	if opts.PopulateItemProducts {
		*pipeline = append(*pipeline, bson.M{
			"$lookup": bson.M{
				"from": productsCollectionName,
				// TODO(research): THIS IS AN ARRAY, HOW WILL IT WORK?
				"localField":   "items.product",
				"foreignField": "_id",
				"as":           "items.populated_product",
			},
		})
	}
	// TODO: if this's true, the product must be true, check how to make it work
	if opts.PopulateItemVariants {
		*pipeline = append(*pipeline, bson.M{
			"$lookup": bson.M{
				"from": productsCollectionName,
				// TODO(research): THIS IS AN ARRAY, HOW WILL IT WORK?
				"localField":   "items.variant",
				"foreignField": "variants._id",
				"as":           "items.populated_variant",
			},
		})
	}
}

func mapOrderToMongoDoc(bu *bsonutils.BsonUtils, ord *order.Order) (bson.D, error) {
	return bu.MarshalBsonD(ord,
		bu.WithObjectID("client"),
		bu.WithObjectID("craftsmen"),
		// TODO: does it work?
		bu.WithObjectID("items._id"),
		bu.WithObjectID("items.product"),
		bu.WithObjectID("items.variant"),
	)
}
