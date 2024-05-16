package merchant

import (
	"errors"
	"time"

	"github.com/omareloui/odinls/internal/interfaces"
)

var (
	ErrMerchantNotFound = errors.New("Merchant Not Found")
	ErrMerchantInvalid  = errors.New("Invalid Merchant")
)

type merchantService struct {
	merchantRepository MerchantRepository
	validator          interfaces.Validator
}

func (s *merchantService) FindMerchant(id string) (*Merchant, error) {
	return s.merchantRepository.FindMerchant(id)
}

func (s *merchantService) CreateMerchant(merchant *Merchant) error {
	if err := s.validator.Validate(merchant); err != nil {
		// TODO: find a way to send the error details "err" with the error
		return ErrMerchantInvalid
	}

	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = time.Now()
	return s.merchantRepository.CreateMerchant(merchant)
}

func NewMerchantService(merchantRepository MerchantRepository, validator interfaces.Validator) MerchantService {
	return &merchantService{merchantRepository: merchantRepository, validator: validator}
}
