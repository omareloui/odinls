package api

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/omareloui/odinls/internal/application/core/domain"
	"github.com/omareloui/odinls/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type Application struct {
	db ports.DBProt
}

func NewApplication(db ports.DBProt) *Application {
	return &Application{db: db}
}

func (a *Application) Register(ctx context.Context, dto domain.Register) (*domain.User, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("unique_email", func(fl validator.FieldLevel) bool {
		matchingEmail, err := a.db.GetUserByEmail(ctx, dto.Email)
		if err != nil {
			return false
		}
		if matchingEmail != nil {
			return false
		}
		return true
	})

	err := validate.Struct(dto)
	if err != nil {
		return nil, err
	}

	return a.db.CreateUser(ctx, dto)
}

func (a *Application) Login(ctx context.Context, dto domain.Login) (*domain.User, error) {
	// TODO: validation

	usr, err := a.db.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(dto.Password))
	isValidPass := err == nil

	if !isValidPass {
		return nil, errors.New("invalid password")
	}

	return usr, nil
}
