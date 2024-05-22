package role

type RoleService interface {
	SeedRoles() error
	GetRoles() ([]Role, error)
	FindRole(id string) (*Role, error)
	FindRoleByName(role string) (*RoleEnum, error)
	CreateRole(*Role) error
}
