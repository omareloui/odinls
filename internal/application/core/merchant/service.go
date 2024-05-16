package merchant

type MerchantService interface {
	GetMerchants() ([]Merchant, error)
	FindMerchant(id string) (*Merchant, error)
	CreateMerchant(merchant *Merchant) error
}
