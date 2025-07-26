package application

import (
	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/application/core/counter"
	"github.com/omareloui/odinls/internal/application/core/material"
	"github.com/omareloui/odinls/internal/application/core/order"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/application/core/supplier"
	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/interfaces"
	repository "github.com/omareloui/odinls/internal/repositories"
)

type Application struct {
	ClientService   client.ClientService
	MaterialService material.MaterialService
	OrderService    order.OrderService
	ProductService  product.ProductService
	SupplierService supplier.SupplierService
	UserService     user.UserService
}

func NewApplication(repo repository.Repository, validator interfaces.Validator, sanitizer interfaces.Sanitizer) *Application {
	counterService := counter.NewCounterService(repo)

	productService := product.NewProductService(repo, validator, sanitizer, counterService)

	return &Application{
		ClientService:   client.NewClientService(repo, validator, sanitizer),
		MaterialService: material.NewMaterialService(repo, validator, sanitizer),
		OrderService:    order.NewOrderService(repo, productService, counterService, validator, sanitizer),
		ProductService:  productService,
		SupplierService: supplier.NewSupplierService(repo, validator, sanitizer),
		UserService:     user.NewUserService(repo, validator, sanitizer),
	}
}
