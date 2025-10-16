package redis

import (
	"context"

	link "github.com/javito2003/shortener_url/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{
		client: client,
	}
}

func (s *Store) Save(ctx context.Context, link *link.Link) error {
	pipe := s.client.Pipeline()

	pipe.HSet(ctx, "short:"+link.ShortCode, link)
	pipe.HSet(ctx, "url:"+link.URL, link)

	_, err := pipe.Exec(ctx)

	return err
}

func (s *Store) GetByUrl(ctx context.Context, url string) (*link.Link, error) {
	var existentLink link.Link

	err := s.client.HGetAll(ctx, "url:"+url).Scan(&existentLink)
	if err != nil {
		return nil, err
	}

	if existentLink.URL == "" {
		return nil, nil
	}

	return &existentLink, nil
}

func (s *Store) FindByShortCode(ctx context.Context, shortCode string) (*link.Link, error) {
	var existentLink link.Link
	err := s.client.HGetAll(ctx, "short:"+shortCode).Scan(&existentLink)
	if err != nil {
		return nil, err
	}

	if existentLink.ShortCode == "" {
		return nil, nil
	}

	return &existentLink, nil
}

func (s *Store) IncrementClickCount(ctx context.Context, shortCode string) error {
	s.client.HIncrBy(ctx, "short:"+shortCode, "click_count", 1)
	return nil
}
