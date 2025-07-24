package material

type (
	Unit                 string
	MaterialCategoryEnum string
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
	MaterialCategoryUnknown       MaterialCategoryEnum = "UNKNOWN"
	MaterialCategoryLeather       MaterialCategoryEnum = "LEATHER"
	MaterialCategoryThread        MaterialCategoryEnum = "THREAD"
	MaterialCategoryDye           MaterialCategoryEnum = "DYE"
	MaterialCategoryHardware      MaterialCategoryEnum = "HARDWARE"
	MaterialCategoryPaper         MaterialCategoryEnum = "PAPER"
	MaterialCategoryAdhesive      MaterialCategoryEnum = "ADHESIVE"
	MaterialCategoryFinish        MaterialCategoryEnum = "FINISH"
	MaterialCategoryLining        MaterialCategoryEnum = "LINING"
	MaterialCategoryPackaging     MaterialCategoryEnum = "PACKAGING"
	MaterialCategoryEmbellishment MaterialCategoryEnum = "EMBELLISHMENT"
	MaterialCategoryConsumable    MaterialCategoryEnum = "CONSUMABLE"
)

func (c MaterialCategoryEnum) View() string {
	return map[MaterialCategoryEnum]string{
		MaterialCategoryUnknown:       "Unknown",
		MaterialCategoryLeather:       "Leather",
		MaterialCategoryThread:        "Thread",
		MaterialCategoryDye:           "Dye",
		MaterialCategoryHardware:      "Hardware",
		MaterialCategoryPaper:         "Paper",
		MaterialCategoryAdhesive:      "Adhesive",
		MaterialCategoryFinish:        "Finish",
		MaterialCategoryLining:        "Lining",
		MaterialCategoryPackaging:     "Packaging",
		MaterialCategoryEmbellishment: "Embellishment",
		MaterialCategoryConsumable:    "Consumable",
	}[c]
}

func MaterialsCategoriesEnums() []MaterialCategoryEnum {
	return []MaterialCategoryEnum{
		MaterialCategoryUnknown, MaterialCategoryLeather,
		MaterialCategoryThread, MaterialCategoryDye,
		MaterialCategoryHardware, MaterialCategoryPaper,
		MaterialCategoryAdhesive, MaterialCategoryFinish,
		MaterialCategoryLining, MaterialCategoryPackaging,
		MaterialCategoryEmbellishment, MaterialCategoryConsumable,
	}
}
