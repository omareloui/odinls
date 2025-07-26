package counter

type CounterRepository interface {
	AddOneToProduct(category string) (uint8, error)
	AddOneToOrder() (uint, error)
}
