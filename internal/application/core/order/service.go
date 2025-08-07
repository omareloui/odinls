package order

import jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"

type OrderService interface {
	GetOrders(claims *jwtadapter.AccessClaims, opts ...RetrieveOptsFunc) ([]Order, error)
	GetOrderByID(claims *jwtadapter.AccessClaims, id string, opts ...RetrieveOptsFunc) (*Order, error)
	CreateOrder(claims *jwtadapter.AccessClaims, ord *Order, opts ...RetrieveOptsFunc) (*Order, error)
	UpdateOrderByID(claims *jwtadapter.AccessClaims, id string, ord *Order, opts ...RetrieveOptsFunc) (*Order, error)
}
