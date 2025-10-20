package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/javito2003/shortener_url/internal/app"
	"github.com/javito2003/shortener_url/internal/app/auth"
	"github.com/javito2003/shortener_url/internal/app/user/mocks"
	"github.com/javito2003/shortener_url/internal/domain"
	"github.com/javito2003/shortener_url/internal/infrastructure/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

// --- 2. Función de Setup (Ayudante) ---

// setupService crea una instancia del servicio con dependencias reales (excepto el repo).
func setupService() (*auth.Service, *mocks.UserRepository) {
	mockUserRepo := new(mocks.UserRepository)

	// Usamos las implementaciones REALES de cómputo para un test más confiable.
	realHasher := security.NewBcryptHasher()
	realTokenizer := security.NewJWTGenerator("un-secreto-de-prueba-muy-seguro")

	// Creamos el servicio a probar
	authSvc := auth.NewService(mockUserRepo, realHasher, realTokenizer)

	return authSvc, mockUserRepo
}

// --- 3. Los Tests ---

func TestAuthService_Login(t *testing.T) {
	assert := assert.New(t)

	// --- Caso 1: Login Exitoso ---
	t.Run("should login successfully with correct credentials", func(t *testing.T) {
		// Arrange
		authSvc, mockUserRepo := setupService()
		realHasher := security.NewBcryptHasher()
		password := "password123"
		hashedPassword, _ := realHasher.Hash(password)

		mockUser := &domain.User{
			ID:        "user-id-123",
			Email:     "test@example.com",
			Password:  hashedPassword,
			FirstName: "javito",
			LastName:  "moreno",
		}

		// Configuramos el mock: si se busca este email, se devuelve el mockUser
		mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(mockUser, true, nil).Once()

		// Act
		token, err := authSvc.Authenticate(context.Background(), "test@example.com", password)

		// Assert
		assert.NoError(err)
		assert.NotNil(token)
		assert.NotEmpty(token)             // Verifica que se emitió un token no vacío
		mockUserRepo.AssertExpectations(t) // Verifica que FindByEmail fue llamado
	})

	// --- Caso 2: Usuario No Encontrado ---
	t.Run("should return unauthorized error if user not found", func(t *testing.T) {
		// Arrange
		authSvc, mockUserRepo := setupService()

		// Configuramos el mock: si se busca este email, se devuelve un error NotFound
		mockUserRepo.On("FindByEmail", mock.Anything, "notfound@example.com").
			Return(nil, false, app.NewNotFoundError("user")).Once()

		// Act
		token, err := authSvc.Authenticate(context.Background(), "notfound@example.com", "password123")

		// Assert
		assert.Error(err)
		assert.Empty(token)

		var appErr *app.AppError
		assert.True(errors.As(err, &appErr), "error should be an AppError")
		assert.Equal(app.ErrUnauthorized, appErr.Type, "error type should be Unauthorized")
		assert.Equal("Invalid credentials", appErr.Message)
		mockUserRepo.AssertExpectations(t)
	})

	// --- Caso 3: Contraseña Incorrecta ---
	t.Run("should return unauthorized error on wrong password", func(t *testing.T) {
		// Arrange
		authSvc, mockUserRepo := setupService()
		realHasher := security.NewBcryptHasher()
		password := "password123"
		hashedPassword, _ := realHasher.Hash(password)

		mockUser := &domain.User{
			ID:        "user-id-123",
			Email:     "test@example.com",
			Password:  hashedPassword,
			FirstName: "javito",
			LastName:  "moreno",
		}

		mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(mockUser, true, nil).Once()

		// Act
		token, err := authSvc.Authenticate(context.Background(), "test@example.com", "WRONG-password")

		// Assert
		assert.Error(err)
		assert.Empty(token)

		var appErr *app.AppError
		assert.True(errors.As(err, &appErr), "error should be an AppError")
		assert.Equal(app.ErrUnauthorized, appErr.Type)
		assert.Equal("Invalid credentials", appErr.Message)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthService_Register(t *testing.T) {
	assert := assert.New(t)

	// --- Caso 1: Login Exitoso ---
	t.Run("should register successfully with correct credentials", func(t *testing.T) {
		// Arrange
		authSvc, mockUserRepo := setupService()
		realHasher := security.NewBcryptHasher()

		firstName := "javito"
		lastName := "moreno"
		email := "test@example.com"
		password := "password123"
		hashedPassword, _ := realHasher.Hash(password)

		mockUser := &domain.User{
			ID:        "user-id-123",
			Email:     email,
			Password:  hashedPassword,
			FirstName: firstName,
			LastName:  lastName,
		}

		// Configuramos el mock: si se busca este email, se devuelve el mockUser
		mockUserRepo.On("FindByEmail", mock.Anything, email).Return(nil, false, nil).Once()
		mockUserRepo.On("Save", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
			return user.Email == email && user.FirstName == firstName && user.LastName == lastName && realHasher.Compare(user.Password, password)
		})).Return(mockUser, nil).Once()

		// Act
		user, err := authSvc.CreateUser(context.Background(), firstName, lastName, email, password)

		// Assert
		assert.NoError(err)
		assert.NotNil(user)
		assert.NotEmpty(user.ID)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return email already in use error", func(t *testing.T) {
		// Arrange
		authSvc, mockUserRepo := setupService()
		realHasher := security.NewBcryptHasher()

		firstName := "javito"
		lastName := "moreno"
		email := "test@example.com"
		password := "password123"
		hashedPassword, _ := realHasher.Hash(password)

		mockUser := &domain.User{
			ID:        "user-id-123",
			Email:     email,
			Password:  hashedPassword,
			FirstName: firstName,
			LastName:  lastName,
		}

		// Configuramos el mock: si se busca este email, se devuelve el mockUser
		mockUserRepo.On("FindByEmail", mock.Anything, email).Return(mockUser, true, nil).Once()
		// Act
		user, err := authSvc.CreateUser(context.Background(), firstName, lastName, email, password)

		// Assert
		var appErr *app.AppError
		assert.True(errors.As(err, &appErr), "error should be an AppError")
		assert.Equal(app.ErrConflict, appErr.Type, "error type should be Conflict")
		assert.Nil(user)

		mockUserRepo.AssertExpectations(t)
	})
}
