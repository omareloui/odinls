package mongo

import (
	"github.com/omareloui/odinls/internal/application/core/material"
)

func (r *repository) GetMaterials(options ...material.RetrieveOptsFunc) ([]material.Material, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	mats, err := GetAll[material.Material](ctx, r.materialsColl)
	if err != nil {
		return nil, err
	}

	opts := material.ParseRetrieveOpts(options...)
	r.populateMaterials(mats, opts)

	return mats, nil
}

func (r *repository) GetMaterialByID(id string, options ...material.RetrieveOptsFunc) (*material.Material, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	m, err := GetByID[material.Material](ctx, r.materialsColl, id)
	if err != nil {
		return nil, err
	}

	opts := material.ParseRetrieveOpts(options...)
	r.populateMaterial(m, opts)

	return m, nil
}

func (r *repository) CreateMaterial(mat *material.Material, options ...material.RetrieveOptsFunc) (*material.Material, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	m, err := InsertStruct(ctx, r.materialsColl, mat)
	if err != nil {
		return nil, err
	}

	opts := material.ParseRetrieveOpts(options...)
	r.populateMaterial(mat, opts)

	return m, err
}

func (r *repository) UpdateMaterialByID(id string, mat *material.Material, options ...material.RetrieveOptsFunc) (*material.Material, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	m, err := UpdateStructByID(ctx, r.materialsColl, id, mat)
	if err != nil {
		return nil, err
	}

	opts := material.ParseRetrieveOpts(options...)
	r.populateMaterial(mat, opts)

	return m, nil
}

func (r *repository) populateMaterials(mats []material.Material, opts *material.RetrieveOpts) {
	for _, material := range mats {
		r.populateMaterial(&material, opts)
	}
}

func (r *repository) populateMaterial(mat *material.Material, opts *material.RetrieveOpts) {
	if opts.PopulateSupplier {
		sup, err := r.GetSupplierByID(mat.SupplierID)
		if err == nil {
			mat.Supplier = sup
		}
	}
}
