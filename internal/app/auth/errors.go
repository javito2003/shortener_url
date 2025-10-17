package auth

import "github.com/javito2003/shortener_url/internal/app"

var (
	ErrInvalidToken       = app.NewUnauthorizedError("Invalid token")
	ErrInvalidCredentials = app.NewUnauthorizedError("Invalid credentials")
	ErrAlreadyLoggedIn    = app.NewConflictError("User already exists")
)
