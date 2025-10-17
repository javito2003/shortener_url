package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/javito2003/shortener_url/internal/app/auth"
	"github.com/javito2003/shortener_url/internal/infrastructure/http/user"
)

type Handler struct {
	userService auth.AuthService
}

func NewHandler(userService auth.AuthService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) registerUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var bodyReq registerUserDTO
		if err := c.ShouldBindJSON(&bodyReq); err != nil {
			c.Error(err)
			return
		}

		userCreated, err := h.userService.CreateUser(
			c.Request.Context(),
			bodyReq.FirstName,
			bodyReq.LastName,
			bodyReq.Email,
			bodyReq.Password,
		)

		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, gin.H{"message": "user registered", "data": user.ToUserResponse(userCreated)})
	}
}

func (h *Handler) authenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var bodyReq loginUserDTO
		if err := c.ShouldBindJSON(&bodyReq); err != nil {
			c.Error(err)
			return
		}

		token, err := h.userService.Authenticate(
			c.Request.Context(),
			bodyReq.Email,
			bodyReq.Password,
		)

		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, gin.H{"message": "user authenticated", "token": token})
	}
}
