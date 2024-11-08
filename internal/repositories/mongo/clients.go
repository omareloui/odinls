package mongo

import (
	"errors"

	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetClients(options ...client.RetrieveOptsFunc) ([]client.Client, error) {
	opts := client.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	filter := bson.M{}

	var cursor *mongo.Cursor
	var err error

	if !opts.PopulateMerchant {
		cursor, err = r.clientsColl.Find(ctx, filter)
	} else {
		cursor, err = r.clientsColl.Aggregate(ctx, bson.A{
			bson.M{
				"$lookup": bson.M{
					"from":         merchantsCollectionName,
					"localField":   "merchant",
					"foreignField": "_id",
					"as":           "populatedMerchant",
				},
			},
			bson.M{"$unwind": "$populatedMerchant"},
		})
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, client.ErrClientNotFound
		}
		return nil, err
	}

	clients := new([]client.Client)
	if err := cursor.All(ctx, clients); err != nil {
		return nil, err
	}

	return *clients, nil
}

func (r *repository) GetClientsByMerchantID(merchantId string, options ...client.RetrieveOptsFunc) ([]client.Client, error) {
	opts := client.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	mrId, err := primitive.ObjectIDFromHex(merchantId)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	filter := bson.M{"merchant": mrId}

	var cursor *mongo.Cursor

	if !opts.PopulateMerchant {
		cursor, err = r.clientsColl.Find(ctx, filter)
	} else {
		cursor, err = r.clientsColl.Aggregate(ctx, bson.A{
			bson.M{"$match": filter},
			bson.M{
				"$lookup": bson.M{
					"from":         merchantsCollectionName,
					"localField":   "merchant",
					"foreignField": "_id",
					"as":           "populatedMerchant",
				},
			},
			bson.M{"$unwind": "$populatedMerchant"},
		})
	}

	if err != nil {
		return nil, err
	}

	clients := new([]client.Client)
	if err := cursor.All(ctx, clients); err != nil {
		return nil, err
	}

	return *clients, nil
}

func (r *repository) GetClientByID(id string, options ...client.RetrieveOptsFunc) (*client.Client, error) {
	opts := client.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	filter := bson.M{"_id": objId}

	cli := &client.Client{}

	err = r.clientsColl.FindOne(ctx, filter).Decode(cli)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, client.ErrClientNotFound
		}
		return nil, err
	}

	if opts.PopulateMerchant {
		r.populateMerchantForClient(cli)
	}

	return cli, nil
}

func (r *repository) CreateClient(cli *client.Client, options ...client.RetrieveOptsFunc) error {
	opts := client.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	doc, err := r.bu.MarshalBsonD(cli, r.bu.WithObjectID("merchant"), r.bu.WithFieldToRemove("notes"))
	if err != nil {
		return err
	}

	res, err := r.clientsColl.InsertOne(ctx, doc)

	if err == nil {
		cli.ID = res.InsertedID.(primitive.ObjectID).Hex()
		if opts.PopulateMerchant {
			r.populateMerchantForClient(cli)
		}
	}

	if ok := mongo.IsDuplicateKeyError(err); ok {
		if se := mongo.ServerError(nil); errors.As(err, &se) {
			if se.HasErrorMessage(" name: ") && se.HasErrorMessage(" merchant: ") {
				return client.ErrClientExistsForMerchant
			}
		}
	}

	return err
}

func (r *repository) UpdateClientByID(id string, cli *client.Client, options ...client.RetrieveOptsFunc) error {
	opts := client.ParseRetrieveOpts(options...)

	ctx, cancel := r.newCtx()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errs.ErrInvalidID
	}

	filter := bson.M{"_id": objId}

	doc, err := r.bu.MarshalBsonD(cli,
		r.bu.WithFieldToRemove("_id"),
		r.bu.WithFieldToRemove("merchant"),
		r.bu.WithFieldToRemove("created_at"))
	if err != nil {
		return err
	}

	res := r.clientsColl.FindOneAndUpdate(ctx, filter, bson.M{"$set": doc})

	err = res.Err()
	if err == nil {
		cli.ID = id
		if opts.PopulateMerchant {
			r.populateMerchantForClient(cli)
		}
	}

	if ok := mongo.IsDuplicateKeyError(err); ok {
		if se := mongo.ServerError(nil); errors.As(err, &se) {
			if se.HasErrorMessage(" name: ") && se.HasErrorMessage(" merchant: ") {
				return client.ErrClientExistsForMerchant
			}
		}
	}

	return err
}

func (r *repository) populateMerchantForClient(cli *client.Client) {
	merchant, err := r.FindMerchant(cli.MerchantID)
	if err == nil {
		cli.Merchant = merchant
	}
}
