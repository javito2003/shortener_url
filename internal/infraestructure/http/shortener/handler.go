package shortener

import (
	"github.com/gin-gonic/gin"
	"github.com/javito2003/shortener_url/internal/app/shortener"
)

type Handler struct {
	shortenerService shortener.Shortener
}

func NewHandler(shortener shortener.Shortener) *Handler {
	return &Handler{shortenerService: shortener}
}

func (h *Handler) createShortener() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateShortenerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request payload"})
			return
		}

		link, err := h.shortenerService.Shorten(c.Request.Context(), req.URL)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create shortener"})
			return
		}

		c.JSON(200, gin.H{"message": "createShortener", "data": link})
	}
}
