package bson_utils

import (
	"errors"
	"slices"
	"strings"

	"github.com/omareloui/odinls/internal/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidBsonKey = errors.New("invalid bson key")

type BsonUtils struct{}

type opts struct {
	objIDKeys  []string
	removeKeys []string
	append     bson.D
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

	for _, k := range o.removeKeys {
		if doc, err = bu.removeKey(doc, k); err != nil && !errors.Is(err, ErrInvalidBsonKey) {
			return nil, err
		}
	}

	return doc, nil
}

func (bu *BsonUtils) WithObjectID(key string) OptsFunc {
	return func(opts *opts) {
		if opts.objIDKeys == nil {
			opts.objIDKeys = []string{}
		}
		opts.objIDKeys = append(opts.objIDKeys, key)
	}
}

func (bu *BsonUtils) WithFieldToAdd(key string, val any) OptsFunc {
	return func(opts *opts) {
		if opts.append == nil {
			opts.append = bson.D{}
		}
		opts.append = append(opts.append, bson.E{Key: key, Value: val})
	}
}

func (bu *BsonUtils) WithFieldToRemove(key string) OptsFunc {
	return func(opts *opts) {
		if opts.removeKeys == nil {
			opts.removeKeys = []string{}
		}
		opts.removeKeys = append(opts.removeKeys, key)
	}
}

func (bu *BsonUtils) setKeyAsObjectID(doc bson.D, key string) error {
	path := strings.Split(key, ".")
	for i, obj := range doc {
		if len(path) > 1 && obj.Key == path[0] {
			if subdoc, ok := obj.Value.(bson.D); ok {
				// FIXME: won't work with array?
				return bu.setKeyAsObjectID(subdoc, strings.Join(path[1:], "."))
			}
		}

		if arr, ok := obj.Value.(bson.A); ok {
			for j, v := range arr {
				for k, subObj := range v.(bson.D) {
					if subObj.Key == key {
						objId, err := getObjectID(obj.Value)
						if err != nil {
							return err
						}
						doc[i].Value.(bson.A)[j].(bson.D)[k].Value = objId
					}
				}
			}
		}

		if obj.Key == key {
			objId, err := getObjectID(obj.Value)
			if err != nil {
				return err
			}
			doc[i].Value = objId
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

func (bu *BsonUtils) parseOpts(funcs ...OptsFunc) *opts {
	o := &opts{objIDKeys: []string{}, removeKeys: []string{}}
	for _, fun := range funcs {
		fun(o)
	}
	return o
}

func getObjectID(val any) (primitive.ObjectID, error) {
	strval, ok := val.(string)
	if !ok {
		return primitive.ObjectID{}, errs.ErrInvalidID
	}
	objId, err := primitive.ObjectIDFromHex(strval)
	if err != nil {
		return primitive.ObjectID{}, errs.ErrInvalidID
	}

	return objId, nil
}
