package redis

import (
	"context"
	"time"

	link "github.com/javito2003/shortener_url/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

type UrlSaved struct {
	ID         string `redis:"id"`
	URL        string `redis:"url"`
	ShortCode  string `redis:"short_code"`
	ClickCount int64  `redis:"click_count"`
}

func NewStore(client *redis.Client) *Store {
	return &Store{
		client: client,
	}
}

func (s *Store) Save(ctx context.Context, link *link.Link) error {
	pipe := s.client.Pipeline()

	var ttl time.Duration
	if link.IsExpired() {
		return nil
	}

	if link.ExpiresAt != nil {
		ttl = time.Until(*link.ExpiresAt)

		if ttl <= 0 {
			return nil
		}
	}

	pipe.HSet(ctx, "short:"+link.ShortCode, UrlSaved{
		ID:         link.ID,
		URL:        link.URL,
		ShortCode:  link.ShortCode,
		ClickCount: int64(link.ClickCount),
	})
	if ttl > 0 {
		pipe.Expire(ctx, "short:"+link.ShortCode, ttl)
	}

	pipe.Set(ctx, "url:"+link.URL, link.ShortCode, ttl)

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
	var linkData UrlSaved
	err = s.client.HGetAll(ctx, "short:"+shortCode).Scan(&linkData)
	if err != nil {
		return nil, false, err
	}

	if linkData == (UrlSaved{}) {
		return nil, false, nil
	}

	return &link.Link{
		ID:         linkData.ID,
		URL:        linkData.URL,
		ShortCode:  linkData.ShortCode,
		ClickCount: int(linkData.ClickCount),
	}, true, nil
}

func (s *Store) IncrementClickCount(ctx context.Context, shortCode string) error {
	err := s.client.HIncrBy(ctx, "clicks:short:"+shortCode, "click_count", 1).Err()
	return err
}

// func (s *Store) IncrementHotness(ctx context.Context, shortCode string) error {
// 	s.client.HIncrBy(ctx, "clicks:hot:"+shortCode, "hotness", 1)
// 	return nil
// }

// func (s *Store) GetHotness(ctx context.Context, limit int64) ([]*ports.LinkData, error) {
// 	var links []*ports.LinkData

// 	var cursor uint64
// 	for {
// 		keys, nextCursor, err := s.client.Scan(ctx, cursor, "clicks:hot:*", 100).Result()
// 		if err != nil {
// 			return nil, err
// 		}

// 		fmt.Println(keys)

// 		cursor = nextCursor
// 		if cursor == 0 {
// 			break
// 		}
// 	}

// 	return links, nil
// }
