package supplier

import (
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/interfaces"
)

type supplierService struct {
	repo      SupplierRepository
	validator interfaces.Validator
	sanitizer interfaces.Sanitizer
}

func NewSupplierService(supplierRepository SupplierRepository, validator interfaces.Validator, sanitizer interfaces.Sanitizer) *supplierService {
	return &supplierService{
		repo:      supplierRepository,
		validator: validator,
		sanitizer: sanitizer,
	}
}

func (s *supplierService) GetSuppliers(claims *jwtadapter.AccessClaims) ([]Supplier, error) {
	if claims == nil || !claims.Role.IsModerator() || !claims.IsCraftsman() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetSuppliers()
}

func (s *supplierService) GetSupplierByID(claims *jwtadapter.AccessClaims, id string) (*Supplier, error) {
	if claims == nil || !claims.Role.IsModerator() || !claims.IsCraftsman() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetSupplierByID(id)
}

func (s *supplierService) CreateSupplier(claims *jwtadapter.AccessClaims, sup *Supplier) (*Supplier, error) {
	if claims == nil || !claims.Role.IsAdmin() {
		return nil, errs.ErrForbidden
	}

	if err := s.sanitizer.SanitizeStruct(sup); err != nil {
		return nil, err
	}

	if err := s.validator.Validate(sup); err != nil {
		return nil, s.validator.ParseError(err)
	}

	return s.repo.CreateSupplier(sup)
}

func (s *supplierService) UpdateSupplierByID(claims *jwtadapter.AccessClaims, id string, sup *Supplier) (*Supplier, error) {
	if claims == nil || !claims.Role.IsAdmin() {
		return nil, errs.ErrForbidden
	}

	if err := s.sanitizer.SanitizeStruct(sup); err != nil {
		return nil, err
	}

	if err := s.validator.Validate(sup); err != nil {
		return nil, s.validator.ParseError(err)
	}

	return s.repo.UpdateSupplierByID(id, sup)
}
