package internal

import (
	"github.com/gin-gonic/gin"
	appShortener "github.com/javito2003/shortener_url/internal/app/shortener"
	"github.com/javito2003/shortener_url/internal/infrastructure/http/shortener"
	"github.com/javito2003/shortener_url/internal/infrastructure/persistence/redis"
)

type Server struct {
	http             *gin.Engine
	ShortenerService appShortener.Shortener
}

const baseUrl = "http://localhost:8080/"

func NewServer() *Server {
	httpServer := gin.Default()
	client, err := redis.Connect()

	if err != nil {
		panic(err)
	}

	store := redis.NewStore(client)
	shortenerService := appShortener.NewService(store, baseUrl)
	server := &Server{
		http:             httpServer,
		ShortenerService: shortenerService,
	}

	shortener.NewRouter(httpServer, shortenerService)

	return server
}

func (s *Server) Run() error {
	return s.http.Run(":8080")
}
