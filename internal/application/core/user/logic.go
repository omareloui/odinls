package user

import (
	"errors"
	"time"

	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/application/core/role"
	"github.com/omareloui/odinls/internal/interfaces"
	"github.com/omareloui/odinls/internal/sanitizer"
	"golang.org/x/crypto/bcrypt"
)

const passwordHashCost = 14

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyExists    = errors.New("user email exists")
	ErrUsernameAlreadyExists = errors.New("user username Exists")
)

type userService struct {
	userRepository  UserRepository
	roleService     role.RoleService
	merchantService merchant.MerchantService
	validator       interfaces.Validator
}

func NewUserService(userRepository UserRepository, merchantService merchant.MerchantService, roleService role.RoleService, validator interfaces.Validator) UserService {
	return &userService{
		userRepository:  userRepository,
		roleService:     roleService,
		merchantService: merchantService,
		validator:       validator,
	}
}

func (s *userService) GetUsers(opts ...RetrieveOptsFunc) ([]User, error) {
	return s.userRepository.GetUsers(opts...)
}

func (s *userService) FindUserByEmailOrUsername(emailOrPassword string, opts ...RetrieveOptsFunc) (*User, error) {
	return s.userRepository.FindUserByEmailOrUsername(sanitizer.TrimAndLowerCaseString(emailOrPassword), opts...)
}

func (s *userService) FindUserByEmailOrUsernameFromUser(usr *User, opts ...RetrieveOptsFunc) (*User, error) {
	return s.userRepository.FindUserByEmailOrUsernameFromUser(&User{
		Username: sanitizer.TrimAndLowerCaseString(usr.Username),
		Email:    sanitizer.TrimAndLowerCaseString(usr.Email),
	}, opts...)
}

func (s *userService) FindUser(id string, opts ...RetrieveOptsFunc) (*User, error) {
	return s.userRepository.FindUser(id, opts...)
}

func (s *userService) CreateUser(usr *User, opts ...RetrieveOptsFunc) error {
	sanitizeUser(usr)

	if err := s.validator.Validate(usr); err != nil {
		return s.validator.ParseError(err)
	}

	passStr := usr.Password
	hash, err := bcrypt.GenerateFromPassword([]byte(passStr), passwordHashCost)
	if err != nil {
		return err
	}

	usr.Password = string(hash)

	usr.CreatedAt = time.Now()
	usr.UpdatedAt = time.Now()

	err = s.userRepository.CreateUser(usr, opts...)
	if err != nil {
		usr.Password = passStr
		usr.ConfirmPassword = passStr
	}

	return err
}

func (s *userService) UpdateUserByID(id string, usr *User, opts ...RetrieveOptsFunc) error {
	sanitizeUser(usr)

	if err := s.validator.Validate(usr); err != nil {
		valerr := s.validator.ParseError(err)
		delete(valerr.Errors, "Password")
		if len(valerr.Errors) > 0 {
			return valerr
		}
	}

	sameEmailUsr, err := s.userRepository.FindUserByEmailOrUsername(usr.Email)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return err
	}
	if sameEmailUsr != nil && sameEmailUsr.ID != id {
		return ErrEmailAlreadyExists
	}

	sameUsernameUsr, err := s.userRepository.FindUserByEmailOrUsername(usr.Username)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return err
	}
	if sameUsernameUsr != nil && sameUsernameUsr.ID != id {
		return ErrUsernameAlreadyExists
	}

	if _, err = s.roleService.FindRole(usr.RoleID); err != nil {
		return err
	}

	if usr.Craftsman != nil {
		if _, err = s.merchantService.FindMerchant(usr.Craftsman.MerchantID); err != nil {
			return err
		}
	}

	return s.userRepository.UpdateUserByID(id, usr, opts...)
}

func (s *userService) UnsetCraftsmanByID(id string) error {
	return s.userRepository.UnsetCraftsmanByID(id)
}

func sanitizeUser(usr *User) {
	usr.Name.First = sanitizer.TrimString(usr.Name.First)
	usr.Name.Last = sanitizer.TrimString(usr.Name.Last)
	usr.Email = sanitizer.TrimAndLowerCaseString(usr.Email)
	usr.Username = sanitizer.TrimAndLowerCaseString(usr.Username)
}
