package supplier

type SupplierRepository interface {
	GetSuppliers() ([]Supplier, error)
	GetSupplierByID(id string) (*Supplier, error)
	CreateSupplier(supplier *Supplier) (*Supplier, error)
	UpdateSupplierByID(id string, supplier *Supplier) (*Supplier, error)
}
