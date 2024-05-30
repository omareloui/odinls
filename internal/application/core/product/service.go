package product

import jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"

type ProductService interface {
	GetProducts(claims *jwtadapter.JwtAccessClaims, opts ...RetrieveOptsFunc) ([]Product, error)
	GetProductByID(claims *jwtadapter.JwtAccessClaims, id string, opts ...RetrieveOptsFunc) (*Product, error)
	GetCurrentMerchantProducts(claims *jwtadapter.JwtAccessClaims, opts ...RetrieveOptsFunc) ([]Product, error)
	CreateProduct(claims *jwtadapter.JwtAccessClaims, prod *Product, opts ...RetrieveOptsFunc) error
	UpdateProductByID(claims *jwtadapter.JwtAccessClaims, id string, prod *Product, opts ...RetrieveOptsFunc) error
}
