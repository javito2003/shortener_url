package shortener

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/javito2003/shortener_url/internal/app/ports"
	link "github.com/javito2003/shortener_url/internal/domain"
)

type Shortener interface {
	Shorten(ctx context.Context, url string) (string, error)
	Resolve(ctx context.Context, shortCode string) (string, error)
}

type Service struct {
	store   ports.Store
	baseUrl string
}

func NewService(store ports.Store, baseUrl string) *Service {
	return &Service{store: store, baseUrl: baseUrl}
}

func genShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}

	return string(result)
}

func genID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), genShortCode())
}

func (s *Service) Shorten(ctx context.Context, url string) (string, error) {
	existingLink, err := s.store.GetByUrl(ctx, url)
	if err != nil {
		return "", err
	}

	if existingLink != nil {
		return s.baseUrl + existingLink.ShortCode, nil
	}

	shortCode := genShortCode()
	for {
		existingLink, err := s.store.FindByShortCode(ctx, shortCode)

		if err != nil {
			return "", err
		}

		if existingLink == nil {
			break
		}

		shortCode = genShortCode()
	}

	err = s.store.Save(ctx, &link.Link{
		ID:        genID(),
		URL:       url,
		ShortCode: shortCode,
	})

	if err != nil {
		return "", err
	}

	return s.baseUrl + shortCode, nil
}

func (s *Service) Resolve(ctx context.Context, shortCode string) (string, error) {
	return "", nil
}
