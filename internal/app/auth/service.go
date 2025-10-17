package auth

import (
	"context"
	"errors"

	"github.com/javito2003/shortener_url/internal/app/user"
	"github.com/javito2003/shortener_url/internal/domain"
)

type AuthService interface {
	Authenticate(ctx context.Context, email, password string) (string, error)
	CreateUser(ctx context.Context, firstName, lastName, email, password string) (*domain.User, error)
}

type Service struct {
	repo   user.UserRepository
	hasher user.PasswordHasher
}

func NewService(repo user.UserRepository, hasher user.PasswordHasher) *Service {
	return &Service{repo: repo, hasher: hasher}
}

func (s *Service) CreateUser(ctx context.Context, firstName, lastName, email, password string) (*domain.User, error) {
	if _, found, err := s.repo.FindByEmail(ctx, email); err != nil {
		return nil, err
	} else if found {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  hashedPassword,
	}

	return s.repo.Save(ctx, user)
}

func (s *Service) Authenticate(ctx context.Context, email, password string) (string, error) {
	user, found, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	same := s.hasher.Compare(user.Password, password)

	if !found || !same {
		return "", errors.New("invalid credentials")
	}

	return "authenticated-token", nil
}
