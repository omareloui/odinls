package merchant

type MerchantService interface {
	GetMerchants() ([]Merchant, error)
	GetMerchantByID(id string) (*Merchant, error)
	UpdateMerchantByID(id string, merchant *Merchant) error
	CreateMerchant(merchant *Merchant) error
}
