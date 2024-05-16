package application

import (
	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/interfaces"
	repository "github.com/omareloui/odinls/internal/repositories"
)

type Application struct {
	userService     user.UserService
	merchantService merchant.MerchantService
}

func NewApplication(repo repository.Repository, validator interfaces.Validator) *Application {
	return &Application{
		userService:     user.NewUserService(repo, validator),
		merchantService: merchant.NewMerchantService(repo, validator),
	}
}
