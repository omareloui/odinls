package counter

import (
	"time"
)

type ProductsCodes map[string]uint8

type Counter struct {
	ID            string        `json:"id" bson:"_id,omitempty"`
	OrdersNumber  uint          `json:"orders_number" bson:"orders_number,omitempty"`
	ProductsCodes ProductsCodes `json:"products_codes" bson:"products_codes,omitempty" validate:"required,dive,keys,required,min=3,max=255,not_blank,endkeys,required"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
