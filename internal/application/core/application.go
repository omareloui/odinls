package application

import (
	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/interfaces"
	repository "github.com/omareloui/odinls/internal/repositories"
)

type Application struct {
	UserService     user.UserService
	MerchantService merchant.MerchantService
}

func NewApplication(repo repository.Repository, validator interfaces.Validator) *Application {
	return &Application{
		UserService:     user.NewUserService(repo, validator),
		MerchantService: merchant.NewMerchantService(repo, validator),
	}
}
