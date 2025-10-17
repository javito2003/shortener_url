package shortener

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"

	link "github.com/javito2003/shortener_url/internal/domain"
)

type Shortener interface {
	Shorten(ctx context.Context, url, user string) (string, error)
	Resolve(ctx context.Context, shortCode string) (string, error)
	GetByUser(ctx context.Context, userID string, limit, skip int32) ([]*link.Link, error)
}

type Service struct {
	baseURL string
	repo    LinkRepository
	cache   LinkCache
}

func NewService(repo LinkRepository, cache LinkCache, baseURL string) *Service {
	return &Service{
		repo:    repo,
		cache:   cache,
		baseURL: baseURL,
	}
}

func (s *Service) Shorten(ctx context.Context, url, user string) (string, error) {
	if existing, found, err := s.cache.GetByUrl(ctx, url); err != nil {
		return "", err
	} else if found {
		return s.baseURL + "/" + existing.ShortCode, nil
	}

	shortCode, err := randomCode(7)
	if err != nil {
		return "", err
	}

	l := &link.Link{URL: url, ShortCode: shortCode, UserID: user}

	savedLink, err := s.repo.Save(ctx, l)
	if err != nil {
		return "", err
	}

	if err := s.cache.Save(ctx, savedLink); err != nil {
		log.Printf("WARN: could not save link to cache: %v", err)
	}

	return s.baseURL + "/" + shortCode, nil
}

func (s *Service) Resolve(ctx context.Context, shortCode string) (string, error) {
	l, found, err := s.cache.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	if !found {
		l, found, err = s.repo.FindByShortCode(ctx, shortCode)
		if err != nil {
			return "", err
		}

		if !found {
			return "", ErrShortLinkNotFound
		}

		if err := s.cache.Save(ctx, l); err != nil {
			log.Printf("WARN: could not save link to cache: %v", err)
		}
	}

	s.cache.IncrementClickCount(ctx, l.ShortCode)

	return l.URL, nil
}

func (s *Service) GetByUser(ctx context.Context, userID string, limit, skip int32) ([]*link.Link, error) {
	links, err := s.repo.GetByUser(ctx, userID, limit, skip)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func randomCode(n int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var out []byte
	for i := 0; i < n; i++ {
		x, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		out = append(out, alphabet[x.Int64()])
	}
	return string(out), nil
}
