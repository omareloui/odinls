package order

import (
	"errors"
	"log"
	"slices"
	"time"

	"github.com/aidarkhanov/nanoid"
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/application/core/counter"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/interfaces"
)

const (
	refAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	refSize     = 8
)

var ErrOrderNotFound = errors.New("order not found")

type orderService struct {
	repo           OrderRepository
	validator      interfaces.Validator
	sanitizer      interfaces.Sanitizer
	productService product.ProductService
	counterService counter.CounterService
}

func NewOrderService(repo OrderRepository, productService product.ProductService, counterService counter.CounterService, validator interfaces.Validator, sanitizer interfaces.Sanitizer) *orderService {
	return &orderService{
		repo:           repo,
		validator:      validator,
		sanitizer:      sanitizer,
		productService: productService,
		counterService: counterService,
	}
}

func (s *orderService) GetOrders(claims *jwtadapter.JwtAccessClaims, options ...RetrieveOptsFunc) ([]Order, error) {
	if claims == nil || !claims.Role.IsModerator() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetOrders(options...)
}

func (s *orderService) GetOrderByID(claims *jwtadapter.JwtAccessClaims, id string, options ...RetrieveOptsFunc) (*Order, error) {
	if claims == nil || !claims.Role.IsModerator() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetOrderByID(id, options...)
}

func (s *orderService) CreateOrder(claims *jwtadapter.JwtAccessClaims, ord *Order, options ...RetrieveOptsFunc) error {
	if claims == nil || !claims.Role.IsAdmin() || !claims.IsCraftsman() {
		return errs.ErrForbidden
	}

	err := s.sanitizer.SanitizeStruct(ord)
	if err != nil {
		return errs.ErrSanitizer
	}

	if err := s.validator.Validate(ord); err != nil {
		return s.validator.ParseError(err)
	}

	for i, item := range ord.Items {
		// Set the price
		prod, err := s.productService.GetProductByIDAndVariantID(claims, item.ProductID, item.VariantID)
		if err != nil {
			return err
		}
		variantIdx := slices.IndexFunc(prod.Variants, func(v product.Variant) bool {
			return v.ID == ord.Items[i].VariantID
		})
		if variantIdx == -1 {
			log.Fatalln("invalid variant index: (searching a variant after getting it back by searching for a product with its id and its variant id)")
		}
		ord.Items[i].Price = prod.Variants[variantIdx].Price

		// Set the default progress status
		ord.Items[i].Progress = ItemProgressNotStarted.String()
	}

	ord.Ref, _ = nanoid.Generate(refAlphabet, refSize)

	ord.Subtotal = ord.calcSubtotal()

	num, err := s.counterService.AddOneToOrder(claims)
	if err != nil {
		return err
	}
	ord.Number = num

	now := time.Now()
	ord.CreatedAt = now
	ord.UpdatedAt = now

	return s.repo.CreateOrder(ord, options...)
}

func (s *orderService) UpdateOrderByID(claims *jwtadapter.JwtAccessClaims, id string, uord *Order, options ...RetrieveOptsFunc) error {
	if claims == nil || !claims.Role.IsAdmin() || !claims.IsCraftsman() {
		return errs.ErrForbidden
	}

	err := s.sanitizer.SanitizeStruct(uord)
	if err != nil {
		return errs.ErrSanitizer
	}

	if err := s.validator.Validate(uord); err != nil {
		return s.validator.ParseError(err)
	}

	uord.UpdatedAt = time.Now()

	return s.repo.UpdateOrderByID(id, uord, options...)
}
