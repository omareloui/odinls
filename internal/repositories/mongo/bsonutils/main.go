// Package bsonutils provides utility functions for working with BSON documents in MongoDB.
package bsonutils

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/omareloui/odinls/internal/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidBsonKey = errors.New("invalid bson key")

type BsonUtils struct{}

type opts struct {
	objIDKeys     []string
	removeKeys    []string
	stringifyKeys []string
	append        bson.D
}

type OptsFunc func(*opts)

func NewBsonUtils() *BsonUtils {
	return &BsonUtils{}
}

func (bu *BsonUtils) MarshalBsonD(rec any, opts ...OptsFunc) (bson.D, error) {
	o := bu.parseOpts(opts...)
	bdoc, err := bson.Marshal(rec)
	if err != nil {
		return nil, err
	}

	var doc bson.D
	err = bson.Unmarshal(bdoc, &doc)
	if err != nil {
		return nil, err
	}

	for _, e := range o.append {
		doc = bu.append(doc, e)
	}

	for _, k := range o.objIDKeys {
		if err := bu.setKeyAsObjectID(doc, k); err != nil {
			return nil, err
		}
	}

	for _, k := range o.stringifyKeys {
		if err := bu.setKeyAsString(doc, k); err != nil {
			return nil, err
		}
	}

	for _, k := range o.removeKeys {
		if doc, err = bu.removeKey(doc, k); err != nil && !errors.Is(err, ErrInvalidBsonKey) {
			return nil, err
		}
	}

	return doc, nil
}

func (bu *BsonUtils) MarshalInsertBsonD(rec any, opts ...OptsFunc) (bson.D, error) {
	now := time.Now()
	opts = append(opts,
		WithFieldToRemove("_id"),
		WithFieldToAdd("created_at", now),
		WithFieldToAdd("updated_at", now),
	)
	return bu.MarshalBsonD(rec, opts...)
}

func (bu *BsonUtils) MarshalUpdateBsonD(rec any, opts ...OptsFunc) (bson.D, error) {
	opts = append(opts,
		WithFieldToRemove("_id"),
		WithFieldToRemove("created_at"),
		WithFieldToAdd("updated_at", time.Now()))
	return bu.MarshalBsonD(rec, opts...)
}

func WithObjectID(key string) OptsFunc {
	return func(opts *opts) {
		opts.objIDKeys = append(opts.objIDKeys, key)
	}
}

func WithStringfied(key string) OptsFunc {
	return func(opts *opts) {
		opts.stringifyKeys = append(opts.stringifyKeys, key)
	}
}

func WithFieldToAdd(key string, val any) OptsFunc {
	return func(opts *opts) {
		opts.append = append(opts.append, bson.E{Key: key, Value: val})
	}
}

func WithFieldToRemove(key string) OptsFunc {
	return func(opts *opts) {
		if opts.removeKeys == nil {
			opts.removeKeys = []string{}
		}
		opts.removeKeys = append(opts.removeKeys, key)
	}
}

func diveAndOverride[T any, K any](doc bson.D, key string, cb func(currValue T) (K, error)) error {
	path := strings.Split(key, ".")

	for i, obj := range doc {
		if len(path) > 1 && obj.Key == path[0] {
			if subdoc, ok := obj.Value.(bson.D); ok {
				return diveAndOverride(subdoc, strings.Join(path[1:], "."), cb)
			}
		}

		if arr, ok := obj.Value.(bson.A); ok {
			for j, v := range arr {
				arrObj, isSubDoc := v.(bson.D)
				if !isSubDoc {
					continue
				}
				for k, subObj := range arrObj {
					if subObj.Key == key {
						newValue, err := cb(subObj.Value.(T))
						if err != nil {
							return err
						}
						doc[i].Value.(bson.A)[j].(bson.D)[k].Value = newValue
					}
				}
			}
		}

		if obj.Key == key {
			newValue, err := cb(obj.Value.(T))
			if err != nil {
				return err
			}
			doc[i].Value = newValue
			return nil
		}
	}
	return ErrInvalidBsonKey
}

func (bu *BsonUtils) append(doc bson.D, e bson.E) bson.D {
	doc, _ = bu.removeKey(doc, e.Key)
	return append(doc, e)
}

func (bu *BsonUtils) removeKey(doc bson.D, key string) (bson.D, error) {
	for i, obj := range doc {
		if obj.Key == key {
			return slices.Delete(doc, i, i+1), nil
		}
	}
	return doc, ErrInvalidBsonKey
}

func (bu *BsonUtils) setKeyAsObjectID(doc bson.D, key string) error {
	return diveAndOverride(doc, key, getObjectID)
}

func (bu *BsonUtils) setKeyAsString(doc bson.D, key string) error {
	return diveAndOverride(doc, key, getStringified)
}

func (bu *BsonUtils) parseOpts(funcs ...OptsFunc) *opts {
	o := &opts{}
	for _, fun := range funcs {
		fun(o)
	}
	return o
}

func getStringified(val fmt.Stringer) (string, error) {
	return val.String(), nil
}

func getObjectID(val any) (primitive.ObjectID, error) {
	strval, ok := val.(string)
	if !ok {
		return primitive.ObjectID{}, errs.ErrInvalidID
	}
	objID, err := primitive.ObjectIDFromHex(strval)
	if err != nil {
		return primitive.ObjectID{}, errs.ErrInvalidID
	}

	return objID, nil
}
