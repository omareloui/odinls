package order

import (
	"errors"

	"github.com/aidarkhanov/nanoid"
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/application/core/counter"
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
	counterService counter.CounterService
}

func NewOrderService(repo OrderRepository, counterService counter.CounterService, validator interfaces.Validator, sanitizer interfaces.Sanitizer) *orderService {
	return &orderService{
		repo:           repo,
		validator:      validator,
		sanitizer:      sanitizer,
		counterService: counterService,
	}
}

func (s *orderService) GetOrders(claims *jwtadapter.JwtAccessClaims, options ...RetrieveOptsFunc) ([]Order, error) {
	if claims == nil || !claims.Role.IsOPAdmin() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetOrders(options...)
}

func (s *orderService) GetOrderByID(claims *jwtadapter.JwtAccessClaims, id string, options ...RetrieveOptsFunc) (*Order, error) {
	if claims == nil || !claims.Role.IsOPAdmin() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetOrderByID(id, options...)
}

func (s *orderService) GetCurrentMerchantOrders(claims *jwtadapter.JwtAccessClaims, options ...RetrieveOptsFunc) ([]Order, error) {
	if claims == nil && !claims.IsCraftsman() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetOrdersByMerchantID(claims.CraftsmanInfo.MerchantID, options...)
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

	num, err := s.counterService.AddOneToOrder(claims)
	if err != nil {
		return err
	}

	ord.Number = num
	ord.MerchantID = claims.CraftsmanInfo.MerchantID
	ord.Ref, _ = nanoid.Generate(refAlphabet, refSize)

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

	ord, err := s.repo.GetOrderByID(id)
	if err != nil {
		return err
	}

	uord.ID = id
	uord.MerchantID = ord.MerchantID
	uord.CreatedAt = ord.CreatedAt

	return s.repo.UpdateOrderByID(id, uord, options...)
}
