package user

import (
	"errors"
	"time"

	"github.com/omareloui/odinls/internal/interfaces"
	"golang.org/x/crypto/bcrypt"
)

const passwordHashCost = 14

var ErrUserNotFound = errors.New("User Not Found")

type userService struct {
	userRepository UserRepository
	validator      interfaces.Validator
}

func (s *userService) FindUserByEmailOrUsername(emailOrPassword string) (*User, error) {
	return s.userRepository.FindUserByEmailOrUsername(emailOrPassword)
}

func (s *userService) FindUser(id string) (*User, error) {
	return s.userRepository.FindUser(id)
}

func (s *userService) CreateUser(usr *User) error {
	if err := s.validator.Validate(usr); err != nil {
		return s.validator.ParseError(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(usr.Password), passwordHashCost)
	if err != nil {
		return err
	}

	usr.Password = string(hash)

	usr.CreatedAt = time.Now()
	usr.UpdatedAt = time.Now()
	return s.userRepository.CreateUser(usr)
}

func NewUserService(userRepository UserRepository, validator interfaces.Validator) UserService {
	return &userService{userRepository: userRepository, validator: validator}
}
