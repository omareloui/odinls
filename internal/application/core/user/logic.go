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
	ErrUserNotFound          = errors.New("User Not Found")
	ErrEmailAlreadyExists    = errors.New("User Email Exists")
	ErrUsernameAlreadyExists = errors.New("User Username Exists")
)

type userService struct {
	userRepository UserRepository
	validator      interfaces.Validator
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

func NewUserService(userRepository UserRepository, validator interfaces.Validator) UserService {
	return &userService{userRepository: userRepository, validator: validator}
}

func sanitizeUser(usr *User) {
	usr.Name.First = sanitizer.TrimString(usr.Name.First)
	usr.Name.Last = sanitizer.TrimString(usr.Name.Last)
	usr.Email = sanitizer.TrimAndLowerCaseString(usr.Email)
	usr.Username = sanitizer.TrimAndLowerCaseString(usr.Username)
}
