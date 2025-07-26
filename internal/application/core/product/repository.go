package product

type ProductRepository interface {
	GetProducts(opts ...RetrieveOptsFunc) ([]Product, error)
	GetProductByID(id string, opts ...RetrieveOptsFunc) (*Product, error)
	GetProductByVariantID(id string, opts ...RetrieveOptsFunc) (*Product, error)
	CreateProduct(prod *Product, opts ...RetrieveOptsFunc) (*Product, error)
	UpdateProductByID(id string, prod *Product, opts ...RetrieveOptsFunc) (*Product, error)
}
