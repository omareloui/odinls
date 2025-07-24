package product

type ProductRepository interface {
	GetProducts(opts ...RetrieveOptsFunc) ([]Product, error)
	GetProductByID(id string, opts ...RetrieveOptsFunc) (*Product, error)
	GetProductByVariantID(id string, opts ...RetrieveOptsFunc) (*Product, error)
	GetProductByIDAndVariantID(id string, variantId string, options ...RetrieveOptsFunc) (*Product, error)
	CreateProduct(prod *Product, opts ...RetrieveOptsFunc) (*Product, error)
	UpdateProductByID(id string, prod *Product, opts ...RetrieveOptsFunc) (*Product, error)
}
