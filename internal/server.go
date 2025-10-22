package internal

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	appAuth "github.com/javito2003/shortener_url/internal/app/auth"
	"github.com/javito2003/shortener_url/internal/app/clicks_worker"
	appShortener "github.com/javito2003/shortener_url/internal/app/shortener"
	appUser "github.com/javito2003/shortener_url/internal/app/user"
	"github.com/javito2003/shortener_url/internal/config"
	"github.com/javito2003/shortener_url/internal/infrastructure/http"
	"github.com/javito2003/shortener_url/internal/infrastructure/http/auth"
	"github.com/javito2003/shortener_url/internal/infrastructure/http/shortener"
	httpUser "github.com/javito2003/shortener_url/internal/infrastructure/http/user"
	"github.com/javito2003/shortener_url/internal/infrastructure/persistence/mongo"
	"github.com/javito2003/shortener_url/internal/infrastructure/persistence/mongo/link"
	"github.com/javito2003/shortener_url/internal/infrastructure/persistence/mongo/user"
	"github.com/javito2003/shortener_url/internal/infrastructure/persistence/redis"
	"github.com/javito2003/shortener_url/internal/infrastructure/security"
)

const expiresMinutes = 60

type Server struct {
	http             *gin.Engine
	ShortenerService appShortener.Shortener
}

func NewServer() *Server {
	httpServer := gin.Default()

	// Clients
	client, err := redis.Connect()
	if err != nil {
		panic(err)
	}
	database, err := mongo.ConnectDatabase()
	if err != nil {
		panic(err)
	}

	// Repositories and Stores
	mongoRepo := link.NewRepository(database)
	userRepo := user.NewRepository(database)
	redisCache := redis.NewStore(client)
	clicksReader := redis.NewClicksReader(client)
	bulkUpdater := link.NewLinkBulkUpdater(database)
	hasher := security.NewBcryptHasher()
	token := security.NewJWTGenerator("secret")

	// Service
	shortenerService := appShortener.NewService(mongoRepo, redisCache, config.AppConfig.BaseURL, expiresMinutes)
	workerService := clicks_worker.NewService(clicksReader, bulkUpdater)
	authService := appAuth.NewService(userRepo, hasher, token)
	userService := appUser.NewUserService(userRepo)

	server := &Server{
		http:             httpServer,
		ShortenerService: shortenerService,
	}

	// Worker initialization
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		workerService.Run(ctx, 10*time.Second)
	}()

	_ = cancel

	// Router
	httpServer.Use(http.ErrorHandler())
	shortener.NewRouter(httpServer, shortenerService, token)
	auth.NewRouter(httpServer, authService)
	httpUser.NewRouter(httpServer, userService, token)

	return server
}

func (s *Server) Run() error {
	return s.http.Run(":8080")
}
