package internal

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/javito2003/shortener_url/internal/app/clicks_worker"
	appShortener "github.com/javito2003/shortener_url/internal/app/shortener"
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

	baseUrl := "http://localhost:8080"

	client, err := redis.Connect()
	if err != nil {
		panic(err)
	}
	database, err := mongo.ConnectDatabase()
	if err != nil {
		panic(err)
	}

	// Infrastructure layer initialization
	mongoRepo := mongo.NewRepository(database)
	redisCache := redis.NewStore(client)

	clicksReader := redis.NewClicksReader(client)
	bulkUpdater := mongo.NewLinkBulkUpdater(database)

	shortenerService := appShortener.NewService(mongoRepo, redisCache, baseUrl)

	workerService := clicks_worker.NewService(clicksReader, bulkUpdater)

	server := &Server{
		http:             httpServer,
		ShortenerService: shortenerService,
	}

	shortener.NewRouter(httpServer, shortenerService)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		workerService.Run(ctx, 10*time.Second)
	}()

	_ = cancel

	return server
}

func (s *Server) Run() error {
	return s.http.Run(":8080")
}
