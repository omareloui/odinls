package user

import (
	"time"
)

type RoleEnum int

const (
	NoAuthority RoleEnum = iota
	Moderator
	Admin
	SuperAdmin
)

func (r RoleEnum) String() string {
	return [...]string{
		"NO_AUTHORITY", "MODERATOR",
		"ADMIN", "SUPER_ADMIN",
	}[r]
}

func RoleFromString(role string) RoleEnum {
	switch role {
	case "NO_AUTHORITY":
		return NoAuthority
	case "MODERATOR":
		return Moderator
	case "ADMIN":
		return Admin
	case "SUPER_ADMIN":
		return SuperAdmin
	default:
		return NoAuthority
	}
}

func (r RoleEnum) IsSuperAdmin() bool {
	return r <= SuperAdmin
}

func (r RoleEnum) IsAdmin() bool {
	return r <= Admin
}

func (r RoleEnum) IsModerator() bool {
	return r <= Moderator
}

type Name struct {
	First string `json:"first" bson:"first" conform:"name" validate:"required,not_blank"`
	Last  string `json:"last" bson:"last" conform:"name" validate:"required,not_blank"`
}

type User struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	Name            Name      `json:"name" bson:"name" validate:"required"`
	Username        string    `json:"username" bson:"username" conform:"trim,lower" validate:"required,min=3,max=64,alphanum_with_underscore,not_blank"`
	Email           string    `json:"email" bson:"email" conform:"email" validate:"required,email,not_blank"`
	Password        string    `json:"password" bson:"password" validate:"required,min=8,max=64,not_blank"`
	ConfirmPassword string    `json:"-" bson:"-" validate:"eqfield=Password"`
	Phone           string    `json:"phone" bson:"phone,omitempty" conform:"num"`
	Role            RoleEnum  `json:"role" bson:"role,omitempty" validate:"required,max=3"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`

	Craftsman *Craftsman `json:"craftsman" bson:"craftsman,omitempty"`
}

type Craftsman struct {
	HourlyRate float64 `json:"hourly_rate" bson:"hourly_rate,omitempty" validate:"required,number"`
}

func (u User) IsCraftsman() bool {
	return u.Craftsman != nil
}
