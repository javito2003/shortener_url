package ports

import (
	"context"

	link "github.com/javito2003/shortener_url/internal/domain"
)

type Store interface {
	Save(ctx context.Context, link *link.Link) error
	FindByShortCode(ctx context.Context, shortCode string) (*link.Link, bool, error)
	GetByUrl(ctx context.Context, url string) (*link.Link, bool, error)
	IncrementClickCount(ctx context.Context, shortCode string) error
}
