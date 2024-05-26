package product

import (
	"fmt"
	"strconv"
	"time"

	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/application/core/user"
)

type CategoryEnum uint8

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

func CategoriesEnums() []CategoryEnum {
	return []CategoryEnum{
		BackPacks, Bags, Bookmarks, Bracelets,
		Cuffs, DeskPads, Folders, HairSliders,
		HandBags, Masks, PhoneCases,
		Tools,
		Wallets,
	}
}

func CategoriesStrings() []string {
	catenums := CategoriesEnums()
	categories := make([]string, len(catenums))
	for _, catenum := range CategoriesEnums() {
		categories = append(categories, catenum.String())
	}
	return categories
}

func CategoriesCodes() []string {
	catenums := CategoriesEnums()
	categories := make([]string, len(catenums))
	for _, catenum := range CategoriesEnums() {
		categories = append(categories, catenum.Code())
	}
	return categories
}

type Product struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	MerchantID  string `json:"merchant_id" bson:"merchant,omitempty"`
	CraftsmanID string `json:"craftsman_id" bson:"craftsman,omitempty"`
	Number      uint8  `json:"number" bson:"number,omitempty"`

	Name        string `json:"name" bson:"name,omitempty" validate:"required,min=3,max=255"`
	Description string `json:"description" bson:"description,omitempty"`
	Category    string `json:"category" bson:"category,omitempty" validate:"required,oneof=BKPK BAGS BKMR BRCT CUFS DKPD FLDR HSLD HNDB MASK FNCS TOLS WLET"`

	Variants []ProductVariant `json:"variants" bson:"variants" validate:"required,min=1,dive"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	Craftsman *user.User         `json:"craftsman" bson:"populatedCraftsman,omitempty"`
	Merchant  *merchant.Merchant `json:"merchant" bson:"populatedMerchant,omitempty"`
}

func (p *Product) Ref() string {
	return fmt.Sprintf("%s%s", p.Category, strconv.Itoa(int(p.Number)))
}

type ProductVariant struct {
	ID                string  `json:"id" bson:"_id,omitempty"`
	Suffix            string  `json:"suffix" bson:"suffix,omitempty" validate:"required,min=2,max=255"`
	Name              string  `json:"name" bson:"name,omitempty" validate:"required,min=3,max=255"`
	Description       string  `json:"description" bson:"description,omitempty"`
	Price             float64 `json:"price" bson:"price"`
	WholesalePrice    float64 `json:"wholesale_price" bson:"wholesale_price"`
	TimeToCraftInMins int     `json:"time_to_craft_in_mins" bson:"time_to_craft_in_mins,omitempty"`
	ProductRef        string  `json:"product_ref" bson:"product_ref,omitempty"`
}

func (p *ProductVariant) Ref() string {
	return fmt.Sprintf("%s-%s", p.ProductRef, p.Suffix)
}
