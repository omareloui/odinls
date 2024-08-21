package role

type RoleService interface {
	SeedRoles() error
	GetRoles() ([]Role, error)
	GetRoleByID(id string) (*Role, error)
	GetRoleByName(name string) (*Role, error)
	MapRoleNameToRoleEnum(role string) (*RoleEnum, error)
	CreateRole(*Role) error
}
