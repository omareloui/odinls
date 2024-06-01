package merchant

import "time"

type Merchant struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name" conform:"trim,title" validate:"required,min=3"`
	Logo      string    `json:"logo" bson:"logo" conform:"trim" validate:"required,http_url"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
