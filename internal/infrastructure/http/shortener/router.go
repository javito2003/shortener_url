package shortener

import (
	"github.com/gin-gonic/gin"
	"github.com/javito2003/shortener_url/internal/app/auth"
	"github.com/javito2003/shortener_url/internal/app/shortener"
	"github.com/javito2003/shortener_url/internal/infrastructure/http/middleware"
)

func NewRouter(router *gin.Engine, shortener shortener.Shortener, token auth.TokenGenerator) {
	handler := NewHandler(shortener)
	router.GET("/:shortCode", handler.resolveShortener())

	{
		shortenGroup := router.Group("/shorten").Use(middleware.AuthMiddleware(token))
		shortenGroup.POST("", handler.createShortener())
		shortenGroup.GET("/", handler.getByUser())
	}
}
