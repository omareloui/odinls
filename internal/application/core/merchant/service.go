package merchant

type MerchantService interface {
	FindMerchant(id string) (*Merchant, error)
	CreateMerchant(merchant *Merchant) error
}
