package ports

import (
	"context"

	"github.com/omareloui/odinls/internal/application/core/domain"
)

type DBProt interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, dto domain.Register) (*domain.User, error)
}
