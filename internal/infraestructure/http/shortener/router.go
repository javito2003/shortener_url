package shortener

import (
	"github.com/gin-gonic/gin"
	"github.com/javito2003/shortener_url/internal/app/shortener"
)

func NewRouter(router *gin.Engine, shortener shortener.Shortener) {
	handler := NewHandler(shortener)

	{
		shortenGroup := router.Group("/shorten")
		shortenGroup.POST("", handler.createShortener())
	}
}
