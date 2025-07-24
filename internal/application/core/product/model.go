package product

import (
	"fmt"
	"math"
	"time"

	"github.com/omareloui/odinls/internal/application/core/material"
)

const (
	hourlyRate = 60

	monthlyFixedCosts = 6000
	monthlyWorkHours  = 176
	hourlyFixedCosts  = monthlyFixedCosts / monthlyWorkHours

	incalculableCostsPercentage = 0.05

	retailProfitPercentage    = 1
	wholesaleProfitPercentage = 0.5
)

type Product struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	Number uint8  `json:"number" bson:"number,omitempty"`

	Name        string       `json:"name" bson:"name,omitempty" conform:"trim,title" validate:"required,min=3,max=255"`
	Description string       `json:"description" bson:"description,omitempty" conform:"trim"`
	Category    CategoryEnum `json:"category" bson:"category,omitempty" conform:"trim,upper" validate:"required"`

	Variants []Variant `json:"variants" bson:"variants" validate:"required,min=1,dive"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (p *Product) SKU() string {
	return fmt.Sprintf("%s%03d", p.Category.Code(), int(p.Number))
}

type MaterialUsage struct {
	MaterialID string  `json:"material_id" bson:"material_id"`
	Quantity   float64 `json:"quantity" bson:"quantity"`

	Material *material.Material `json:"material" bson:"populated_material"`
}

type Variant struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Suffix      string `json:"suffix" bson:"suffix,omitempty" conform:"trim,lower" validate:"required,min=2,max=255"`
	Name        string `json:"name" bson:"name,omitempty" conform:"trim,title" validate:"required,min=3,max=255"`
	Description string `json:"description" conform:"trim" bson:"description,omitempty"`

	// The options that make this variant
	Options map[string]string `json:"options" bson:"options,omitempty"`

	MaterialUsage []MaterialUsage `json:"material_usage" bson:"material_usage"`

	Price          float64 `json:"price" bson:"price"`
	WholesalePrice float64 `json:"wholesale_price" bson:"wholesale_price"`

	TimeToCraft time.Duration `json:"time_to_craft" bson:"time_to_craft,omitempty"`
	ProductSKU  string        `json:"-" bson:"-"`
}

func (v *Variant) SKU() string {
	return fmt.Sprintf("%s-%s", v.ProductSKU, v.Suffix)
}

func (v *Variant) MaterialCost() float64 {
	var sum float64 = 0
	for _, u := range v.MaterialUsage {
		if u.Material == nil {
			return 0
		}
		sum += u.Quantity * u.Material.PricePerUnit
	}
	return sum * (1 + incalculableCostsPercentage)
}

func (v *Variant) TimeCost() float64 {
	return hourlyRate * v.TimeToCraft.Hours()
}

func (v *Variant) FixedCost() float64 {
	return hourlyFixedCosts * v.TimeToCraft.Hours()
}

func (v *Variant) TotalCost() float64 {
	return v.TimeCost() + v.MaterialCost() + v.FixedCost()
}

func (v *Variant) EstPrice() float64 {
	return v.estPrice(retailProfitPercentage)
}

func (v *Variant) EstWholesalePrice() float64 {
	return v.estPrice(wholesaleProfitPercentage)
}

func (v *Variant) estPrice(profitPercentage float64) float64 {
	return math.Floor((v.TotalCost()*(1+profitPercentage))/5) * 5
}

func (v *Variant) Profit(price float64) float64 {
	return price - v.TotalCost()
}

func (v *Variant) MaxDiscountPercentage(price float64) float64 {
	profit := v.Profit(price)
	return profit / price * 100
}
