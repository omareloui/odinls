package bson_utils

import (
	"errors"
	"slices"

	"github.com/omareloui/odinls/internal/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidBsonKey = errors.New("invalid bson key")

type BsonUtils struct{}

type opts struct {
	objIDKeys  []string
	removeKeys []string
}

type optsFunc func(*opts)

func NewBsonUtils() *BsonUtils {
	return &BsonUtils{}
}

func (bu *BsonUtils) MarshalBsonD(rec any, opts ...optsFunc) (bson.D, error) {
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

	for _, k := range o.objIDKeys {
		if err := bu.setKeyAsObjectID(doc, k); err != nil {
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

func (bu *BsonUtils) AddObjectIDKey(key string) optsFunc {
	return func(opts *opts) {
		if opts.objIDKeys == nil {
			opts.objIDKeys = []string{}
		}
		opts.objIDKeys = append(opts.objIDKeys, key)
	}
}

func (bu *BsonUtils) RemoveKey(key string) optsFunc {
	return func(opts *opts) {
		if opts.removeKeys == nil {
			opts.removeKeys = []string{}
		}
		opts.removeKeys = append(opts.removeKeys, key)
	}
}

func (bu *BsonUtils) setKeyAsObjectID(doc bson.D, key string) error {
	for i, obj := range doc {
		if obj.Key == key {
			strval, ok := obj.Value.(string)
			if !ok {
				return errs.ErrInvalidID
			}
			objId, err := primitive.ObjectIDFromHex(strval)
			if err != nil {
				return errs.ErrInvalidID
			}

			doc[i].Value = objId
			return nil
		}
	}
	return ErrInvalidBsonKey
}

func (bu *BsonUtils) removeKey(doc bson.D, key string) (bson.D, error) {
	for i, obj := range doc {
		if obj.Key == key {
			return slices.Delete(doc, i, i+1), nil
		}
	}
	return doc, ErrInvalidBsonKey
}

func (bu *BsonUtils) parseOpts(funcs ...optsFunc) *opts {
	o := &opts{objIDKeys: []string{}, removeKeys: []string{}}
	for _, fun := range funcs {
		fun(o)
	}
	return o
}
