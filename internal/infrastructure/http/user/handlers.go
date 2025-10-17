package user

import (
	"github.com/gin-gonic/gin"
	"github.com/javito2003/shortener_url/internal/app/user"
)

type Handler struct {
	userService user.UserService
}

func NewHandler(userService user.UserService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) getMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		user, err := h.userService.GetMe(c.Request.Context(), userID)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.JSON(200, gin.H{"message": "user fetched", "data": ToUserResponse(user)})
	}
}
