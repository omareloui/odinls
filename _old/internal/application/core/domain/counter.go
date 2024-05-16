package domain

type Counter struct {
	ID            ID
	Merchant      ID
	OrderNumber   int16
	ProductsCodes string
}
