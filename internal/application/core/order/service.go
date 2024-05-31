package order

import jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"

type OrderService interface {
	GetOrders(claims *jwtadapter.JwtAccessClaims, opts ...RetrieveOptsFunc) ([]Order, error)
	GetOrderByID(claims *jwtadapter.JwtAccessClaims, id string, opts ...RetrieveOptsFunc) (*Order, error)
	GetCurrentMerchantOrders(claims *jwtadapter.JwtAccessClaims, opts ...RetrieveOptsFunc) ([]Order, error)
	CreateOrder(claims *jwtadapter.JwtAccessClaims, ord *Order, opts ...RetrieveOptsFunc) error
	UpdateOrderByID(claims *jwtadapter.JwtAccessClaims, id string, ord *Order, opts ...RetrieveOptsFunc) error
}
