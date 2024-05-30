package product

type ProductRepository interface {
	GetProducts(opts ...RetrieveOptsFunc) ([]Product, error)
	GetProductByID(id string, opts ...RetrieveOptsFunc) (*Product, error)
	GetProductsByMerchantID(id string, opts ...RetrieveOptsFunc) ([]Product, error)
	CreateProduct(prod *Product, opts ...RetrieveOptsFunc) error
	UpdateProductByID(id string, prod *Product, opts ...RetrieveOptsFunc) error
}
