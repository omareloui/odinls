package product

import jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"

type ProductService interface {
	GetProducts(*jwtadapter.JwtAccessClaims, ...RetrieveOptsFunc) ([]Product, error)
	GetProductByID(*jwtadapter.JwtAccessClaims, string, ...RetrieveOptsFunc) (*Product, error)
	GetCurrentMerchantProducts(*jwtadapter.JwtAccessClaims, ...RetrieveOptsFunc) ([]Product, error)
	CreateProduct(*jwtadapter.JwtAccessClaims, *Product, ...RetrieveOptsFunc) error
	UpdateClientByID(*jwtadapter.JwtAccessClaims, string, *Product, ...RetrieveOptsFunc) error
}
