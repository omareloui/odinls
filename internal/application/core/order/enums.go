package order

type (
	PriceAddonKindEnum string
	StatusEnum         string
	ItemProgressEnum   string
)

const (
	StatusPendingConfirmation StatusEnum = "PENDING_CONFIRMATION"
	StatusConfirmed           StatusEnum = "CONFIRMED"
	StatusInProgress          StatusEnum = "IN_PROGRESS"
	StatusPendingShipment     StatusEnum = "PENDING_SHIPMENT"
	StatusShipping            StatusEnum = "SHIPPING"
	StatusCompleted           StatusEnum = "COMPLETED"
	StatusCanceled            StatusEnum = "CANCELED"
	StatusExpired             StatusEnum = "EXPIRED"
)

func (s StatusEnum) View() string {
	v := map[StatusEnum]string{
		StatusPendingConfirmation: "Pending Confirmation",
		StatusConfirmed:           "Confirmed",
		StatusInProgress:          "In_progress",
		StatusPendingShipment:     "Pending Shipment",
		StatusShipping:            "Shipping",
		StatusCompleted:           "Completed",
		StatusCanceled:            "Canceled",
		StatusExpired:             "Expired",
	}[s]
	if v == "" {
		return StatusPendingConfirmation.View()
	}
	return v
}

func StatusesEnums() []StatusEnum {
	return []StatusEnum{
		StatusPendingConfirmation, StatusConfirmed,
		StatusInProgress, StatusPendingShipment, StatusShipping,
		StatusCompleted, StatusCanceled, StatusExpired,
	}
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
	ItemProgressNotStarted       ItemProgressEnum = "NOT_STARTED"
	ItemProgressDesigning        ItemProgressEnum = "DESIGNING"
	ItemProgressPendingMaterials ItemProgressEnum = "PENDING_MATERIALS"
	ItemProgressCrafting         ItemProgressEnum = "CRAFTING"
	ItemProgressLaserCarving     ItemProgressEnum = "LASER_CARVING"
	ItemProgressOnHold           ItemProgressEnum = "ON_HOLD"
	ItemProgressDone             ItemProgressEnum = "DONE"
)

func (s ItemProgressEnum) View() string {
	v := map[ItemProgressEnum]string{
		ItemProgressNotStarted:       "Not Started",
		ItemProgressDesigning:        "Designing",
		ItemProgressPendingMaterials: "Pending Materials",
		ItemProgressCrafting:         "Crafting",
		ItemProgressLaserCarving:     "Laser Carving",
		ItemProgressOnHold:           "On Hold",
		ItemProgressDone:             "Done",
	}[s]
	if v == "" {
		return ItemProgressNotStarted.View()
	}
	return v
}

func ItemsProgressEnums() []ItemProgressEnum {
	return []ItemProgressEnum{
		ItemProgressNotStarted, ItemProgressDesigning,
		ItemProgressPendingMaterials, ItemProgressCrafting,
		ItemProgressLaserCarving, ItemProgressOnHold,
		ItemProgressDone,
	}
}

func ItemsProgressViews() []string {
	itemsProgressEnums := ItemsProgressEnums()
	itemsProgress := make([]string, len(itemsProgressEnums))
	for _, enum := range itemsProgressEnums {
		itemsProgress = append(itemsProgress, enum.View())
	}
	return itemsProgress
}

const (
	PriceAddonKindFees     PriceAddonKindEnum = "FEES"
	PriceAddonKindTaxes    PriceAddonKindEnum = "TAXES"
	PriceAddonKindShipping PriceAddonKindEnum = "SHIPPING"
	PriceAddonKindDiscount PriceAddonKindEnum = "DISCOUNT"
)

func (p PriceAddonKindEnum) View() string {
	v := map[PriceAddonKindEnum]string{
		PriceAddonKindFees:     "Fees",
		PriceAddonKindTaxes:    "Taxes",
		PriceAddonKindShipping: "Shipping",
		PriceAddonKindDiscount: "Discount",
	}[p]
	if v == "" {
		return PriceAddonKindFees.View()
	}
	return v
}

func PriceAddonKindEnums() []PriceAddonKindEnum {
	return []PriceAddonKindEnum{
		PriceAddonKindFees,
		PriceAddonKindTaxes,
		PriceAddonKindShipping,
		PriceAddonKindDiscount,
	}
}

func PriceAddonsViews() []string {
	priceAddonKindEnums := PriceAddonKindEnums()
	priceAddons := make([]string, len(priceAddonKindEnums))
	for _, enum := range priceAddonKindEnums {
		priceAddons = append(priceAddons, enum.View())
	}
	return priceAddons
}
