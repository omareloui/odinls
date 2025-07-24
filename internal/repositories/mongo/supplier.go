package mongo

import (
	"github.com/omareloui/odinls/internal/application/core/supplier"
)

func (r *repository) GetSuppliers() ([]supplier.Supplier, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	return GetAll[supplier.Supplier](ctx, r.suppliersColl)
}

func (r *repository) GetSupplierByID(id string) (*supplier.Supplier, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	return GetByID[supplier.Supplier](ctx, r.suppliersColl, id)
}

func (r *repository) CreateSupplier(sup *supplier.Supplier) (*supplier.Supplier, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	return InsertStruct(ctx, r.suppliersColl, sup)
}

func (r *repository) UpdateSupplierByID(id string, sup *supplier.Supplier) (*supplier.Supplier, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	return UpdateStructByID(ctx, r.suppliersColl, id, sup)
}
