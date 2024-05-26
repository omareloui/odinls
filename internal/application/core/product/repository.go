package product

type ProductRepository interface {
	GetProducts(...RetrieveOptsFunc) ([]Product, error)
	GetProductByID(string, ...RetrieveOptsFunc) (*Product, error)
	GetProductsByMerchantID(string, ...RetrieveOptsFunc) ([]Product, error)
	CreateProduct(*Product, ...RetrieveOptsFunc) error
	UpdateProductByID(string, *Product, ...RetrieveOptsFunc) error
}
