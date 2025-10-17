package user

import (
	"github.com/gin-gonic/gin"
	"github.com/javito2003/shortener_url/internal/app/auth"
	"github.com/javito2003/shortener_url/internal/app/user"
	"github.com/javito2003/shortener_url/internal/infrastructure/http/middleware"
)

func NewRouter(router *gin.Engine, userService user.UserService, token auth.TokenGenerator) {
	handler := NewHandler(userService)

	{
		userGroup := router.Group("/users").Use(middleware.AuthMiddleware(token))
		userGroup.GET("/me", handler.getMe())
	}
}
