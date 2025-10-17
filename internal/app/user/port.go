package user

import (
	"context"

	"github.com/javito2003/shortener_url/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, bool, error)
	GetById(ctx context.Context, id string) (*domain.User, bool, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) bool
}
