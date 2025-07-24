package product

type CategoryEnum string

const (
	Unknown     CategoryEnum = "UNKNOWN"
	BackPacks   CategoryEnum = "BACK_PACKS"
	Bags        CategoryEnum = "BAGS"
	Bookmarks   CategoryEnum = "BOOKMARKS"
	Bracelets   CategoryEnum = "BRACELETS"
	Cuffs       CategoryEnum = "CUFFS"
	DeskPads    CategoryEnum = "DESK_PADS"
	Folders     CategoryEnum = "FOLDERS"
	HairSliders CategoryEnum = "HAIR_SLIDERS"
	HandBags    CategoryEnum = "HAND_BAGS"
	Masks       CategoryEnum = "MASKS"
	PhoneCases  CategoryEnum = "PHONE_CASES"
	Tools       CategoryEnum = "TOOLS"
	Wallets     CategoryEnum = "WALLETS"
)

func (c CategoryEnum) View() string {
	return map[CategoryEnum]string{
		Unknown:     "Unknown",
		BackPacks:   "Back Packs",
		Bags:        "Bags",
		Bookmarks:   "Bookmarks",
		Bracelets:   "Bracelets",
		Cuffs:       "Cuffs",
		DeskPads:    "Desk Pads",
		Folders:     "Folders",
		HairSliders: "Hair Sliders",
		HandBags:    "Hand Bags",
		Masks:       "Masks",
		PhoneCases:  "Phone Cases",
		Tools:       "Tools",
		Wallets:     "Wallets",
	}[c]
}

func (c CategoryEnum) Code() string {
	return map[CategoryEnum]string{
		Unknown:     "????",
		BackPacks:   "BKPK",
		Bags:        "BAGS",
		Bookmarks:   "BKMR",
		Bracelets:   "BRCT",
		Cuffs:       "CUFS",
		DeskPads:    "DKPD",
		Folders:     "FLDR",
		HairSliders: "HSLD",
		HandBags:    "HNDB",
		Masks:       "MASK",
		PhoneCases:  "FNCS",
		Tools:       "TOLS",
		Wallets:     "WLET",
	}[c]
}

func CategoriesEnums() []CategoryEnum {
	return []CategoryEnum{
		Unknown,
		BackPacks,
		Bags,
		Bookmarks,
		Bracelets,
		Cuffs,
		DeskPads,
		Folders,
		HairSliders,
		HandBags,
		Masks,
		PhoneCases,
		Tools,
		Wallets,
	}
}

func CategoriesViews() []string {
	catenums := CategoriesEnums()
	categories := make([]string, len(catenums))
	for _, catenum := range CategoriesEnums() {
		categories = append(categories, catenum.View())
	}
	return categories
}

func CategoriesCodes() []string {
	catenums := CategoriesEnums()
	categories := make([]string, len(catenums))
	for _, catenum := range catenums {
		categories = append(categories, catenum.Code())
	}
	return categories
}

