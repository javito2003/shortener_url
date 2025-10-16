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

	// Save by short code (key: shortCode, fields: url, shortCode)
	pipe.HSet(ctx, "short:"+link.ShortCode, map[string]interface{}{
		"url":       link.URL,
		"shortCode": link.ShortCode,
		"id":        link.ID,
	})

	// Save by URL for reverse lookup (key: url hash, fields: url, shortCode)
	pipe.HSet(ctx, "url:"+link.URL, map[string]interface{}{
		"url":       link.URL,
		"shortCode": link.ShortCode,
		"id":        link.ID,
	})

	_, err := pipe.Exec(ctx)

	return err
}

func (s *Store) GetByUrl(ctx context.Context, url string) (*link.Link, error) {
	result := s.client.HGetAll(ctx, "url:"+url)
	if result.Err() != nil {
		return nil, result.Err()
	}

	data := result.Val()
	if len(data) == 0 {
		return nil, nil
	}

	return &link.Link{
		ID:        data["id"],
		URL:       data["url"],
		ShortCode: data["shortCode"],
	}, nil
}

func (s *Store) FindByShortCode(ctx context.Context, shortCode string) (*link.Link, error) {
	result := s.client.HGetAll(ctx, "short:"+shortCode)
	if result.Err() != nil {
		return nil, result.Err()
	}

	data := result.Val()
	if len(data) == 0 {
		return nil, nil
	}

	return &link.Link{
		ID:        data["id"],
		URL:       data["url"],
		ShortCode: data["shortCode"],
	}, nil
}
