package internal

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/javito2003/shortener_url/internal/app/clicks_worker"
	appShortener "github.com/javito2003/shortener_url/internal/app/shortener"
	"github.com/javito2003/shortener_url/internal/config"
	"github.com/javito2003/shortener_url/internal/infrastructure/http"
	"github.com/javito2003/shortener_url/internal/infrastructure/http/shortener"
	"github.com/javito2003/shortener_url/internal/infrastructure/persistence/mongo"
	"github.com/javito2003/shortener_url/internal/infrastructure/persistence/redis"
)

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
	mongoRepo := mongo.NewRepository(database)
	redisCache := redis.NewStore(client)
	clicksReader := redis.NewClicksReader(client)
	bulkUpdater := mongo.NewLinkBulkUpdater(database)

	// Service
	shortenerService := appShortener.NewService(mongoRepo, redisCache, config.AppConfig.BaseURL)
	workerService := clicks_worker.NewService(clicksReader, bulkUpdater)

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
	shortener.NewRouter(httpServer, shortenerService)

	return server
}

func (s *Server) Run() error {
	return s.http.Run(":8080")
}
