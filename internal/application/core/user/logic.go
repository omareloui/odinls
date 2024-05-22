package user

import (
	"errors"
	"time"

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
	userRepository UserRepository
	validator      interfaces.Validator
}

func NewUserService(userRepository UserRepository, validator interfaces.Validator) UserService {
	return &userService{userRepository: userRepository, validator: validator}
}

func (s *userService) GetUsers() ([]User, error) {
	return s.userRepository.GetUsers()
}

func (s *userService) FindUserByEmailOrUsername(emailOrPassword string) (*User, error) {
	return s.userRepository.FindUserByEmailOrUsername(sanitizer.TrimAndLowerCaseString(emailOrPassword))
}

func (s *userService) FindUserByEmailOrUsernameFromUser(usr *User) (*User, error) {
	return s.userRepository.FindUserByEmailOrUsernameFromUser(&User{
		Username: sanitizer.TrimAndLowerCaseString(usr.Username),
		Email:    sanitizer.TrimAndLowerCaseString(usr.Email),
	})
}

func (s *userService) FindUser(id string) (*User, error) {
	return s.userRepository.FindUser(id)
}

func (s *userService) CreateUser(usr *User) error {
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

	err = s.userRepository.CreateUser(usr)
	if err != nil {
		usr.Password = passStr
		usr.ConfirmPassword = passStr
	}

	return err
}

func (s *userService) UpdateUserByID(id string, usr *User) error {
	type updateUser struct {
		Name     Name   `validate:"required"`
		Username string `validate:"required,min=3,max=64,alphanum_with_underscore,not_blank"`
		Email    string `validate:"required,email,not_blank"`
		Phone    string
		Role     string
	}

	sanitizeUser(usr)

	u := &updateUser{
		Name:     usr.Name,
		Username: usr.Username,
		Email:    usr.Email,
		Phone:    usr.Phone,
		Role:     usr.Role,
	}

	if err := s.validator.Validate(u); err != nil {
		return s.validator.ParseError(err)
	}

	sameEmailUsr, err := s.userRepository.FindUserByEmailOrUsername(u.Email)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return err
	}
	if sameEmailUsr != nil && sameEmailUsr.ID != id {
		return ErrEmailAlreadyExists
	}

	sameUsernameUsr, err := s.userRepository.FindUserByEmailOrUsername(u.Username)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return err
	}
	if sameUsernameUsr != nil && sameUsernameUsr.ID != id {
		return ErrUsernameAlreadyExists
	}

	return s.userRepository.UpdateUserByID(id, usr)
}

func sanitizeUser(usr *User) {
	usr.Name.First = sanitizer.TrimString(usr.Name.First)
	usr.Name.Last = sanitizer.TrimString(usr.Name.Last)
	usr.Email = sanitizer.TrimAndLowerCaseString(usr.Email)
	usr.Username = sanitizer.TrimAndLowerCaseString(usr.Username)
}
