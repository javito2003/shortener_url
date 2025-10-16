package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func Connect() (*redis.Client, error) {
	if client != nil {
		return client, nil
	}

	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := client.Ping(context.Background()).Result()

	return client, err
}
