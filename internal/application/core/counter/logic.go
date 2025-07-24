package counter

import (
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/errs"
)

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

	return s.repo.AddOneToProduct(category)
}

func (s *counterService) AddOneToOrder(claims *jwtadapter.JwtAccessClaims) (uint, error) {
	if claims == nil || !claims.IsCraftsman() || !claims.Role.IsAdmin() {
		return 0, errs.ErrForbidden
	}

	return s.repo.AddOneToOrder()
}
