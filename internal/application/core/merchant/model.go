package merchant

import "time"

type Merchant struct {
	ID               string    `json:"id" bson:"_id,omitempty"`
	Name             string    `json:"name" bson:"name" conform:"trim,title" validate:"required,min=3"`
	Logo             string    `json:"logo" bson:"logo" conform:"trim" validate:"required,http_url"`
	HourlyRate       float64   `json:"hourly_rate" bson:"hourly_rate,omitempty" validate:"required,number,min=1"`
	ProfitPercentage float64   `json:"profit_percentage" bson:"profit_percentage,omitempty" validate:"required,number,min=1"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}
