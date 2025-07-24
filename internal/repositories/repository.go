package repository

import (
	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/application/core/counter"
	"github.com/omareloui/odinls/internal/application/core/material"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/application/core/supplier"
	"github.com/omareloui/odinls/internal/application/core/user"
	// "github.com/omareloui/odinls/internal/application/core/order"
)

type Repository interface {
	client.ClientRepository
	counter.CounterRepository
	material.MaterialRepository
	product.ProductRepository
	supplier.SupplierRepository
	// order.OrderRepository
	user.UserRepository
}
