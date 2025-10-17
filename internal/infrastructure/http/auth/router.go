package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/javito2003/shortener_url/internal/app/auth"
)

func NewRouter(router *gin.Engine, userService auth.AuthService) {
	handler := NewHandler(userService)
	{
		userGroup := router.Group("/auth")
		userGroup.POST("/register", handler.registerUser())
		userGroup.POST("/login", handler.authenticateUser())
	}
}
