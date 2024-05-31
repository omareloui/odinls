package application

import (
	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/application/core/counter"
	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/application/core/order"
	"github.com/omareloui/odinls/internal/application/core/product"
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
	ProductService  product.ProductService
	OrderService    order.OrderService
}

func NewApplication(repo repository.Repository, validator interfaces.Validator) *Application {
	role := role.NewRoleService(repo, validator)
	merchantService := merchant.NewMerchantService(repo, validator)
	counterService := counter.NewCounterService(repo)

	return &Application{
		UserService:     user.NewUserService(repo, merchantService, role, validator),
		RoleService:     role,
		MerchantService: merchantService,
		ClientService:   client.NewClientService(repo, validator),
		ProductService:  product.NewProductService(repo, validator, counterService),
		OrderService:    order.NewOrderService(repo, validator, counterService),
	}
}
