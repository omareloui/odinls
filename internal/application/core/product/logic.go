package product

import (
	"errors"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/application/core/counter"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/interfaces"
	"github.com/omareloui/odinls/internal/sanitizer"
)

var ErrProductNotFound = errors.New("product not found")

type productService struct {
	repo           ProductRepository
	validator      interfaces.Validator
	counterService counter.CounterService
}

func NewProductService(repo ProductRepository, validator interfaces.Validator, counterService counter.CounterService) *productService {
	return &productService{
		repo:           repo,
		validator:      validator,
		counterService: counterService,
	}
}

func (s *productService) GetProducts(claims *jwtadapter.JwtAccessClaims, options ...RetrieveOptsFunc) ([]Product, error) {
	if claims == nil || !claims.Role.IsOPAdmin() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetProducts(options...)
}

func (s *productService) GetProductByID(claims *jwtadapter.JwtAccessClaims, id string, options ...RetrieveOptsFunc) (*Product, error) {
	if claims == nil || !claims.Role.IsOPAdmin() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetProductByID(id, options...)
}

func (s *productService) GetCurrentMerchantProducts(claims *jwtadapter.JwtAccessClaims, options ...RetrieveOptsFunc) ([]Product, error) {
	if claims == nil && !claims.IsCraftsman() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetProductsByMerchantID(claims.CraftsmanInfo.MerchantID, options...)
}

func (s *productService) CreateProduct(claims *jwtadapter.JwtAccessClaims, prod *Product, options ...RetrieveOptsFunc) error {
	if claims == nil || !claims.Role.IsAdmin() || !claims.IsCraftsman() {
		return errs.ErrForbidden
	}

	sanitizeProduct(prod)

	if err := s.validator.Validate(prod); err != nil {
		return s.validator.ParseError(err)
	}

	num, err := s.counterService.AddOneToProduct(claims, prod.Category)
	if err != nil {
		return err
	}

	prod.Number = num

	prod.CraftsmanID = claims.ID
	prod.MerchantID = claims.CraftsmanInfo.MerchantID

	for i := range prod.Variants {
		prod.Variants[i].ProductRef = prod.Ref()
	}

	return s.repo.CreateProduct(prod, options...)
}

func (s *productService) UpdateClientByID(claims *jwtadapter.JwtAccessClaims, id string, prod *Product, options ...RetrieveOptsFunc) error {
	if claims == nil || !claims.Role.IsAdmin() || !claims.IsCraftsman() {
		return errs.ErrForbidden
	}

	sanitizeProduct(prod)

	if err := s.validator.Validate(prod); err != nil {
		return s.validator.ParseError(err)
	}

	for i := range prod.Variants {
		prod.Variants[i].ProductRef = prod.Ref()
	}

	return s.repo.UpdateProductByID(id, prod, options...)
}

func sanitizeProduct(prod *Product) {
	prod.Name = sanitizer.TrimString(prod.Name)
	prod.Description = sanitizer.TrimString(prod.Description)
	for i := range prod.Variants {
		prod.Variants[i].Suffix = sanitizer.LowerCase(sanitizer.TrimString(prod.Variants[i].Suffix))
		prod.Variants[i].Name = sanitizer.TrimString(prod.Variants[i].Name)
		prod.Variants[i].Description = sanitizer.TrimString(prod.Variants[i].Description)
		prod.Variants[i].ProductRef = sanitizer.TrimString(prod.Variants[i].ProductRef)
	}
}
