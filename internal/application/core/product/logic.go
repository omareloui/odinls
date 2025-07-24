package product

import (
	"slices"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/application/core/counter"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/interfaces"
)

type productService struct {
	repo           ProductRepository
	validator      interfaces.Validator
	sanitizer      interfaces.Sanitizer
	counterService counter.CounterService
}

func NewProductService(repo ProductRepository, validator interfaces.Validator, sanitizer interfaces.Sanitizer, counterService counter.CounterService) *productService {
	return &productService{
		repo:           repo,
		validator:      validator,
		sanitizer:      sanitizer,
		counterService: counterService,
	}
}

func (s *productService) GetProducts(claims *jwtadapter.JwtAccessClaims, options ...RetrieveOptsFunc) ([]Product, error) {
	return s.repo.GetProducts(options...)
}

func (s *productService) GetProductByID(claims *jwtadapter.JwtAccessClaims, id string, options ...RetrieveOptsFunc) (*Product, error) {
	return s.repo.GetProductByID(id, options...)
}

func (s *productService) GetProductByVariantID(claims *jwtadapter.JwtAccessClaims, id string, options ...RetrieveOptsFunc) (*Product, error) {
	return s.repo.GetProductByVariantID(id, options...)
}

func (s *productService) GetProductByIDAndVariantID(claims *jwtadapter.JwtAccessClaims, id string, variantId string, options ...RetrieveOptsFunc) (*Product, error) {
	return s.repo.GetProductByIDAndVariantID(id, variantId, options...)
}

func (s *productService) CreateProduct(claims *jwtadapter.JwtAccessClaims, prod *Product, options ...RetrieveOptsFunc) (*Product, error) {
	if claims == nil || !claims.Role.IsAdmin() || !claims.IsCraftsman() {
		return nil, errs.ErrForbidden
	}

	err := s.sanitizer.SanitizeStruct(prod)
	if err != nil {
		return nil, errs.ErrSanitizer
	}

	if err := s.validator.Validate(prod); err != nil {
		return nil, s.validator.ParseError(err)
	}

	num, err := s.counterService.AddOneToProduct(claims, prod.Category.Code())
	if err != nil {
		return nil, err
	}

	prod.Number = num

	for i := range prod.Variants {
		prod.Variants[i].ProductSKU = prod.SKU()
		if prod.Variants[i].Price == 0 {
			// TODO: populate the used materials to set the price
			prod.Variants[i].Price = prod.Variants[i].EstPrice()
		}
		if prod.Variants[i].WholesalePrice == 0 {
			// TODO: populate the used materials to set the price
			prod.Variants[i].Price = prod.Variants[i].EstWholesalePrice()
		}
	}

	return s.repo.CreateProduct(prod, options...)
}

func (s *productService) UpdateProductByID(claims *jwtadapter.JwtAccessClaims, id string, uprod *Product, options ...RetrieveOptsFunc) (*Product, error) {
	if claims == nil || !claims.Role.IsAdmin() || !claims.IsCraftsman() {
		return nil, errs.ErrForbidden
	}

	err := s.sanitizer.SanitizeStruct(uprod)
	if err != nil {
		return nil, errs.ErrSanitizer
	}

	if err := s.validator.Validate(uprod); err != nil {
		return nil, s.validator.ParseError(err)
	}

	prod, err := s.repo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	if prod.Category != uprod.Category {
		newnum, err := s.counterService.AddOneToProduct(claims, uprod.Category.Code())
		if err != nil {
			return nil, err
		}
		uprod.Number = newnum
	} else {
		uprod.Number = prod.Number
	}

	uprod.ID = id
	uprod.CreatedAt = prod.CreatedAt

	for i := range uprod.Variants {
		prod.Variants[i].ProductSKU = prod.SKU()
		if prod.Variants[i].Price == 0 {
			// TODO: populate the used materials to set the price
			prod.Variants[i].Price = prod.Variants[i].EstPrice()
		}
		if prod.Variants[i].WholesalePrice == 0 {
			// TODO: populate the used materials to set the price
			prod.Variants[i].Price = prod.Variants[i].EstWholesalePrice()
		}
	}

	// This keeps the variant even if the new update data doesn't
	// include the same variant.
	for _, variant := range prod.Variants {
		idx := slices.IndexFunc(uprod.Variants, func(uvariant Variant) bool {
			return variant.ID == uvariant.ID
		})
		if idx == -1 {
			uprod.Variants = append(uprod.Variants, variant)
		}
	}

	return s.repo.UpdateProductByID(id, uprod, options...)
}
