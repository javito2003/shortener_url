package shortener_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/javito2003/shortener_url/internal/app/shortener"
	"github.com/javito2003/shortener_url/internal/app/shortener/mocks"
	"github.com/javito2003/shortener_url/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const baseUrl = "http://short.ly"

func setupService() (shortener.Shortener, *mocks.LinkRepository, *mocks.LinkCache) {
	mockShortenerRepo := new(mocks.LinkRepository)
	mockCache := new(mocks.LinkCache)
	service := shortener.NewService(mockShortenerRepo, mockCache, baseUrl)

	return service, mockShortenerRepo, mockCache
}

func TestShortenerService_Shorten(t *testing.T) {
	assert := assert.New(t)
	service, mockRepo, mockCache := setupService()

	t.Run("Should return existing short URL from cache", func(t *testing.T) {
		ctx := context.Background()

		shortCode := "abc1234"
		mockCache.On("GetByUrl", mock.Anything, "http://example.com").Return(&domain.Link{
			ID:         "id",
			ClickCount: 0,
			UserID:     "userId",
			URL:        "http://example.com",
			ShortCode:  shortCode,
		}, true, nil).Once()

		shortURL, err := service.Shorten(ctx, "http://example.com", "userId")

		assert.NoError(err)
		assert.Equal(fmt.Sprintf("%s/%s", baseUrl, shortCode), shortURL)
	})

	t.Run("Should create new short URL when not in cache", func(t *testing.T) {
		ctx := context.Background()
		shortCode := "def5678"

		mockCache.On("GetByUrl", mock.Anything, "http://newurl.com").Return(&domain.Link{}, false, nil).Once()

		mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(l *domain.Link) bool {
			return l.URL == "http://newurl.com" && l.UserID == "userId"
		})).Return(func(ctx context.Context, l *domain.Link) *domain.Link {
			l.ID = "new-id"
			l.ShortCode = shortCode
			return l
		}, nil).Once()

		mockCache.On("Save", mock.Anything, mock.MatchedBy(func(l *domain.Link) bool {
			return l.URL == "http://newurl.com" && l.UserID == "userId"
		})).Return(nil).Once()

		shortURL, err := service.Shorten(ctx, "http://newurl.com", "userId")

		assert.NoError(err)
		assert.Equal(fmt.Sprintf("%s/%s", baseUrl, shortCode), shortURL)
	})
}

func TestShortenerService_Resolve(t *testing.T) {
	assert := assert.New(t)
	service, mockRepo, mockCache := setupService()

	t.Run("Should return URL from cache if found", func(t *testing.T) {
		ctx := context.Background()
		shortCode := "abc1234"
		expectedURL := "http://example.com"

		mockCache.On("FindByShortCode", mock.Anything, shortCode).Return(&domain.Link{
			ID:         "id",
			ClickCount: 0,
			UserID:     "userId",
			URL:        expectedURL,
			ShortCode:  shortCode,
		}, true, nil).Once()

		mockCache.On("IncrementClickCount", mock.Anything, shortCode).Return(nil).Once()

		url, err := service.Resolve(ctx, shortCode)

		assert.NoError(err)
		assert.Equal(expectedURL, url)
	})

	t.Run("Should return error if short code not found", func(t *testing.T) {
		ctx := context.Background()
		shortCode := "nonexistent"

		mockCache.On("FindByShortCode", mock.Anything, shortCode).Return(&domain.Link{}, false, nil).Once()
		mockRepo.On("FindByShortCode", mock.Anything, shortCode).Return(&domain.Link{}, false, nil).Once()

		url, err := service.Resolve(ctx, shortCode)

		assert.Equal(shortener.ErrShortLinkNotFound, err)
		assert.Empty(url)
	})
}

func TestShortenerService_GetByUser(t *testing.T) {
	assert := assert.New(t)
	service, mockRepo, _ := setupService()

	t.Run("Should return links for user", func(t *testing.T) {
		ctx := context.Background()
		userID := "user123"
		limit := int32(10)
		skip := int32(0)

		expectedLinks := []*domain.Link{
			{ID: "1", URL: "http://example1.com", ShortCode: "code1", UserID: userID},
			{ID: "2", URL: "http://example2.com", ShortCode: "code2", UserID: userID},
		}

		mockRepo.On("GetByUser", mock.Anything, userID, limit, skip).Return(expectedLinks, nil).Once()

		links, err := service.GetByUser(ctx, userID, limit, skip)

		assert.NoError(err)
		assert.Equal(expectedLinks, links)
	})
}
