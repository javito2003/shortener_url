package redis

import (
	"context"

	"github.com/javito2003/shortener_url/internal/config"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func Connect() (*redis.Client, error) {
	if client != nil {
		return client, nil
	}

	client = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Redis.Address,
		Username: config.AppConfig.Redis.Username,
		Password: config.AppConfig.Redis.Password,
	})

	_, err := client.Ping(context.Background()).Result()

	return client, err
}
