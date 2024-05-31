package merchant

import (
	"errors"
	"time"

	"github.com/omareloui/odinls/internal/interfaces"
	"github.com/omareloui/odinls/internal/sanitizer"
)

var ErrMerchantNotFound = errors.New("merchant not found")

type merchantService struct {
	merchantRepository MerchantRepository
	validator          interfaces.Validator
}

func NewMerchantService(merchantRepository MerchantRepository, validator interfaces.Validator) MerchantService {
	return &merchantService{merchantRepository: merchantRepository, validator: validator}
}

func (s *merchantService) GetMerchants() ([]Merchant, error) {
	return s.merchantRepository.GetMerchants()
}

func (s *merchantService) GetMerchantByID(id string) (*Merchant, error) {
	return s.merchantRepository.FindMerchant(id)
}

func (s *merchantService) UpdateMerchantByID(id string, merchant *Merchant) error {
	sanitizeMerchant(merchant)

	if err := s.validator.Validate(merchant); err != nil {
		return s.validator.ParseError(err)
	}

	return s.merchantRepository.UpdateMerchantByID(id, merchant)
}

func (s *merchantService) CreateMerchant(merchant *Merchant) error {
	sanitizeMerchant(merchant)

	if err := s.validator.Validate(merchant); err != nil {
		return s.validator.ParseError(err)
	}

	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = time.Now()
	return s.merchantRepository.CreateMerchant(merchant)
}

func sanitizeMerchant(m *Merchant) {
	m.Name = sanitizer.TrimString(m.Name)
	m.Logo = sanitizer.TrimString(m.Logo)
}
