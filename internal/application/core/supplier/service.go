package supplier

import jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"

type SupplierService interface {
	GetSuppliers(claims *jwtadapter.AccessClaims) ([]Supplier, error)
	GetSupplierByID(claims *jwtadapter.AccessClaims, id string) (*Supplier, error)
	CreateSupplier(claims *jwtadapter.AccessClaims, supplier *Supplier) (*Supplier, error)
	UpdateSupplierByID(claims *jwtadapter.AccessClaims, id string, supplier *Supplier) (*Supplier, error)
}
