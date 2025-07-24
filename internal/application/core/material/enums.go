package material

type (
	Unit         string
	CategoryEnum string
)

const (
	UnitUnknown Unit = "Unknown"
	UnitMl      Unit = "ml"
	UnitFt2     Unit = "ft²"
	UnitM       Unit = "m"
	UnitCm      Unit = "cm"
	UnitPiece   Unit = "piece"
	UnitGram    Unit = "g"
)

const (
	CategoryUnknown       CategoryEnum = "UNKNOWN"
	CategoryLeather       CategoryEnum = "LEATHER"
	CategoryThread        CategoryEnum = "THREAD"
	CategoryDye           CategoryEnum = "DYE"
	CategoryHardware      CategoryEnum = "HARDWARE"
	CategoryPaper         CategoryEnum = "PAPER"
	CategoryAdhesive      CategoryEnum = "ADHESIVE"
	CategoryFinish        CategoryEnum = "FINISH"
	CategoryLining        CategoryEnum = "LINING"
	CategoryPackaging     CategoryEnum = "PACKAGING"
	CategoryEmbellishment CategoryEnum = "EMBELLISHMENT"
	CategoryConsumable    CategoryEnum = "CONSUMABLE"
)

func (c CategoryEnum) View() string {
	return map[CategoryEnum]string{
		CategoryUnknown:       "Unknown",
		CategoryLeather:       "Leather",
		CategoryThread:        "Thread",
		CategoryDye:           "Dye",
		CategoryHardware:      "Hardware",
		CategoryPaper:         "Paper",
		CategoryAdhesive:      "Adhesive",
		CategoryFinish:        "Finish",
		CategoryLining:        "Lining",
		CategoryPackaging:     "Packaging",
		CategoryEmbellishment: "Embellishment",
		CategoryConsumable:    "Consumable",
	}[c]
}

func MaterialsCategoriesEnums() []CategoryEnum {
	return []CategoryEnum{
		CategoryUnknown, CategoryLeather,
		CategoryThread, CategoryDye,
		CategoryHardware, CategoryPaper,
		CategoryAdhesive, CategoryFinish,
		CategoryLining, CategoryPackaging,
		CategoryEmbellishment, CategoryConsumable,
	}
}
