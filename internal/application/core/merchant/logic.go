package merchant

import (
	"errors"
	"time"

	"github.com/omareloui/odinls/internal/interfaces"
)

var ErrMerchantNotFound = errors.New("Merchant Not Found")

type merchantService struct {
	merchantRepository MerchantRepository
	validator          interfaces.Validator
}

func (s *merchantService) GetMerchants() ([]Merchant, error) {
	return s.merchantRepository.GetMerchants()
}

func (s *merchantService) FindMerchant(id string) (*Merchant, error) {
	return s.merchantRepository.FindMerchant(id)
}

func (s *merchantService) CreateMerchant(merchant *Merchant) error {
	if err := s.validator.Validate(merchant); err != nil {
		return s.validator.ParseError(err)
	}

	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = time.Now()
	return s.merchantRepository.CreateMerchant(merchant)
}

func NewMerchantService(merchantRepository MerchantRepository, validator interfaces.Validator) MerchantService {
	return &merchantService{merchantRepository: merchantRepository, validator: validator}
}
