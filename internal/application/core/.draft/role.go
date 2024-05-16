package domain

type RoleEnum int

const (
	OPAdmin RoleEnum = iota
	SuperAdmin
	Admin
	Moderator
	NoAuthority
)

func (r *RoleEnum) String() string {
	return [...]string{
		"op_admin", "super_admin",
		"admin", "moderator",
		"no_authority",
	}[*r]
}

type Role struct {
	ID   ID
	Name RoleEnum
	// Permissions []Permission
}

// type Permission struct {
// 	Subject       string // enum
// 	Action        string // enum
// 	IsCraftsman   bool
// 	IsOwn         bool
// 	IsOwnMerchant bool
// }
