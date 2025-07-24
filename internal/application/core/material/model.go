// Package material is meant for the materials, materials categories, and
// units used to produce a leather product
package material

import (
	"time"

	"github.com/omareloui/odinls/internal/application/core/supplier"
)

type Material struct {
	ID string `json:"id" bson:"_id,omitempty"`

	Name        string       `json:"name" bson:"name,omitempty" conform:"trim,title" validate:"required,min=3,max=255"`
	Description string       `json:"description" bson:"description,omitempty" conform:"trim"`
	Category    CategoryEnum `json:"category" bson:"category,omitempty" conform:"trim,upper"`

	Unit         Unit    `json:"unit" bson:"unit" conform:"trim,lower" validate:"required"`
	PricePerUnit float64 `json:"price_per_unit" bson:"price_per_unit" validate:"required,min=0"`

	QuantityOnHand  float64 `json:"quantity_on_hand" bson:"quantity_on_hand" validate:"min=0"`
	ReorderLevel    float64 `json:"reorder_level" bson:"reorder_level" validate:"min=0"`
	ReorderQuantity float64 `json:"reorder_quantity" bson:"reorder_quantity" validate:"min=0"`

	Tags []string `json:"tags" bson:"tags,omitempty"`

	SupplierID string `json:"supplier_id" bson:"supplier" validate:"required,mongodb"`

	LastPriceUpdate time.Time `json:"last_price_update" bson:"last_price_update"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	Supplier *supplier.Supplier `json:"supplier" bson:"populated_supplier,omitempty"`
}
