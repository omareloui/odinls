package mongo

import (
	"errors"
	"fmt"

	"github.com/omareloui/odinls/internal/application/core/counter"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/errs"
	"go.mongodb.org/mongo-driver/bson"
)

const amountToIncrement = 1

func (r *repository) createCounter(cntr *counter.Counter) (*counter.Counter, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	return InsertStruct(ctx, r.countersColl, cntr)
}

func (r *repository) getCounter() (*counter.Counter, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	cntr, err := GetOne[counter.Counter](ctx, r.countersColl, bson.M{})
	if err != nil {
		if errors.Is(err, errs.ErrDocumentNotFound) {
			codes := product.CategoriesCodes()
			pcodes := make(counter.ProductsCodes, len(codes))
			for _, code := range codes {
				pcodes[code] = 0
			}
			return r.createCounter(&counter.Counter{ProductsCodes: pcodes})
		}

		return nil, err
	}

	return cntr, nil
}

func (r *repository) AddOneToProduct(category string) (uint8, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	filter := bson.M{}
	update := bson.M{
		"$inc": bson.M{fmt.Sprintf("products_codes.%s", category): amountToIncrement},
	}

	err := UpdateOne[counter.Counter](ctx, r.countersColl, filter, update)
	if err != nil {
		return 0, err
	}

	cntr, err := r.getCounter()
	if err != nil {
		return 0, err
	}

	return cntr.ProductsCodes[category], nil
}

func (r *repository) AddOneToOrder() (uint, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	filter := bson.M{}
	update := bson.M{
		"$inc": bson.M{"orders_number": amountToIncrement},
	}

	err := UpdateOne[counter.Counter](ctx, r.countersColl, filter, update)
	if err != nil {
		return 0, err
	}

	cntr, err := r.getCounter()
	if err != nil {
		return 0, err
	}

	return cntr.OrdersNumber, nil
}
