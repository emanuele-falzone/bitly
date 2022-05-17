package redis

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/go-redis/redis/v8"
)

type RedisRedirectionRespository struct {
	client *redis.Client
}

func NewRedirectionRespository(connectionString string) (redirection.Repository, error) {
	// Parse connection string and check for errors
	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	// Create a new redis client
	client := redis.NewClient(opt)

	// Return a new RedisRedirectionRespository
	return RedisRedirectionRespository{client: client}, nil
}

func (r RedisRedirectionRespository) Create(ctx context.Context, a redirection.Redirection) error {
	// Save the redirection in Redis
	ok, err := r.client.SetNX(ctx, a.Key, a.Location, 0).Result()
	if err != nil {
		// There was some problem with Redis return error
		return &internal.Error{Code: internal.ErrInternal, Op: "RedisRedirectionRespository: Create", Err: err}
	}
	if !ok {
		// A redirection with the same key already exists return error
		return &internal.Error{Code: internal.ErrConflict, Op: "RedisRedirectionRespository: Create"}
	}

	// Return nil to signal that the creation was executed
	return nil
}

func (r RedisRedirectionRespository) Delete(ctx context.Context, a redirection.Redirection) error {
	// Delete the redirection from Redis
	_, err := r.client.GetDel(ctx, a.Key).Result()
	if err == redis.Nil {
		// Cannot delete a redirection that does not exists return ErrNotFound
		return &internal.Error{Code: internal.ErrNotFound, Op: "RedisRedirectionRespository: Delete"}
	}
	if err != nil {
		// There was some problem with Redis return error
		return &internal.Error{Code: internal.ErrInternal, Op: "RedisRedirectionRespository: Delete", Err: err}
	}

	// Return nil to signal that the deletion was executed
	return nil
}

func (r RedisRedirectionRespository) FindByKey(ctx context.Context, key string) (redirection.Redirection, error) {
	// Get the location associated with the key
	location, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Cannot find a redirection that does not exists return ErrNotFound
		return redirection.Redirection{}, &internal.Error{Code: internal.ErrNotFound, Op: "RedisRedirectionRespository: FindByKey"}
	}
	if err != nil {
		// There was some problem with Redis return error
		return redirection.Redirection{}, &internal.Error{Code: internal.ErrInternal, Op: "RedisRedirectionRespository: FindByKey", Err: err}
	}

	// Use key and location to create a new redirection
	return redirection.New(key, location)
}
