package domain

type Product struct {
	ID          ID
	Merchant    ID
	Craftsman   ID
	Name        string
	Description string
	Category    string // enum
	Number      int
	Variants    []ProductVariant
}

type ProductVariant struct {
	ID                   ID
	Suffix               string
	Name                 string
	Description          string
	Price                float64
	WholesalePrice       float64
	TimeToCraftInMinutes int
}
