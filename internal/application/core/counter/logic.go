package counter

import (
	"errors"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/interfaces"
)

var (
	ErrCounterNotFound        = errors.New("counter not found")
	ErrAlreadyExistingCounter = errors.New("already existing counter for this merchant")
)

type counterService struct {
	repo      CounterRepository
	validator interfaces.Validator
}

func NewCounterService(repo CounterRepository, validator interfaces.Validator) *counterService {
	return &counterService{
		repo:      repo,
		validator: validator,
	}
}

func (s *counterService) AddOneToProduct(claims *jwtadapter.JwtAccessClaims, category string) (uint8, error) {
	if claims == nil || !claims.IsCraftsman() || !claims.Role.IsAdmin() {
		return 0, errs.ErrForbidden
	}

	merId := claims.CraftsmanInfo.MerchantID
	// TODO: if the cat key doesn't exist add it
	count, err := s.repo.AddOneToProduct(merId, category)
	if errors.Is(ErrCounterNotFound, err) {
		_, err := s.createCounter(claims, category)
		if err != nil {
			return 0, err
		}
		return s.AddOneToProduct(claims, category)
	}

	return count, nil
}

func (s *counterService) AddOneToOrder(claims *jwtadapter.JwtAccessClaims) (uint, error) {
	if claims == nil || !claims.IsCraftsman() || !claims.Role.IsAdmin() {
		return 0, errs.ErrForbidden
	}

	merId := claims.CraftsmanInfo.MerchantID
	count, err := s.repo.AddOneToOrder(merId)
	if errors.Is(ErrCounterNotFound, err) {
		_, err := s.createCounter(claims)
		if err != nil {
			return 0, err
		}
		return s.AddOneToOrder(claims)
	}

	return count, nil
}

func (s *counterService) createCounter(claims *jwtadapter.JwtAccessClaims, categories ...string) (*Counter, error) {
	if claims == nil || !claims.IsCraftsman() || !claims.Role.IsAdmin() {
		return nil, errs.ErrForbidden
	}

	merId := claims.CraftsmanInfo.MerchantID
	cntr := &Counter{
		MerchantID:    merId,
		OrdersNumber:  0,
		ProductsCodes: ProductsCodes{},
	}

	for _, cat := range categories {
		cntr.ProductsCodes[cat] = 0
	}

	err := s.repo.CreateCounter(cntr)
	if err != nil {
		return nil, err
	}

	return cntr, nil
}
