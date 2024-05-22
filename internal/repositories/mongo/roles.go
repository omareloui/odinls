package mongo

import (
	"errors"
	"time"

	"github.com/omareloui/odinls/internal/application/core/role"
	"github.com/omareloui/odinls/internal/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetRoles() ([]role.Role, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	filter := bson.M{}
	cursor, err := r.rolesColl.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, role.ErrRoleNotFound
		}
		return nil, err
	}

	roles := new([]role.Role)
	if err := cursor.All(ctx, roles); err != nil {
		return nil, err
	}

	return *roles, nil
}

func (r *repository) FindRole(id string) (*role.Role, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	rol := &role.Role{}
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrInvalidID
	}
	filter := bson.M{"_id": objId}
	err = r.rolesColl.FindOne(ctx, filter).Decode(rol)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, role.ErrRoleNotFound
		}
		return nil, err
	}

	return rol, nil
}

func (r *repository) FindRoleByName(name string) (*role.Role, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	rol := &role.Role{}
	filter := bson.M{"name": name}
	err := r.rolesColl.FindOne(ctx, filter).Decode(rol)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, role.ErrRoleNotFound
		}
		return nil, err
	}

	return rol, nil
}

func (r *repository) CreateRole(dto *role.Role) error {
	ctx, cancel := r.newCtx()
	defer cancel()

	res, err := r.rolesColl.InsertOne(ctx, dto)

	if err == nil {
		dto.ID = res.InsertedID.(primitive.ObjectID).Hex()
	}

	if ok := mongo.IsDuplicateKeyError(err); ok {
		if se := mongo.ServerError(nil); errors.As(err, &se) {
			if se.HasErrorMessage("{ name: ") {
				return role.ErrRoleNameAlreadyExists
			}
		}
	}

	return err
}

func (r *repository) SeedRoles(dto []string) error {
	ctx, cancel := r.newCtx()
	defer cancel()

	now := time.Now()

	roles := []interface{}{}

	for _, rol := range dto {
		roles = append(roles, role.Role{
			Name:      rol,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	_, err := r.rolesColl.InsertMany(ctx, roles)

	return err
}
