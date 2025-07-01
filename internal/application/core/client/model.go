package client

import (
	"time"

	"github.com/omareloui/odinls/internal/application/core/merchant"
)

type Client struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	MerchantID string `json:"merchant_id" bson:"merchant,omitempty"`

	Name               string      `json:"name" bson:"name" conform:"title,trim" validate:"required,min=3,max=255,not_blank"`
	Notes              string      `json:"notes" conform:"trim" bson:"notes,omitempty"`
	ContactInfo        ContactInfo `json:"contact_info" bson:"contact_info,omitempty"`
	WholesaleAsDefault bool        `json:"wholesale_as_default" bson:"wholesale_as_default" validate:"boolean"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	Merchant *merchant.Merchant `json:"merchant" bson:"populatedMerchant,omitempty"`
}

type ContactInfo struct {
	PhoneNumbers map[string]string `json:"phone_number" bson:"phone_number,omitempty" conform:"num" validate:"dive,keys,required,min=3,max=255,not_blank,endkeys,required,min=3,max=255,not_blank"`
	Emails       map[string]string `json:"emails" bson:"emails,omitempty" conform:"email" validate:"dive,keys,required,min=3,max=255,not_blank,endkeys,required,email"`
	Links        map[string]string `json:"links" bson:"links,omitempty" conform:"trim" validate:"dive,keys,required,min=3,max=255,not_blank,min=3,max=255,endkeys,required,http_url"`
	Locations    map[string]string `json:"locations" bson:"locations,omitempty" conform:"trim" validate:"dive,keys,required,min=3,max=255,not_blank,endkeys,required"`
}

func (c Client) HasContactInfo() bool {
	return len(c.ContactInfo.Emails) > 0 ||
		len(c.ContactInfo.Locations) > 0 ||
		len(c.ContactInfo.Links) > 0 ||
		len(c.ContactInfo.PhoneNumbers) > 0
}
