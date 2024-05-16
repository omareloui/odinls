package merchant

type MerchantRepository interface {
	FindMerchant(id string) (*Merchant, error)
	CreateMerchant(merchant *Merchant) error
}
