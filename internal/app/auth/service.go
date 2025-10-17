package auth

import (
	"context"

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
	token  TokenGenerator
}

func NewService(repo user.UserRepository, hasher user.PasswordHasher, token TokenGenerator) *Service {
	return &Service{repo: repo, hasher: hasher, token: token}
}

func (s *Service) CreateUser(ctx context.Context, firstName, lastName, email, password string) (*domain.User, error) {
	if _, found, err := s.repo.FindByEmail(ctx, email); err != nil {
		return nil, err
	} else if found {
		return nil, ErrAlreadyLoggedIn
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
		return "", ErrInvalidCredentials
	}

	same := s.hasher.Compare(user.Password, password)

	if !found || !same {
		return "", ErrInvalidCredentials
	}

	token, err := s.token.Generate(ctx, user.ID, nil)
	if err != nil {
		return "", err
	}

	return token, nil
}
