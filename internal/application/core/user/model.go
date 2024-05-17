package user

import (
	"time"
)

type Name struct {
	First string `json:"first" bson:"first" validate:"required,not_blank"`
	Last  string `json:"last" bson:"last" validate:"required,not_blank"`
}

type User struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	Name            Name      `json:"name" bson:"name" validate:"required"`
	Username        string    `json:"username" bson:"username" validate:"required,min=3,max=64,alphanum_with_underscore,not_blank"`
	Email           string    `json:"email" bson:"email" validate:"required,email,not_blank"`
	Password        string    `json:"password" bson:"password" validate:"required,min=8,max=64,not_blank"`
	ConfirmPassword string    `json:"-" bson:"-" validate:"required,min=8,eqfield=Password,not_blank"`
	Phone           string    `json:"phone" bson:"phone,omitempty"`
	Role            string    `json:"role" bson:"role,omitempty"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}
