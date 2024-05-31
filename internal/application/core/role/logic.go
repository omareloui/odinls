package role

import (
	"errors"
	"time"

	"github.com/omareloui/odinls/internal/interfaces"
	"github.com/omareloui/odinls/internal/sanitizer"
)

var (
	ErrInvalidRole           = errors.New("invalid role")
	ErrRoleNotFound          = errors.New("role not found")
	ErrRoleNameAlreadyExists = errors.New("role name already exists")
)

type roleService struct {
	roleRepository RoleRepository
	validator      interfaces.Validator
}

func NewRoleService(roleRepository RoleRepository, validator interfaces.Validator) RoleService {
	return &roleService{
		roleRepository: roleRepository,
		validator:      validator,
	}
}

func (s *roleService) GetRoles() ([]Role, error) {
	return s.roleRepository.GetRoles()
}

func (s *roleService) GetRoleByID(id string) (*Role, error) {
	return s.roleRepository.FindRole(id)
}

func (s *roleService) CreateRole(role *Role) error {
	sanitizeRole(role)

	if err := s.validator.Validate(role); err != nil {
		return s.validator.ParseError(err)
	}

	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()
	return s.roleRepository.CreateRole(role)
}

func (s *roleService) SeedRoles() error {
	roles := []RoleEnum{
		OPAdmin, SuperAdmin,
		Admin, Moderator,
		NoAuthority,
	}

	rolestrs := []string{}

	for _, role := range roles {
		rolestrs = append(rolestrs, role.String())
	}

	return s.roleRepository.SeedRoles(rolestrs)
}

func (s *roleService) GetRoleByName(role string) (*RoleEnum, error) {
	var r RoleEnum
	var err error

	switch role {
	case OPAdmin.String():
		r = OPAdmin
	case SuperAdmin.String():
		r = SuperAdmin
	case Admin.String():
		r = Admin
	case Moderator.String():
		r = Moderator
	case NoAuthority.String():
		r = NoAuthority
	default:
		err = ErrInvalidRole
	}

	return &r, err
}

func sanitizeRole(r *Role) {
	r.Name = sanitizer.UpperCase(sanitizer.TrimString(r.Name))
}
