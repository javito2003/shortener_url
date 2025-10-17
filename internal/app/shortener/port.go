package shortener

import (
	"context"

	link "github.com/javito2003/shortener_url/internal/domain"
)

type LinkRepository interface {
	Save(ctx context.Context, link *link.Link) (*link.Link, error)
	FindByShortCode(ctx context.Context, shortCode string) (*link.Link, bool, error)
	GetByUrl(ctx context.Context, url string) (*link.Link, bool, error)
	GetByUser(ctx context.Context, userID string, limit, skip int32) ([]*link.Link, error)
}

type LinkCache interface {
	Save(ctx context.Context, link *link.Link) error
	FindByShortCode(ctx context.Context, shortCode string) (*link.Link, bool, error)
	GetByUrl(ctx context.Context, url string) (*link.Link, bool, error)
	IncrementClickCount(ctx context.Context, linkId string) error
}
