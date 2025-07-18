package counter

type CounterRepository interface {
	CreateCounter(*Counter) error
	GetCounterByID(id string) (*Counter, error)
	AddOneToProduct(category string) (uint8, error)
	AddOneToOrder() (uint, error)
}
