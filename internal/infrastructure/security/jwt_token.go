package security

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/javito2003/shortener_url/internal/app/auth"
)

type JWTGenerator struct {
	secretKey []byte
}

func NewJWTGenerator(secret string) *JWTGenerator {
	return &JWTGenerator{secretKey: []byte(secret)}
}

type jwtClaims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

func (j *JWTGenerator) Generate(ctx context.Context, userID string, ttl *time.Duration) (string, error) {
	registeredClaims := jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}

	if ttl != nil {
		registeredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(*ttl))
	}

	claims := jwtClaims{
		UserID:           userID,
		RegisteredClaims: registeredClaims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTGenerator) Validate(ctx context.Context, tokenString string) (*auth.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil {
		return nil, auth.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		return &auth.Claims{
			UserID: claims.UserID,
		}, nil
	}

	return nil, auth.ErrInvalidToken
}

// Verificación estática para asegurar que la interfaz se cumple
var _ auth.TokenGenerator = (*JWTGenerator)(nil)
