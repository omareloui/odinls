package order

import (
	"time"

	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/application/core/user"
)

type Order struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	Ref    string `json:"ref" bson:"ref"`
	Number uint   `json:"number" bson:"number"`

	MerchantID   string   `json:"merchant_id" bson:"merchant"`
	CraftsmenIDs []string `json:"craftsmen_ids" bson:"craftsmen"`
	ClientID     string   `json:"client_id" bson:"client" validate:"required,mongodb"`

	Status      string       `json:"status" bson:"status" validate:"required,oneof=pending_confirmation confirmed in_progress pending_shipment shipping pending_payment completed canceled expired"`
	Items       []Item       `json:"items" bson:"items" validate:"required,min=1,dive,required"`
	PriceAddons []PriceAddon `json:"price_addons" bson:"price_addons" validate:"dive,required"`

	ReceivedAmounts []ReceivedAmount `json:"received_amounts" bson:"received_amounts" validate:"dive"`

	Timeline Timeline `json:"timeline" bson:"timeline" validate:"required"`
	Note     string   `json:"note" bson:"note,omitempty"`

	Subtotal float64 `json:"subtotal" bson:"subtotal,omitempty"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	Merchant  *merchant.Merchant `json:"merchant" bson:"populatedMerchant,omitempty"`
	Client    *client.Client     `json:"client" bson:"populatedClient,omitempty"`
	Craftsmen []user.User        `json:"craftsmen" bson:"populatedCraftsmen,omitempty"`
}

type Item struct {
	ID          string  `json:"id" bson:"_id,omitempty"`
	ProductID   string  `json:"product_id" bson:"product" validate:"required,mongodb"`
	VariantID   string  `json:"variant_id" bson:"variant" validate:"required,mongodb"`
	Price       float64 `json:"price" bson:"price"`
	CustomPrice float64 `json:"custom_price" bson:"custom_price" validate:"gte=0"`
	Progress    string  `json:"progress" bson:"progress" validate:"required,oneof=not_started designing pending_material crafting laser_carving on_hold done"`

	Product *product.Product `json:"product" bson:"populatedProduct"`
	Variant *product.Variant `json:"variant" bson:"populatedVariant"`
}

type PriceAddon struct {
	Kind         string  `json:"kind" bson:"kind" validate:"required,oneof=fees taxes shipping discount"`
	Amount       float64 `json:"amount" bson:"amount" validate:"required,gte=1"`
	IsPercentage bool    `json:"is_percentage" bson:"is_percentage"`
}

type Timeline struct {
	IssuanceDate time.Time `json:"issuance_date" bson:"issuance_date,omitempty" validate:""`
	DueDate      time.Time `json:"due_date" bson:"due_date,omitempty" validate:""`
	Deadline     time.Time `json:"deadline" bson:"deadline,omitempty" validate:""`
	DoneOn       time.Time `json:"done_on" bson:"done_on,omitempty" validate:"ltcsfield=IssuanceDate"`
	ShippedOn    time.Time `json:"shipped_on" bson:"shipped_on,omitempty" validate:"ltcsfield=IssuanceDate"`
	ResolvedOn   time.Time `json:"resolved_on" bson:"resolved_on,omitempty" validate:"ltcsfield=IssuanceDate"`
}

type ReceivedAmount struct {
	Amount float64   `json:"amount" bson:"amount" validate:"required,lte=0"`
	Date   time.Time `json:"date" bson:"date" validate:"required"`
}

func (o *Order) CalcSubtotal() float64 {
	var total float64

	var itemsTotal float64

	var discountsPercentage float64
	var taxesPercentage float64
	var feesPercentage float64
	var shippingPercentage float64

	for _, item := range o.Items {
		if item.CustomPrice > 0 {
			itemsTotal += item.CustomPrice
		} else {
			itemsTotal += item.Price
		}
	}

	for _, addon := range o.PriceAddons {
		if !addon.IsPercentage {
			if addon.Kind == PriceAddonKindDiscount.String() {
				total -= addon.Amount
			} else {
				total += addon.Amount
			}
			continue
		}

		switch addon.Kind {
		case PriceAddonKindTaxes.String():
			taxesPercentage = addon.Amount / 100
		case PriceAddonKindFees.String():
			feesPercentage = addon.Amount / 100
		case PriceAddonKindShipping.String():
			shippingPercentage = addon.Amount / 100
		case PriceAddonKindDiscount.String():
			discountsPercentage = addon.Amount / 100
		}
	}

	total += total * (1 + feesPercentage)
	total += total * (1 + shippingPercentage)
	total -= total * (1 + discountsPercentage)
	total += total * (1 + taxesPercentage)

	return total
}

func (o *Order) RemainingAmount() float64 {
	var paid float64
	for _, received := range o.ReceivedAmounts {
		paid += received.Amount
	}
	return o.CalcSubtotal() - paid
}
