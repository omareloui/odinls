package order

import (
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

func (s *orderService) GetOrders(claims *jwtadapter.AccessClaims, options ...RetrieveOptsFunc) ([]Order, error) {
	if claims == nil || !claims.Role.IsModerator() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetOrders(options...)
}

func (s *orderService) GetOrderByID(claims *jwtadapter.AccessClaims, id string, options ...RetrieveOptsFunc) (*Order, error) {
	if claims == nil || !claims.Role.IsModerator() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetOrderByID(id, options...)
}

func (s *orderService) CreateOrder(claims *jwtadapter.AccessClaims, ord *Order, options ...RetrieveOptsFunc) (*Order, error) {
	if claims == nil || !claims.Role.IsAdmin() || !claims.IsCraftsman() {
		return nil, errs.ErrForbidden
	}

	ord.Ref, _ = nanoid.Generate(refAlphabet, refSize)

	num, err := s.counterService.AddOneToOrder(claims)
	if err != nil {
		return nil, err
	}
	ord.Number = num

	if ord.Timeline.IssuanceDate.IsZero() {
		ord.Timeline.IssuanceDate = time.Now()
	}

	for i, item := range ord.Items {
		prod, err := s.productService.GetProductByVariantID(claims, item.Snapshot.VariantID)
		if err != nil {
			return nil, err
		}

		variantIdx := slices.IndexFunc(prod.Variants, func(v product.Variant) bool {
			return v.ID == ord.Items[i].Snapshot.VariantID
		})
		if variantIdx == -1 {
			log.Fatalln("invalid variant index: (searching a variant after getting it back by searching for a product with its id and its variant id)")
		}
		variant := prod.Variants[variantIdx]

		ord.Items[i].Snapshot.ProductID = prod.ID
		ord.Items[i].Snapshot.ProductName = prod.Name
		ord.Items[i].Snapshot.Category = prod.Category

		ord.Items[i].Snapshot.SKU = variant.SKU()
		ord.Items[i].Snapshot.VariantName = variant.Name
		ord.Items[i].Snapshot.Options = variant.Options

		ord.Items[i].Snapshot.Price = variant.Price

		ord.Items[i].Snapshot.TimeToCraft = variant.TimeToCraft

		ord.Items[i].Progress = ItemProgressNotStarted
	}

	err = s.sanitizer.SanitizeStruct(ord)
	if err != nil {
		return nil, errs.ErrSanitizer
	}

	if err := s.validator.Validate(ord); err != nil {
		return nil, err
	}

	return s.repo.CreateOrder(ord, options...)
}

func (s *orderService) UpdateOrderByID(claims *jwtadapter.AccessClaims, id string, uord *Order, options ...RetrieveOptsFunc) (*Order, error) {
	if claims == nil || !claims.Role.IsAdmin() || !claims.IsCraftsman() {
		return nil, errs.ErrForbidden
	}

	err := s.sanitizer.SanitizeStruct(uord)
	if err != nil {
		return nil, errs.ErrSanitizer
	}

	if err := s.validator.Validate(uord); err != nil {
		return nil, err
	}

	return s.repo.UpdateOrderByID(id, uord, options...)
}
