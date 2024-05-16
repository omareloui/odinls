package repository

import (
	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/application/core/user"
)

type Repository interface {
	merchant.MerchantRepository
	user.UserRepository
}
