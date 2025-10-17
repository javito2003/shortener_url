package redis

import (
	"context"
	"strconv" // Required for string to int conversion
	"strings"

	"github.com/redis/go-redis/v9"
)

type ClicksReader struct {
	client *redis.Client
}

func NewClicksReader(client *redis.Client) *ClicksReader {
	return &ClicksReader{client: client}
}

func (r *ClicksReader) FetchAndClear(ctx context.Context) (map[string]int64, error) {
	counts := make(map[string]int64)
	pattern := "clicks:short:*" // Asumo que tus contadores usan INCR y no HINCRBY

	// 1. Primero, obtenemos todas las claves que coinciden con el patrón
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return counts, nil
	}

	pipe := r.client.Pipeline()

	getCmds := make([]*redis.StringCmd, len(keys))
	for i, key := range keys {
		getCmds[i] = pipe.HGet(ctx, key, "click_count")
	}

	pipe.Del(ctx, keys...)

	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	for i, cmd := range getCmds {
		valStr, err := cmd.Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			continue
		}

		count, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			continue
		}

		shortCode := strings.TrimPrefix(keys[i], "clicks:short:")
		counts[shortCode] = count
	}

	return counts, nil
}
