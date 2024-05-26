package counter

type CounterRepository interface {
	CreateCounter(*Counter) error

	GetCounterByID(id string) (*Counter, error)
	GetCounterByMerchantID(merchantId string) (*Counter, error)

	AddOneToProduct(merchantId, category string) (uint8, error)
	AddOneToOrder(merchantId string) (uint, error)
}
