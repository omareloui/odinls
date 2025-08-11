// Package supplier is meant for any leatherwork suppliers
package supplier

import "time"

type Supplier struct {
	ID string `json:"id" formfield:"-" bson:"_id"`

	Name     string `json:"name" formfield:"name" bson:"name"`
	Location string `json:"location" formfield:"location" bson:"location,omitempty"`

	Tags []string `json:"tags" formfield:"tags" bson:"tags"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
