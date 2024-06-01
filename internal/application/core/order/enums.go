package order

type (
	PriceAddonKindEnum uint8
	StatusEnum         uint8
	ItemProgressEnum   uint8
)

const (
	StatusPendingConfirmation StatusEnum = iota
	StatusConfirmed
	StatusInProgress
	StatusPendingShipment
	StatusShipping
	StatusPendingPayment
	StatusCompleted
	StatusCanceled
	StatusExpired
)

func (s *StatusEnum) String() string {
	return [...]string{
		"pending_confirmation", "confirmed",
		"in_progress", "pending_shipment",
		"shipping", "pending_payment",
		"completed", "canceled", "expired",
	}[*s]
}

func (s *StatusEnum) View() string {
	return [...]string{
		"Pending Confirmation", "Confirmed",
		"In Progress", "Pending Shipment",
		"Shipping", "Pending Payment",
		"Completed", "Canceled", "Expired",
	}[*s]
}

func StatusesEnums() []StatusEnum {
	return []StatusEnum{
		StatusPendingConfirmation, StatusConfirmed, StatusInProgress,
		StatusPendingShipment, StatusShipping, StatusPendingPayment,
		StatusCompleted, StatusCanceled, StatusExpired,
	}
}

func StatusesStrings() []string {
	statusesEnums := StatusesEnums()
	statuses := make([]string, len(statusesEnums))
	for _, enum := range statusesEnums {
		statuses = append(statuses, enum.String())
	}
	return statuses
}

func StatusesViews() []string {
	statusesEnums := StatusesEnums()
	statuses := make([]string, len(statusesEnums))
	for _, enum := range statusesEnums {
		statuses = append(statuses, enum.View())
	}
	return statuses
}

const (
	ItemProgressNotStarted ItemProgressEnum = iota
	ItemProgressDesigning
	ItemProgressPendingMaterials
	ItemProgressCrafting
	ItemProgressLaserCarving
	ItemProgressOnHold
	ItemProgressDone
)

func (s *ItemProgressEnum) View() string {
	return [...]string{
		"Not Started", "Designing", "Pending Material",
		"Crafting", "Laser Carving", "On Hold", "Done",
	}[*s]
}

func (s *ItemProgressEnum) String() string {
	return [...]string{
		"not_started", "designing", "pending_material",
		"crafting", "laser_carving", "on_hold", "done",
	}[*s]
}

func ItemsProgressEnums() []ItemProgressEnum {
	return []ItemProgressEnum{
		ItemProgressNotStarted, ItemProgressDesigning, ItemProgressPendingMaterials,
		ItemProgressCrafting, ItemProgressLaserCarving, ItemProgressOnHold, ItemProgressDone,
	}
}

func ItemsProgressStrings() []string {
	itemProgressEnums := ItemsProgressEnums()
	progress := make([]string, len(itemProgressEnums))
	for _, enum := range itemProgressEnums {
		progress = append(progress, enum.String())
	}
	return progress
}

func ItemsProgressViews() []string {
	itemProgressEnums := ItemsProgressEnums()
	progress := make([]string, len(itemProgressEnums))
	for _, enum := range itemProgressEnums {
		progress = append(progress, enum.View())
	}
	return progress
}

const (
	PriceAddonKindFees PriceAddonKindEnum = iota
	PriceAddonKindTaxes
	PriceAddonKindShipping
	PriceAddonKindDiscount
)

func (p PriceAddonKindEnum) String() string {
	return [...]string{
		"fees", "taxes",
		"shipping", "discount",
	}[p]
}

func (p PriceAddonKindEnum) View() string {
	return [...]string{
		"Fees", "Taxes",
		"Shipping", "Discount",
	}[p]
}

func PriceAddonKindEnums() []PriceAddonKindEnum {
	return []PriceAddonKindEnum{
		PriceAddonKindFees, PriceAddonKindTaxes,
		PriceAddonKindShipping, PriceAddonKindDiscount,
	}
}

func PriceAddonsStrings() []string {
	enums := PriceAddonKindEnums()
	kinds := make([]string, len(enums))
	for _, enum := range enums {
		kinds = append(kinds, enum.String())
	}
	return kinds
}

func PriceAddonsViews() []string {
	enums := PriceAddonKindEnums()
	kinds := make([]string, len(enums))
	for _, enum := range enums {
		kinds = append(kinds, enum.View())
	}
	return kinds
}
