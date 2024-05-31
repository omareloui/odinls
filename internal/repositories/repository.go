package repository

import (
	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/application/core/counter"
	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/application/core/order"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/application/core/role"
	"github.com/omareloui/odinls/internal/application/core/user"
)

type Repository interface {
	merchant.MerchantRepository
	user.UserRepository
	role.RoleRepository
	client.ClientRepository
	counter.CounterRepository
	product.ProductRepository
	order.OrderRepository
}
