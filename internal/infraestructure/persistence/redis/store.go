package redis

import (
	"context"
	"fmt"

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

	fmt.Println(existentLink.ShortCode)

	if existentLink.URL == "" {
		return nil, nil
	}

	return &existentLink, nil
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
