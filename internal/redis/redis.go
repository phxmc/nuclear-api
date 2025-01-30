package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context, addr, password string, db int) (*redis.Client, error) {
	options := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	}

	client := redis.NewClient(options)

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}
