package role

type RoleRepository interface {
	SeedRoles(roles []string) error
	GetRoles() ([]Role, error)
	FindRole(id string) (*Role, error)
	FindRoleByName(name string) (*Role, error)
	CreateRole(roles *Role) error
}
