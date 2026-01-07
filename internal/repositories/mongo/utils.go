package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/logger"
	"github.com/omareloui/odinls/internal/repositories/mongo/bsonutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func Get[T any](ctx context.Context, coll *mongo.Collection, filter *bson.M) ([]T, error) {
	l := logger.FromCtx(ctx)

	if filter == nil {
		filter = &bson.M{}
	}

	var cursor *mongo.Cursor
	var err error

	cursor, err = coll.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Info("no document found for the given filter", zap.String("collection", coll.Name()), zap.Any("filter", filter))
			return nil, errs.ErrDocumentNotFound
		}
		l.Warn("error finding in collection", zap.String("collection", coll.Name()), zap.Error(err), zap.Any("filter", filter))
		return nil, err
	}

	res := new([]T)
	if err := cursor.All(ctx, res); err != nil {
		l.Error("error decoding after finding", zap.String("collection", coll.Name()), zap.Error(err), zap.Any("filter", filter))
		return nil, err
	}

	l.Info("found documents", zap.String("collection", coll.Name()), zap.Any("filter", filter), zap.Int("found count", len(*res)))
	return *res, nil
}

func GetAll[T any](ctx context.Context, coll *mongo.Collection) ([]T, error) {
	return Get[T](ctx, coll, nil)
}

func GetByIDs[T any](ctx context.Context, coll *mongo.Collection, ids []string) ([]T, error) {
	objIDs := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errs.ErrInvalidID
		}

		objIDs[i] = objID
	}
	filter := bson.M{"_id": bson.M{"$in": objIDs}}

	return Get[T](ctx, coll, &filter)
}

func GetOne[T any](ctx context.Context, coll *mongo.Collection, filter bson.M) (*T, error) {
	l := logger.FromCtx(ctx)

	res := new(T)

	if err := coll.FindOne(ctx, filter).Decode(res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Info("no document found for the given filter", zap.String("collection", coll.Name()), zap.Any("filter", filter))
			return nil, errs.ErrDocumentNotFound
		}
		l.Error("error finding one documentid", zap.String("collection", coll.Name()), zap.Error(err), zap.Any("filter", filter))
		return nil, err
	}

	l.Info("found document one", zap.String("collection", coll.Name()), zap.Any("filter", filter), zap.Any("record", res))

	return res, nil
}

func GetByID[T any](ctx context.Context, coll *mongo.Collection, id string) (*T, error) {
	l := logger.FromCtx(ctx)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		l.Error("error getting id from hex", zap.String("collection", coll.Name()), zap.Error(err), zap.String("hex", id))
		return nil, errs.ErrInvalidID
	}

	filter := bson.M{"_id": objID}

	return GetOne[T](ctx, coll, filter)
}

