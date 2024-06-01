package order

import (
	"time"

	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/application/core/user"
)

// TODO: add a price
// TODO: remove custom price on the order level
type Order struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	Ref    string `json:"ref" bson:"ref"`
	Number uint   `json:"number" bson:"number"`

	MerchantID   string   `json:"merchant_id" bson:"merchant"`
	CraftsmenIDs []string `json:"craftsmen_ids" bson:"craftsmen"`
	ClientID     string   `json:"client_id" bson:"client" validate:"mongodb"`

	Status      string       `json:"status" bson:"status" validate:"required,oneof=pending_confirmation confirmed in_progress pending_shipment shipping pending_payment completed canceled expired"`
	Items       []Item       `json:"items" bson:"items" validate:"required,min=1,dive,required"`
	PriceAddons []PriceAddon `json:"price_addons" bson:"price_addons" validate:"dive"`

	CustomPrice     float64          `json:"custom_price" bson:"custom_price,omitempty"`
	ReceivedAmounts []ReceivedAmount `json:"received_amounts" bson:"received_amounts" validate:"dive"`

	Timeline Timeline `json:"timeline" bson:"timeline" validate:"required"`
	Note     string   `json:"note" bson:"note,omitempty"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	Merchant  *merchant.Merchant `json:"merchant" bson:"populatedMerchant,omitempty"`
	Client    *client.Client     `json:"client" bson:"populatedClient,omitempty"`
	Craftsmen []user.User        `json:"craftsmen" bson:"populatedCraftsmen,omitempty"`
}

type Item struct {
	ID          string  `json:"id" bson:"_id,omitempty"`
	ProductID   string  `json:"product_id" bson:"product" validate:"mongodb"`
	VariantID   string  `json:"variant_id" bson:"variant" validate:"mongodb"`
	Price       float64 `json:"price" bson:"price"`
	CustomPrice float64 `json:"custom_price" bson:"custom_price"`
	Progress    string  `json:"progress" bson:"progress" validate:"required,oneof=not_started designing pending_material crafting laser_carving on_hold done"`

	Product *product.Product `json:"product" bson:"populatedProduct"`
	Variant *product.Variant `json:"variant" bson:"populatedVariant"`
}

type PriceAddon struct {
	Kind         string  `json:"kind" bson:"kind" validate:"required,oneof=fees taxes shipping discount"`
	Amount       float64 `json:"amount" bson:"amount" validate:"required"`
	IsPercentage bool    `json:"is_percentage" bson:"is_percentage"`
}

type Timeline struct {
	IssuanceDate time.Time `json:"issuance_date" bson:"issuance_date,omitempty"`
	DueDate      time.Time `json:"due_date" bson:"due_date,omitempty"`
	Deadline     time.Time `json:"deadline" bson:"deadline,omitempty"`
	DoneOn       time.Time `json:"done_on" bson:"done_on,omitempty" validate:"ltcsfield=IssuanceDate"`
	ShippedOn    time.Time `json:"shipped_on" bson:"shipped_on,omitempty" validate:"ltcsfield=IssuanceDate"`
	ResolvedOn   time.Time `json:"resolved_on" bson:"resolved_on,omitempty" validate:"ltcsfield=IssuanceDate"`
}

type ReceivedAmount struct {
	Amount float64   `json:"amount" bson:"amount"`
	Date   time.Time `json:"date" bson:"date"`
}
