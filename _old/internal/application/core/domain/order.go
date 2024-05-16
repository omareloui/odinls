package domain

import "time"

type Order struct {
	ID  ID
	Ref string

	Merchant   ID
	AssignedTo []ID
	Client     ID
	Items      []Item
	Number     int
	Status     string // enum

	CustomPrice    float64
	ReceivedAmount float64
	PriceAddons    []PriceAddon

	Timeline Timeline
	Notes    string
}

type Timeline struct {
	IssuanceDate time.Time
	DueDate      time.Time
	Deadline     time.Time
	DoneOn       time.Time
	ShippedOn    time.Time
	ResolvedOn   time.Time
}

type PriceAddon struct {
	Kind         string // enum
	Amount       float64
	IsPercentage bool
}

type Item struct {
	ID          ID
	Product     ID
	Variant     ID
	Price       float64
	CustomPrice float64
	Status      string // enum
	Quantity    int
}
