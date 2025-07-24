package supplier

import jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"

type SupplierService interface {
	GetSuppliers(claims *jwtadapter.JwtAccessClaims) ([]Supplier, error)
	GetSupplierByID(claims *jwtadapter.JwtAccessClaims, id string) (*Supplier, error)
	CreateSupplier(claims *jwtadapter.JwtAccessClaims, supplier *Supplier) (*Supplier, error)
	UpdateSupplierByID(claims *jwtadapter.JwtAccessClaims, id string, supplier *Supplier) (*Supplier, error)
}
