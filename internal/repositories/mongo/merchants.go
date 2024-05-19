package mongo

import (
	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetMerchants() ([]merchant.Merchant, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	m := []merchant.Merchant{}

	cursor, err := r.merchantsColl.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func (r *repository) FindMerchant(id string) (*merchant.Merchant, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	m := &merchant.Merchant{}
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrInvalidID
	}
	filter := bson.M{"_id": objId}
	err = r.merchantsColl.FindOne(ctx, filter).Decode(m)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchant.ErrMerchantNotFound
		}
		return nil, err
	}

	return m, nil
}

func (r *repository) CreateMerchant(merchant *merchant.Merchant) error {
	ctx, cancel := r.newCtx()
	defer cancel()

	res, err := r.merchantsColl.InsertOne(ctx, merchant)

	if err == nil {
		merchant.ID = res.InsertedID.(primitive.ObjectID).Hex()
	}

	return err
}

func (r *repository) UpdateMerchantByID(id string, mer *merchant.Merchant) error {
	ctx, cancel := r.newCtx()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errs.ErrInvalidID
	}

	filter := bson.D{{Key: "_id", Value: objId}}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.M{
				"name": mer.Name,
				"logo": mer.Logo,
			},
		},
	}

	updated, err := r.merchantsColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if updated.ModifiedCount == 0 {
		return merchant.ErrMerchantNotFound
	}
	return nil
}