func InsertStruct[T any](ctx context.Context, coll *mongo.Collection, rec *T, opts ...bsonutils.OptsFunc) (*T, error) {
	l := logger.FromCtx(ctx)

	bu := bsonutils.NewBsonUtils()
	doc, err := bu.MarshalInsertBsonD(rec, opts...)
	if err != nil {
		l.Error("error marshaling the record to bson d", zap.String("collection", coll.Name()), zap.Error(err), zap.Any("record", rec))
		return nil, err
	}

	res, err := coll.InsertOne(ctx, doc)
	if err != nil {
		if ok := mongo.IsDuplicateKeyError(err); ok {
			if se := mongo.ServerError(nil); errors.As(err, &se) {
				l.Warn("error duplicate document on inserting the document", zap.String("collection", coll.Name()), zap.Error(err), zap.Error(se), zap.Any("record", rec))
			}
			return nil, errs.ErrDocumentAlreadyExists
		}
		l.Error("error inserting document", zap.String("collection", coll.Name()), zap.Error(err), zap.Any("record", rec))
		return nil, err
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()

	l.Info("inserted the document", zap.String("collection", coll.Name()), zap.String("new_id", id), zap.Any("record", rec))

	newDoc, err := GetByID[T](ctx, coll, id)
	if err != nil {
		l.Error("error getting the inserted document", zap.String("collection", coll.Name()), zap.Error(err), zap.String("id", id), zap.Any("record", rec))
		return nil, err
	}

	return newDoc, nil
}

func UpdateStructByID[T any](ctx context.Context, coll *mongo.Collection, id string, rec *T, opts ...bsonutils.OptsFunc) (*T, error) {
	bu := bsonutils.NewBsonUtils()
	doc, err := bu.MarshalUpdateBsonD(rec, opts...)
	if err != nil {
		l := logger.FromCtx(ctx)
		l.Error("error marshaling the record to bson d", zap.String("collection", coll.Name()), zap.Error(err), zap.String("id", id), zap.Any("record", rec))
		return nil, err
	}

	return UpdateByID[T](ctx, coll, id, bson.M{"$set": doc})
}

func UpdateByID[T any](ctx context.Context, coll *mongo.Collection, id string, update bson.M) (*T, error) {
	l := logger.FromCtx(ctx)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		l.Error("error getting id from hex", zap.String("collection", coll.Name()), zap.Error(err), zap.String("hex", id))
		return nil, errs.ErrInvalidID
	}

	filter := bson.M{"_id": objID}

	err = UpdateOne[T](ctx, coll, filter, update)
	if err != nil {
		return nil, err
	}

	updatedDoc, err := GetByID[T](ctx, coll, id)
	if err != nil {
		l.Error("error getting the updated document", zap.String("collection", coll.Name()), zap.Error(err), zap.String("id", id), zap.Any("update", update))
		return nil, err
	}

	return updatedDoc, nil
}

func UpdateOne[T any](ctx context.Context, coll *mongo.Collection, filter bson.M, update bson.M) error {
	l := logger.FromCtx(ctx)

	res := coll.FindOneAndUpdate(ctx, filter, update)

	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Info("no document found for the given filter", zap.String("collection", coll.Name()), zap.Any("filter", filter))
			return errs.ErrDocumentNotFound
		}

		if ok := mongo.IsDuplicateKeyError(err); ok {
			if se := mongo.ServerError(nil); errors.As(err, &se) {
				l.Warn("error duplicate document on updating the document", zap.String("collection", coll.Name()), zap.Error(err), zap.Any("filter", filter), zap.Error(se), zap.Any("update", update))
			}
			return errs.ErrDocumentAlreadyExists
		}
		l.Error("error updating document", zap.String("collection", coll.Name()), zap.Error(err), zap.Any("filter", filter), zap.Any("update", update))
		return err
	}

	l.Info("updated the document", zap.String("collection", coll.Name()), zap.Any("filter", filter), zap.Any("update", update))

	return nil
}

type populateOpts struct {
	include      bool
	isMany       bool
	from         string
	localField   string
	foreignField string
	as           string
}

func PopulateAggregation[T any](ctx context.Context, coll *mongo.Collection, filterStages bson.A, popOpts ...populateOpts) ([]T, error) {
	pipeline := filterStages

	for _, pop := range popOpts {
		if pop.include {
			pipeline = append(pipeline, bson.M{
				"$lookup": bson.M{
					"from":         pop.from,
					"localField":   pop.localField,
					"foreignField": pop.foreignField,
					"as":           pop.as,
				},
			})

			if !pop.isMany {
				pipeline = append(pipeline, bson.M{
					"$unwind": bson.M{
						"path":                       fmt.Sprintf("$%s", pop.as),
						"preserveNullAndEmptyArrays": true,
					},
				})
			}
		}
	}

	cur, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	docs := new([]T)
	if err := cur.All(ctx, docs); err != nil {
		return nil, err
	}
	return *docs, nil
}

func PopulateAggregationByID[T any](ctx context.Context, coll *mongo.Collection, id string, popOpts ...populateOpts) (*T, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	pipeline := bson.A{
		bson.M{"$match": bson.M{"_id": objID}},
	}

	docs, err := PopulateAggregation[T](ctx, coll, pipeline, popOpts...)
	if err != nil {
		return nil, err
	}

	return &docs[0], nil
}
