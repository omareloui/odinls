package merchant

import (
	"errors"
	"time"

	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/interfaces"
)

var ErrMerchantNotFound = errors.New("merchant not found")

type merchantService struct {
	merchantRepository MerchantRepository
	validator          interfaces.Validator
	sanitizer          interfaces.Sanitizer
}

func NewMerchantService(merchantRepository MerchantRepository, validator interfaces.Validator, sanitizer interfaces.Sanitizer) MerchantService {
	return &merchantService{merchantRepository: merchantRepository, validator: validator, sanitizer: sanitizer}
}

func (s *merchantService) GetMerchants() ([]Merchant, error) {
	return s.merchantRepository.GetMerchants()
}

func (s *merchantService) GetMerchantByID(id string) (*Merchant, error) {
	return s.merchantRepository.FindMerchant(id)
}

func (s *merchantService) UpdateMerchantByID(id string, merchant *Merchant) error {
	err := s.sanitizer.SanitizeStruct(merchant)
	if err != nil {
		return errs.ErrSanitizer
	}

	if err := s.validator.Validate(merchant); err != nil {
		return s.validator.ParseError(err)
	}

	merchant.UpdatedAt = time.Now()

	return s.merchantRepository.UpdateMerchantByID(id, merchant)
}

func (s *merchantService) CreateMerchant(merchant *Merchant) error {
	err := s.sanitizer.SanitizeStruct(merchant)
	if err != nil {
		return errs.ErrSanitizer
	}

	if err := s.validator.Validate(merchant); err != nil {
		return s.validator.ParseError(err)
	}

	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = time.Now()
	return s.merchantRepository.CreateMerchant(merchant)
}
