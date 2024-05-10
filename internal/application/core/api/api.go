package api

import (
	"context"
	"errors"

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
	// TODO: validation
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
