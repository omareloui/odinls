package product

import (
	"errors"
	"slices"

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

func (s *productService) UpdateClientByID(claims *jwtadapter.JwtAccessClaims, id string, uprod *Product, options ...RetrieveOptsFunc) error {
	if claims == nil || !claims.Role.IsAdmin() || !claims.IsCraftsman() {
		return errs.ErrForbidden
	}

	sanitizeProduct(uprod)

	if err := s.validator.Validate(uprod); err != nil {
		return s.validator.ParseError(err)
	}

	prod, err := s.repo.GetProductByID(id)
	if err != nil {
		return err
	}

	if prod.Category != uprod.Category {
		newnum, err := s.counterService.AddOneToProduct(claims, uprod.Category)
		if err != nil {
			return err
		}
		uprod.Number = newnum
	} else {
		uprod.Number = prod.Number
	}

	uprod.ID = id
	uprod.MerchantID = prod.MerchantID
	uprod.CraftsmanID = prod.CraftsmanID
	uprod.CreatedAt = prod.CreatedAt

	for i := range uprod.Variants {
		uprod.Variants[i].ProductRef = uprod.Ref()
	}

	// This keeps the variant even if the new update data doesn't
	// include the same variant.
	for _, variant := range prod.Variants {
		idx := slices.IndexFunc(uprod.Variants, func(uvariant ProductVariant) bool {
			return variant.ID == uvariant.ID
		})
		if idx == -1 {
			uprod.Variants = append(uprod.Variants, variant)
		}
	}

	return s.repo.UpdateProductByID(id, uprod, options...)
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
