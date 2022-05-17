package util

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func ClearRedis(ctx context.Context, connectionString string) error {
	// Parse connection string and check for errors
	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		return err
	}

	// Create a new redis client
	client := redis.NewClient(opt)

	// Flush all keys
	err = client.FlushAll(ctx).Err()

	return err
}
