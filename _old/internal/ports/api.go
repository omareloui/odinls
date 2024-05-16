package ports

import (
	"context"

	"github.com/omareloui/odinls/internal/application/core/domain"
)

type APIPort interface {
	Register(ctx context.Context, dto domain.Register) (*domain.User, error)
	Login(ctx context.Context, dto domain.Login) (*domain.User, error)
}
