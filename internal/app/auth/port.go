package auth

import (
	"context"
	"time"
)

type Claims struct {
	UserID string
}

type TokenGenerator interface {
	Generate(ctx context.Context, userID string, ttl *time.Duration) (string, error)
	Validate(ctx context.Context, tokenString string) (*Claims, error)
}
