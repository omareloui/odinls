package role

import (
	"time"
)

type RoleEnum int

const (
	OPAdmin RoleEnum = iota
	SuperAdmin
	Admin
	Moderator
	NoAuthority
)

func (r RoleEnum) String() string {
	return [...]string{
		"OP_ADMIN", "SUPER_ADMIN",
		"ADMIN", "MODERATOR",
		"NO_AUTHORITY",
	}[r]
}

type Role struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name" validate:"required,oneof=OP_ADMIN SUPER_ADMIN ADMIN MODERATOR NO_AUTHORITY"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (r Role) IsOPAdmin() bool {
	return r.Name == OPAdmin.String()
}

func (r Role) IsSuperAdmin() bool {
	n := r.Name
	return n == OPAdmin.String() || n == SuperAdmin.String()
}

func (r Role) IsAdmin() bool {
	n := r.Name
	return n == OPAdmin.String() || n == SuperAdmin.String() || n == Admin.String()
}

func (r Role) IsModerator() bool {
	n := r.Name
	return n == OPAdmin.String() || n == SuperAdmin.String() || n == Admin.String() || n == Moderator.String()
}
