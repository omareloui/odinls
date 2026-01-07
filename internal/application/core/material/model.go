// Package material is meant for the materials, materials categories, and
// units used to produce a leather product
package material

import (
	"time"

	"github.com/omareloui/odinls/internal/application/core/supplier"
)

type Material struct {
	ID string `json:"id" bson:"_id,omitempty" formfield:"-"`

	Name        string       `json:"name" bson:"name,omitempty" formfield:"name" conform:"trim,title" validate:"required,min=3,max=255"`
	Description string       `json:"description" bson:"description,omitempty" formfield:"description" conform:"trim"`
	Category    CategoryEnum `json:"category" bson:"category,omitempty" formfield:"category" conform:"trim,upper"`

	Unit         Unit    `json:"unit" bson:"unit" conform:"trim,lower" formfield:"unit" validate:"required"`
	PricePerUnit float64 `json:"price_per_unit" bson:"price_per_unit" formfield:"price_per_unit" validate:"required,min=0"`

	QuantityOnHand  float64 `json:"quantity_on_hand" bson:"quantity_on_hand" formfield:"quantity_on_hand" validate:"min=0"`
	ReorderLevel    float64 `json:"reorder_level" bson:"reorder_level" formfield:"reorder_level" validate:"min=0"`
	ReorderQuantity float64 `json:"reorder_quantity" bson:"reorder_quantity" formfield:"reorder_quantity" validate:"min=0"`

	Tags []string `json:"tags" bson:"tags,omitempty" formfield:"tags"`

	SupplierID string `json:"supplier_id" bson:"supplier" formfield:"supplier_id" validate:"required,mongodb"`

	LastPriceUpdate time.Time `json:"last_price_update" bson:"last_price_update" formfield:"last_price_update"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	Supplier *supplier.Supplier `json:"supplier" bson:"populated_supplier,omitempty" formfield:"-"`
}
