package order

import (
	"fmt"
	"os/user"
	"slices"
	"time"

	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/application/core/product"
)

const splitRefOnIdx = 4

type Order struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	Ref    string `json:"ref" bson:"ref,omitempty"`
	Number uint   `json:"number" bson:"number,omitempty"`

	ClientID      string `json:"client_id" bson:"client" validate:"required,mongodb"`
	CustomerName  string `json:"customer_name,omitzero" bson:"customer_name,omitempty"`
	CustomerEmail string `json:"customer_email,omitzero" bson:"customer_email,omitempty"`
	CustomerPhone string `json:"customer_phone,omitzero" bson:"customer_phone,omitempty"`

	Status      StatusEnum   `json:"status" bson:"status" validate:"required"`
	Items       []Item       `json:"items" bson:"items" validate:"required,min=1,dive"`
	PriceAddons []PriceAddon `json:"price_addons" bson:"price_addons,omitempty" validate:"dive"`

	ReceivedAmounts []ReceivedAmount `json:"received_amounts" bson:"received_amounts,omitempty" validate:"dive"`

	Note string `json:"note" bson:"note,omitempty"`

	Timeline Timeline `json:"timeline" bson:"timeline"`

	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at,omitempty"`

	Client *client.Client `json:"client" bson:"populated_client,omitempty"`
}

type PriceAddon struct {
	Kind         PriceAddonKindEnum `json:"kind" bson:"kind" validate:"required"`
	Amount       float64            `json:"amount" bson:"amount" validate:"required,gte=1"`
	IsPercentage bool               `json:"is_percentage" bson:"is_percentage"`
}

type Timeline struct {
	IssuanceDate  time.Time `json:"issuance_date" bson:"issuance_date" validate:"required,gte=now"`
	ScheduledDate time.Time `json:"scheduled_date,omitzero" bson:"scheduled_date,omitempty" validate:"omitempty,gtfield=IssuanceDate"`
	DoneOn        time.Time `json:"done_on,omitzero" bson:"done_on,omitempty" validate:"omitempty,gtfield=IssuanceDate"`
	ShippedOn     time.Time `json:"shipped_on,omitzero" bson:"shipped_on,omitempty" validate:"omitempty,gtfield=IssuanceDate"`
	ResolvedOn    time.Time `json:"resolved_on,omitzero" bson:"resolved_on,omitempty" validate:"omitempty,gtfield=IssuanceDate"`
	DueDate       time.Time `json:"due_date,omitzero" bson:"due_date,omitempty" validate:"omitempty,gtfield=IssuanceDate"`
}

type ReceivedAmount struct {
	Amount float64   `json:"amount" bson:"amount" validate:"required,lte=0"`
	Date   time.Time `json:"date" bson:"date" validate:"required"`
}

type Item struct {
	ID       string           `json:"id" bson:"_id,omitempty" validate:"omitempty,mongodb"`
	Progress ItemProgressEnum `json:"progress" bson:"progress" validate:"omitempty"`

	CraftsmanID string `json:"craftsman_id,omitzero" bson:"craftsman,omitempty" validate:"omitempty,mongodb"`

	CustomUnitPrice float64 `json:"custom_price" bson:"custom_price" validate:"gte=0"`
	Quantity        uint16  `json:"quantity" bson:"quantity"`

	Snapshot ItemSnapshot `json:"snapshot" bson:"snapshot,omitempty"`

	Product   *product.Product `json:"product" bson:"populated_product"`
	Craftsman *user.User       `json:"craftsman" bson:"populated_craftsman,omitempty"`
}

type ItemSnapshot struct {
	ProductID string `json:"product_id" bson:"product,omitempty"`

	ProductName string               `json:"name" bson:"name,omitempty" conform:"trim,title" validate:"required,min=3,max=255"`
	Category    product.CategoryEnum `json:"category" bson:"category,omitempty" conform:"trim,upper" validate:"required"`

	VariantID   string            `json:"variant_id" bson:"variant_id,omitempty" validate:"required,mongodb"`
	VariantName string            `json:"variant_name" bson:"variant_name,omitempty" conform:"trim,title" validate:"required,min=3,max=255"`
	SKU         string            `json:"sku" bson:"sku,omitempty"`
	Options     map[string]string `json:"options" bson:"options,omitempty"`

	Price float64 `json:"price" bson:"price" validate:"required,gte=0"`

	TimeToCraft time.Duration `json:"time_to_craft" bson:"time_to_craft,omitempty"`
}

func (o *Order) RefView() string {
	if o.Ref == "" {
		return ""
	}
	return fmt.Sprintf("%s-%s", o.Ref[:splitRefOnIdx], o.Ref[splitRefOnIdx:])
}

func (o *Order) Subtotal() float64 {
	var sum float64
	for _, item := range o.Items {
		sum += item.TotalPrice()
	}
	return sum
}

func (o *Order) TimeToCraft() time.Duration {
	var duration time.Duration
	for i, item := range o.Items {
		if item.Snapshot.TimeToCraft > 0 {
			duration += item.Snapshot.TimeToCraft * time.Duration(item.Quantity)
		} else if item.Product != nil && len(item.Product.Variants) > i {
			vIdx := slices.IndexFunc(item.Product.Variants, func(v product.Variant) bool {
				return v.ID == item.Snapshot.VariantID
			})
			if vIdx < 0 || vIdx >= len(item.Product.Variants) {
				continue
			}
			duration += item.Product.Variants[vIdx].TimeToCraft * time.Duration(item.Quantity)
		}
	}
	return duration
}

func (o *Order) TotalPrice() float64 {
	var total float64

	var subtotal float64 = o.Subtotal()

	var discountsPercentage float64
	var taxesPercentage float64
	var feesPercentage float64
	var shippingPercentage float64

	for _, addon := range o.PriceAddons {
		if !addon.IsPercentage {
			if addon.Kind == PriceAddonKindDiscount {
				total -= addon.Amount
			} else {
				total += addon.Amount
			}
			continue
		}

		switch addon.Kind {
		case PriceAddonKindTaxes:
			taxesPercentage = addon.Amount / 100
		case PriceAddonKindFees:
			feesPercentage = addon.Amount / 100
		case PriceAddonKindShipping:
			shippingPercentage = addon.Amount / 100
		case PriceAddonKindDiscount:
			discountsPercentage = addon.Amount / 100
		}
	}

	total += subtotal
	total += subtotal * feesPercentage
	total += subtotal * shippingPercentage
	total -= subtotal * discountsPercentage

	total += total * taxesPercentage

	return total
}

func (o *Order) RemainingAmount() float64 {
	var paid float64
	for _, received := range o.ReceivedAmounts {
		paid += received.Amount
	}
	return o.TotalPrice() - paid
}

func (o *Order) NotFullyPaid() bool {
	return o.RemainingAmount() > 0
}

func (i *Item) TotalPrice() float64 {
	if i.CustomUnitPrice > 0 {
		return i.CustomUnitPrice * float64(i.Quantity)
	}
	return i.Snapshot.Price * float64(i.Quantity)
}
