package product

import (
	"fmt"
	"strconv"
)

type CategoryEnum int

const (
	BackPacks CategoryEnum = iota
	Bags
	Bookmarks
	Bracelets
	Cuffs
	DeskPads
	Folders
	HairSliders
	HandBags
	Masks
	PhoneCases
	Tools
	Wallets
)

func (c *CategoryEnum) String() string {
	return [...]string{
		"Back Packs", "Bags", "Bookmarks", "Bracelets",
		"Cuffs", "Desk Pads", "Folders", "Hair Sliders",
		"Hand Bags", "Masks", "Phone Cases", "Tools",
		"Wallets",
	}[*c]
}

func (c *CategoryEnum) Code() string {
	return [...]string{
		"BKPK", "BAGS", "BKMR", "BRCT",
		"CUFS", "DKPD", "FLDR", "HSLD",
		"HNDB", "MASK", "FNCS", "TOLS",
		"WLET",
	}[*c]
}

type Product struct {
	ID          string
	MerchantID  string
	CraftsmanID string
	Number      int

	Name        string
	Description string
	Category    string // enum

	Variants []ProductVariant
}

func (p *Product) Ref() string {
	return fmt.Sprintf("%s%s", p.Category, strconv.Itoa(p.Number))
}

type ProductVariant struct {
	ID                string
	Suffix            string // TODO(research): replace this with a auto generated field
	Name              string
	Description       string
	Price             float64
	WholesalePrice    float64
	TimeToCraftInMins int

	ProductRef string
}

func (p *ProductVariant) Ref() string {
	return fmt.Sprintf("%s-%s", p.ProductRef, p.Suffix)
}
