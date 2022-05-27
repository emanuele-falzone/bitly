package redis

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/go-redis/redis/v8"
)

type RedirectionRepository struct {
	client *redis.Client
}

func NewRedirectionRepository(connectionString string) (*RedirectionRepository, error) {
	// Parse connection string and check for errors
	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, &internal.Error{
			Code: internal.ErrInvalid,
			Op:   "NewRedirectionRepository",
			Err:  err,
		}
	}

	// Create a new redis client
	client := redis.NewClient(opt)

	// Return a new RedirectionRepository
	return &RedirectionRepository{client: client}, nil
}

func (r *RedirectionRepository) Create(ctx context.Context, a redirection.Redirection) error {
	// Save the redirection in Redis
	ok, err := r.client.SetNX(ctx, a.Key, a.Location, 0).Result()

	if err != nil {
		// There was some problem with Redis return error
		return &internal.Error{
			Code: internal.ErrInternal,
			Op:   "RedirectionRepository: Create",
			Err:  err,
		}
	}

	if !ok {
		// A redirection with the same key already exists return error
		return &internal.Error{
			Code: internal.ErrConflict,
			Op:   "RedirectionRepository: Create",
		}
	}

	// Return nil to signal that the creation was executed
	return nil
}

func (r *RedirectionRepository) Delete(ctx context.Context, a redirection.Redirection) error {
	// Delete the redirection from Redis
	_, err := r.client.GetDel(ctx, a.Key).Result()

	if err == redis.Nil {
		// Cannot delete a redirection that does not exists return ErrNotFound
		return &internal.Error{
			Code: internal.ErrNotFound,
			Op:   "RedirectionRepository: Delete",
		}
	}

	if err != nil {
		// There was some problem with Redis return error
		return &internal.Error{
			Code: internal.ErrInternal,
			Op:   "RedirectionRepository: Delete",
			Err:  err,
		}
	}

	// Return nil to signal that the deletion was executed
	return nil
}

func (r *RedirectionRepository) FindByKey(ctx context.Context, key string) (redirection.Redirection, error) {
	// Get the location associated with the key
	location, err := r.client.Get(ctx, key).Result()

	if err == redis.Nil {
		// Cannot find a redirection that does not exists return ErrNotFound
		return redirection.Redirection{}, &internal.Error{
			Code: internal.ErrNotFound,
			Op:   "RedirectionRepository: FindByKey",
		}
	}

	if err != nil {
		// There was some problem with Redis return error
		return redirection.Redirection{}, &internal.Error{
			Code: internal.ErrInternal,
			Op:   "RedirectionRepository: FindByKey",
			Err:  err,
		}
	}

	// Use key and location to create a new redirection
	return redirection.New(key, location)
}

func (r *RedirectionRepository) FindAll(ctx context.Context) ([]redirection.Redirection, error) {
	// Get all the keys from redis
	keys, err := r.client.Keys(ctx, "*").Result()

	if err != nil {
		// There was some problem with Redis return error
		return nil, &internal.Error{
			Code: internal.ErrInternal,
			Op:   "RedirectionRepository: FindAll",
			Err:  err,
		}
	}

	// Create new empty result set
	result := []redirection.Redirection{}

	// Iterate over keys
	for _, key := range keys {
		// Find the value associated with the specified key
		value, err := r.FindByKey(ctx, key)
		if err != nil {
			// There was some problem with Redis return error
			return nil, &internal.Error{
				Code: internal.ErrInternal,
				Op:   "RedirectionRepository: FindAll",
				Err:  err,
			}
		}
		// Append the value to the result
		result = append(result, value)
	}

	// Return result
	return result, nil
}
