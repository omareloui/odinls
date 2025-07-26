package user

import (
	"time"
)

type Name struct {
	First string `json:"first" bson:"first" conform:"name" validate:"required,not_blank"`
	Last  string `json:"last" bson:"last" conform:"name" validate:"required,not_blank"`
}

func (n Name) FullName() string {
	if n.First == "" && n.Last == "" {
		return ""
	}
	if n.First == "" {
		return n.Last
	}
	if n.Last == "" {
		return n.First
	}
	return n.First + " " + n.Last
}

type User struct {
	ID              string   `json:"id" bson:"_id,omitempty"`
	Name            Name     `json:"name" bson:"name" validate:"required"`
	Username        string   `json:"username" bson:"username" conform:"trim,lower" validate:"required,min=3,max=64,alphanum_with_underscore,not_blank"`
	Email           string   `json:"email" bson:"email" conform:"email" validate:"required,email,not_blank"`
	Password        string   `json:"password" bson:"password" validate:"required,min=8,max=64,not_blank"`
	ConfirmPassword string   `json:"-" bson:"-" validate:"eqfield=Password"`
	Phone           string   `json:"phone" bson:"phone,omitempty" conform:"num"`
	Role            RoleEnum `json:"role" bson:"role"`

	Picture string `json:"picture_url" bson:"picture_url,omitempty" validate:"omitempty,url"`

	OAuthID       string        `json:"oauth_id,omitzero" bson:"oauth_id,omitempty" validate:"omitempty"`
	OAuthProvider OAuthProvider `json:"oauth_provider,omitzero" bson:"oauth_provider,omitempty" validate:"omitempty"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	Craftsman *Craftsman `json:"craftsman" bson:"craftsman,omitempty"`
}

type Craftsman struct {
	HourlyRate float64 `json:"hourly_rate" bson:"hourly_rate,omitempty" validate:"required,number"`
}

func (u User) IsCraftsman() bool {
	return u.Craftsman != nil
}
