package user

import "github.com/javito2003/shortener_url/internal/app"

var (
	ErrInvalidUserId = app.NewBadRequestError("invalid user id")
	ErrUserNotFound  = app.NewNotFoundError("user not found")
)
