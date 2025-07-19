package counter

type CounterRepository interface {
	CreateCounter(*Counter) error
	GetCounter() (*Counter, error)
	AddOneToProduct(category string) (uint8, error)
	AddOneToOrder() (uint, error)
}
