package order

type OrderRepository interface {
	GetOrders(opts ...RetrieveOptsFunc) ([]Order, error)
	GetOrderByID(id string, opts ...RetrieveOptsFunc) (*Order, error)
	GetOrdersByMerchantID(id string, opts ...RetrieveOptsFunc) ([]Order, error)
	CreateOrder(ord *Order, opts ...RetrieveOptsFunc) error
	UpdateOrderByID(id string, ord *Order, opts ...RetrieveOptsFunc) error
}
