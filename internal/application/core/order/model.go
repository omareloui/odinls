package order

import (
	"fmt"
	"time"

	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/application/core/user"
)

const splitRefOnIdx = 4

type Order struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	Ref    string `json:"ref" bson:"ref,omitempty"`
	Number uint   `json:"number" bson:"number,omitempty"`

	MerchantID   string   `json:"merchant_id" bson:"merchant,omitempty" validate:"omitempty,mongodb"`
	CraftsmenIDs []string `json:"craftsmen_ids" bson:"craftsmen,omitempty" validate:"omitempty,mongodb"`
	ClientID     string   `json:"client_id" bson:"client" validate:"required,mongodb"`

	Status      string       `json:"status" bson:"status" validate:"required,oneof=pending_confirmation confirmed in_progress pending_shipment shipping pending_payment completed canceled expired"`
	Items       []Item       `json:"items" bson:"items" validate:"required,min=1,dive"`
	PriceAddons []PriceAddon `json:"price_addons" bson:"price_addons,omitempty" validate:"dive"`

	ReceivedAmounts []ReceivedAmount `json:"received_amounts" bson:"received_amounts,omitempty" validate:"dive"`

	Timeline Timeline `json:"timeline" bson:"timeline"`
	Note     string   `json:"note" bson:"note,omitempty"`

	Subtotal float64 `json:"subtotal" bson:"subtotal,omitempty"`

	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at,omitempty"`

	Merchant  *merchant.Merchant `json:"merchant" bson:"populatedMerchant,omitempty"`
	Client    *client.Client     `json:"client" bson:"populatedClient,omitempty"`
	Craftsmen []user.User        `json:"craftsmen" bson:"populatedCraftsmen,omitempty"`
}

func (o *Order) RefView() string {
	if o.Ref == "" {
		return ""
	}
	return fmt.Sprintf("%s-%s", o.Ref[:splitRefOnIdx], o.Ref[splitRefOnIdx:])
}

type Item struct {
	ID          string  `json:"id" bson:"_id,omitempty" validate:"omitempty,mongodb"`
	ProductID   string  `json:"product_id" bson:"product" validate:"required,mongodb"`
	VariantID   string  `json:"variant_id" bson:"variant" validate:"required,mongodb"`
	Price       float64 `json:"price" bson:"price"`
	CustomPrice float64 `json:"custom_price" bson:"custom_price" validate:"gte=0"`
	Progress    string  `json:"progress" bson:"progress" validate:"omitempty,oneof=not_started designing pending_material crafting laser_carving on_hold done"`

	Product *product.Product `json:"product" bson:"populatedProduct"`
	Variant *product.Variant `json:"variant" bson:"populatedVariant"`
}

type PriceAddon struct {
	Kind         string  `json:"kind" bson:"kind" validate:"required,oneof=fees taxes shipping discount"`
	Amount       float64 `json:"amount" bson:"amount" validate:"required,gte=1"`
	IsPercentage bool    `json:"is_percentage" bson:"is_percentage"`
}

type Timeline struct {
	IssuanceDate time.Time `json:"issuance_date" bson:"issuance_date,omitempty" validate:"omitempty"`
	DueDate      time.Time `json:"due_date" bson:"due_date,omitempty" validate:"omitempty"`
	Deadline     time.Time `json:"deadline" bson:"deadline,omitempty" validate:"omitempty"`
	DoneOn       time.Time `json:"done_on" bson:"done_on,omitempty" validate:"omitempty,gtfield=IssuanceDate"`
	ShippedOn    time.Time `json:"shipped_on" bson:"shipped_on,omitempty" validate:"omitempty,gtfield=IssuanceDate"`
	ResolvedOn   time.Time `json:"resolved_on" bson:"resolved_on,omitempty" validate:"omitempty,gtfield=IssuanceDate"`
}

type ReceivedAmount struct {
	Amount float64   `json:"amount" bson:"amount" validate:"required,lte=0"`
	Date   time.Time `json:"date" bson:"date" validate:"required"`
}

func (o *Order) calcSubtotal() float64 {
	var total float64

	var itemsTotal float64

	for _, item := range o.Items {
		if item.CustomPrice > 0 {
			itemsTotal += item.CustomPrice
		} else {
			itemsTotal += item.Price
		}
	}

	var discountsPercentage float64
	var taxesPercentage float64
	var feesPercentage float64
	var shippingPercentage float64

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

	total += itemsTotal
	total += itemsTotal * feesPercentage
	total += itemsTotal * shippingPercentage
	total -= itemsTotal * discountsPercentage

	total += total * taxesPercentage

	return total
}

func (o *Order) RemainingAmount() float64 {
	var paid float64
	for _, received := range o.ReceivedAmounts {
		paid += received.Amount
	}
	return o.calcSubtotal() - paid
}
