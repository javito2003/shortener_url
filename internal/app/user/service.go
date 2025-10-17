package user

import (
	"context"

	"github.com/javito2003/shortener_url/internal/domain"
)

type UserService interface {
	GetMe(ctx context.Context, userID string) (*domain.User, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetMe(ctx context.Context, userID string) (*domain.User, error) {
	user, found, err := s.userRepo.GetById(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrUserNotFound
	}

	return user, nil
}
