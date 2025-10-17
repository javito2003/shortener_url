package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/javito2003/shortener_url/internal/app/auth"
)

func AuthMiddleware(tokenGen auth.TokenGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(auth.ErrInvalidToken)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(auth.ErrInvalidToken)
			c.Abort()
			return
		}

		claims, err := tokenGen.Validate(c.Request.Context(), parts[1])
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		// Attach user ID to the context for other handlers
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
