package redis

import (
	"context"
	"strconv"

	link "github.com/javito2003/shortener_url/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

func mapUrlSaved(link *link.Link) map[string]interface{} {
	return map[string]interface{}{
		"id":          link.ID,
		"url":         link.URL,
		"short_code":  link.ShortCode,
		"click_count": link.ClickCount,
	}
}

func mapToLink(data map[string]string) *link.Link {
	var clickCount int64
	if data["click_count"] != "" {
		clickCount, _ = strconv.ParseInt(data["click_count"], 10, 64)
	}

	return &link.Link{
		ID:         data["id"],
		URL:        data["url"],
		ShortCode:  data["shortCode"],
		ClickCount: int(clickCount),
	}
}

func NewStore(client *redis.Client) *Store {
	return &Store{
		client: client,
	}
}

func (s *Store) Save(ctx context.Context, link *link.Link) error {
	pipe := s.client.Pipeline()

	pipe.HSet(ctx, "short:"+link.ShortCode, mapUrlSaved(link))
	pipe.Set(ctx, "url:"+link.URL, link.ShortCode, 0)

	_, err := pipe.Exec(ctx)

	return err
}

func (s *Store) GetByUrl(ctx context.Context, url string) (existentLink *link.Link, found bool, err error) {
	code, err := s.client.Get(ctx, "url:"+url).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	existentLink, found, err = s.FindByShortCode(ctx, code)
	return existentLink, found, err
}

func (s *Store) FindByShortCode(ctx context.Context, shortCode string) (existentLink *link.Link, found bool, err error) {
	val, err := s.client.HGetAll(ctx, "short:"+shortCode).Result()
	if err != nil {
		return nil, false, err
	}

	if len(val) == 0 {
		return nil, false, nil
	}
	linkData := mapToLink(val)
	return linkData, true, nil
}

func (s *Store) IncrementClickCount(ctx context.Context, shortCode string) error {
	s.client.HIncrBy(ctx, "short:"+shortCode, "click_count", 1)
	return nil
}
