package user

import (
	"time"
)

type Name struct {
	First string `json:"first" bson:"first" validate:"required"`
	Last  string `json:"last" bson:"last" validate:"required"`
}

type User struct {
	ID        string    `json:"id" bson:"_id"`
	Name      Name      `json:"name" bson:"name" validate:"required"`
	Username  string    `json:"username" bson:"username" validate:"required"`
	Email     string    `json:"email" bson:"email" validate:"required"`
	Password  string    `json:"password" bson:"password" validate:"required"`
	Phone     string    `json:"phone" bson:"phone"`
	Role      string    `json:"role" bson:"role"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
