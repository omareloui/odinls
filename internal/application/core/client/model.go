package client

import (
	"time"

	"github.com/omareloui/odinls/internal/application/core/merchant"
)

type Client struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	MerchantID string `json:"merchant_id" bson:"merchant,omitempty"`

	Name               string      `json:"name" bson:"name" validate:"required,min=3,max=255,not_blank"`
	Notes              string      `json:"notes" bson:"notes"`
	ContactInfo        ContactInfo `json:"contact_info" bson:"contact_info" validate:"required"`
	WholesaleAsDefault bool        `json:"wholesale_as_default" bson:"wholesale_as_default" validate:"boolean"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	Merchant *merchant.Merchant `json:"merchant" bson:"populatedMerchant"`
}

type ContactInfo struct {
	PhoneNumbers map[string]string `json:"phone_number" bson:"phone_number" validate:"dive,keys,required,min=3,max=255,not_blank,endkeys,required,min=3,max=255,not_blank"`
	Emails       map[string]string `json:"emails" bson:"emails" validate:"dive,keys,required,min=3,max=255,not_blank,endkeys,required,email"`
	Links        map[string]string `json:"links" bson:"links" validate:"dive,keys,required,min=3,max=255,not_blank,min=3,max=255,endkeys,required,http_url"`
	Locations    map[string]string `json:"locations" bson:"locations" validate:"dive,keys,required,min=3,max=255,not_blank,endkeys,required"`
}
