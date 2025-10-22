package shortener

import "github.com/javito2003/shortener_url/internal/app"

var (
	ErrShortLinkNotFound = app.NewNotFoundError("Short link not found")
	ErrShortLinkExpired  = app.NewBadRequestError("Short link has expired")
)
