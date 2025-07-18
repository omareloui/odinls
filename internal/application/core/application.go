package application

import (
	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/application/core/counter"
	"github.com/omareloui/odinls/internal/application/core/order"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/interfaces"
	repository "github.com/omareloui/odinls/internal/repositories"
)

type Application struct {
	UserService    user.UserService
	ClientService  client.ClientService
	ProductService product.ProductService
	OrderService   order.OrderService
}

func NewApplication(repo repository.Repository, validator interfaces.Validator, sanitizer interfaces.Sanitizer) *Application {
	counterService := counter.NewCounterService(repo)

	productService := product.NewProductService(repo, validator, sanitizer, counterService)

	return &Application{
		UserService:    user.NewUserService(repo, validator, sanitizer),
		ClientService:  client.NewClientService(repo, validator, sanitizer),
		ProductService: productService,
		OrderService:   order.NewOrderService(repo, productService, counterService, validator, sanitizer),
	}
}
