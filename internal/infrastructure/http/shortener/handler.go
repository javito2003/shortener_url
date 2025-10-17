package shortener

import (
	"errors"
	"net/http"

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
			c.Error(err)
			return
		}

		userId := c.GetString("userID")

		link, err := h.shortenerService.Shorten(c.Request.Context(), req.URL, userId)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, gin.H{"message": "createShortener", "data": link})
	}
}

func (h *Handler) resolveShortener() gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")
		if shortCode == "" {
			c.Error(errors.New("missing short code"))
			return
		}

		url, err := h.shortenerService.Resolve(c.Request.Context(), shortCode)
		if err != nil {
			c.Error(err)
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (h *Handler) getByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		limit := int32(10)
		skip := int32(0)

		links, err := h.shortenerService.GetByUser(c.Request.Context(), userID, limit, skip)
		if err != nil {
			c.Error(err)
			return
		}

		var response []*linkResponse
		for _, l := range links {
			response = append(response, toLinkResponse(l))
		}

		c.JSON(200, gin.H{"message": "getByUser", "data": response})
	}
}
