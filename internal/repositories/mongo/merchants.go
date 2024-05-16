package mongo

import (
	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/errmsgs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) FindMerchant(id string) (*merchant.Merchant, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	m := &merchant.Merchant{}
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errmsgs.ErrInvalidID
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

	// TODO: find a better way to map from struct to bson
	_, err := r.merchantsColl.InsertOne(
		ctx,
		bson.M{
			"name":       merchant.Name,
			"logo":       merchant.Logo,
			"created_at": merchant.CreatedAt,
			"updated_at": merchant.UpdatedAt,
		},
	)

	return err
}
