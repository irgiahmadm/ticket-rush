package ports

import (
	"auth-service/internal/core/domain"
	"context"

	"github.com/google/uuid"
)

// Input Port (Service)
type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, email, password string) error
	GetUser(ctx context.Context, id uuid.UUID) (domain.User, error)
}

// Output Port (DB)
type UserRepository interface {
	Save(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (domain.User, error)
}

// Output Port (Token Generator)
type TokenGenerator interface {
	GenerateToken(user *domain.User) (string, error)
}
