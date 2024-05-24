package application

import (
	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/application/core/role"
	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/interfaces"
	repository "github.com/omareloui/odinls/internal/repositories"
)

type Application struct {
	UserService     user.UserService
	MerchantService merchant.MerchantService
	RoleService     role.RoleService
	ClientService   client.ClientService
}

func NewApplication(repo repository.Repository, validator interfaces.Validator) *Application {
	roleService := role.NewRoleService(repo, validator)

	return &Application{
		UserService:     user.NewUserService(repo, roleService, validator),
		RoleService:     roleService,
		MerchantService: merchant.NewMerchantService(repo, validator),
		ClientService:   client.NewClientService(repo, validator),
	}
}
