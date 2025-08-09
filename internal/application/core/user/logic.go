package user

import (
	"errors"

	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/interfaces"
	"golang.org/x/crypto/bcrypt"
)

const passwordHashCost = 14

type userService struct {
	repo      UserRepository
	validator interfaces.Validator
	sanitizer interfaces.Sanitizer
}

func NewUserService(userRepository UserRepository, validator interfaces.Validator, sanitizer interfaces.Sanitizer) UserService {
	return &userService{
		repo:      userRepository,
		validator: validator,
		sanitizer: sanitizer,
	}
}

func (s *userService) GetUsers() ([]User, error) {
	return s.repo.GetUsers()
}

func (s *userService) GetUserByEmailOrUsername(emailOrUsername string) (*User, error) {
	return s.repo.GetUserByEmailOrUsername(s.sanitizer.Lower(s.sanitizer.Trim(emailOrUsername)))
}

func (s *userService) GetUserByEmailOrUsernameFromUser(usr *User) (*User, error) {
	err := s.sanitizer.SanitizeStruct(usr)
	if err != nil {
		return nil, errs.ErrSanitizer
	}
	return s.repo.GetUserByEmailOrUsernameFromUser(usr)
}

func (s *userService) GetUserByID(id string) (*User, error) {
	return s.repo.GetUser(id)
}

func (s *userService) CreateUser(usr *User) (*User, error) {
	err := s.sanitizer.SanitizeStruct(usr)
	if err != nil {
		return nil, errs.ErrSanitizer
	}

	if err := s.validator.Validate(usr); err != nil {
		return nil, err
	}

	passStr := usr.Password
	hash, err := bcrypt.GenerateFromPassword([]byte(passStr), passwordHashCost)
	if err != nil {
		return nil, err
	}

	usr.Password = string(hash)

	return s.repo.CreateUser(usr)
}

func (s *userService) UpdateUserByID(id string, usr *User) (*User, error) {
	err := s.sanitizer.SanitizeStruct(usr)
	if err != nil {
		return nil, errs.ErrSanitizer
	}

	if err := s.validator.Validate(usr); err != nil {
		delete(err.Errors, "Password")
		if len(err.Errors) > 0 {
			return nil, err
		}
	}

	sameEmailUsr, err := s.repo.GetUserByEmailOrUsername(usr.Email)
	if err != nil && !errors.Is(err, errs.ErrDocumentNotFound) {
		return nil, err
	}
	if sameEmailUsr != nil && sameEmailUsr.ID != id {
		return nil, errs.ErrDocumentAlreadyExists
	}

	sameUsernameUsr, err := s.repo.GetUserByEmailOrUsername(usr.Username)
	if err != nil && !errors.Is(err, errs.ErrDocumentNotFound) {
		return nil, err
	}
	if sameUsernameUsr != nil && sameUsernameUsr.ID != id {
		return nil, errs.ErrDocumentAlreadyExists
	}

	if usr.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(usr.Password), passwordHashCost)
		if err != nil {
			return nil, err
		}

		usr.Password = string(hash)
	}

	return s.repo.UpdateUserByID(id, usr)
}

func (s *userService) UnsetCraftsmanByID(id string) (*User, error) {
	return s.repo.UnsetCraftsmanByID(id)
}
