package user_test

import (
	"context"
	"testing"

	"github.com/javito2003/shortener_url/internal/app/user"
	"github.com/javito2003/shortener_url/internal/app/user/mocks"
	"github.com/javito2003/shortener_url/internal/domain"
	"github.com/stretchr/testify/assert"
)

func setupService() (user.UserService, *mocks.UserRepository) {
	mockRepo := new(mocks.UserRepository)
	service := user.NewUserService(mockRepo)

	return service, mockRepo
}

func TestUserService_GetMe(t *testing.T) {
	service, mockRepo := setupService()
	assert := assert.New(t)

	t.Run("should return user when found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := "user-123"
		expectedUser := &domain.User{
			ID:    userID,
			Email: "user@example.com",
		}
		mockRepo.On("GetById", ctx, userID).Return(expectedUser, true, nil)

		// Act
		user, err := service.GetMe(ctx, userID)

		// Assert
		assert.NoError(err)
		assert.Equal(expectedUser, user)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := "non-existent-user"
		mockRepo.On("GetById", ctx, userID).Return(nil, false, nil)

		// Act
		userFound, err := service.GetMe(ctx, userID)

		// Assert
		assert.Nil(userFound)
		assert.Equal(user.ErrUserNotFound, err)
	})
}
