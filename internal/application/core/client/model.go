package client

import (
	"time"

	"github.com/omareloui/odinls/internal/application/core/merchant"
)

type User struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	MerchantID string    `json:"merchant_id" bson:"merchant,omitempty" validate:"required,mongodb"`
	Name       string    `json:"name" bson:"name" validate:"required,min=3,max=255,not_blank"`
	Email      string    `json:"email" bson:"email" validate:"email"`
	Phone      string    `json:"phone" bson:"phone,omitempty"`
	Notes      string    `json:"notes" bson:"notes" validate:"not_blank"`
	Locations  []string  `json:"locations" bson:"locations" validate:"dive,required,not_blank"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`

	Merchant *merchant.Merchant `json:"merchant"`
}
