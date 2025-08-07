package product

import jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"

type ProductService interface {
	GetProducts(claims *jwtadapter.AccessClaims, opts ...RetrieveOptsFunc) ([]Product, error)
	GetProductByID(claims *jwtadapter.AccessClaims, id string, opts ...RetrieveOptsFunc) (*Product, error)
	GetProductByVariantID(claims *jwtadapter.AccessClaims, id string, opts ...RetrieveOptsFunc) (*Product, error)
	CreateProduct(claims *jwtadapter.AccessClaims, prod *Product, opts ...RetrieveOptsFunc) (*Product, error)
	UpdateProductByID(claims *jwtadapter.AccessClaims, id string, prod *Product, opts ...RetrieveOptsFunc) (*Product, error)
}
