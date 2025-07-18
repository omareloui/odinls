package counter

import (
	"errors"
	"time"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/errs"
)

var ErrCounterNotFound = errors.New("counter not found")

type counterService struct {
	repo CounterRepository
}

func NewCounterService(repo CounterRepository) *counterService {
	return &counterService{
		repo: repo,
	}
}

func (s *counterService) AddOneToProduct(claims *jwtadapter.JwtAccessClaims, category string) (uint8, error) {
	if claims == nil || !claims.IsCraftsman() || !claims.Role.IsAdmin() {
		return 0, errs.ErrForbidden
	}

	count, err := s.repo.AddOneToProduct(category)
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

	count, err := s.repo.AddOneToOrder()
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

	now := time.Now()

	cntr := &Counter{
		OrdersNumber:  0,
		ProductsCodes: ProductsCodes{},
		CreatedAt:     now,
		UpdatedAt:     now,
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
