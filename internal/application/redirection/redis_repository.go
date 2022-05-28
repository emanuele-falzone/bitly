package redirection

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/go-redis/redis/v8"
)

// RedisRepository is a redirection repository that store values in redis
type RedisRepository struct {
	client *redis.Client
}

// NewInMemoryRepository creates a new redirection repository that store values in redis
func NewRedisRepository(connection string) (*RedisRepository, error) {
	// Parse connection string and check for errors
	opt, err := redis.ParseURL(connection)
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
	return &RedisRepository{client: client}, nil
}

func (repo *RedisRepository) Create(ctx context.Context, value *Redirection) error {
	// Save the redirection in Redis
	ok, err := repo.client.SetNX(ctx, value.Key, value.Location, 0).Result()

	switch {
	case err != nil:
		// There was some problem with Redis return error
		return &internal.Error{
			Code: internal.ErrInternal,
			Op:   "RedisRepository: Create",
			Err:  err,
		}
	case !ok:
		// A redirection with the same key already exists return error
		return &internal.Error{
			Code: internal.ErrConflict,
			Op:   "RedisRepository: Create",
		}
	default:
		// Return nil to signal that the creation was executed
		return nil
	}
}

func (repo *RedisRepository) Delete(ctx context.Context, value *Redirection) error {
	// Delete the redirection from Redis
	_, err := repo.client.GetDel(ctx, value.Key).Result()

	switch {
	case err == redis.Nil:
		// Cannot delete a redirection that does not exists return ErrNotFound
		return &internal.Error{
			Code: internal.ErrNotFound,
			Op:   "RedisRepository: Delete",
		}
	case err != nil:
		// There was some problem with Redis return error
		return &internal.Error{
			Code: internal.ErrInternal,
			Op:   "RedisRepository: Delete",
			Err:  err,
		}
	default:
		// Return nil to signal that the deletion was executed
		return nil
	}
}

func (repo *RedisRepository) FindByKey(ctx context.Context, key string) (*Redirection, error) {
	// Get the location associated with the key
	location, err := repo.client.Get(ctx, key).Result()

	switch {
	case err == redis.Nil:
		// Cannot find a redirection that does not exists return ErrNotFound
		return nil, &internal.Error{
			Code: internal.ErrNotFound,
			Op:   "RedisRepository: FindByKey",
		}
	case err != nil:
		// There was some problem with Redis return error
		return nil, &internal.Error{
			Code: internal.ErrInternal,
			Op:   "RedisRepository: FindByKey",
			Err:  err,
		}
	default:
		// Use key and location to create a new redirection
		return New(key, location)
	}
}

func (repo *RedisRepository) FindAll(ctx context.Context) ([]*Redirection, error) {
	// Get all the keys from redis
	keys, err := repo.client.Keys(ctx, "*").Result()

	if err != nil {
		// There was some problem with Redis return error
		return nil, &internal.Error{
			Code: internal.ErrInternal,
			Op:   "RedisRepository: FindAll",
			Err:  err,
		}
	}

	// Create new empty result set
	result := make([]*Redirection, len(keys))

	// Iterate over keys
	for i, key := range keys {
		// Find the value associated with the specified key
		if result[i], err = repo.FindByKey(ctx, key); err != nil {
			// There was some problem with Redis return error
			return nil, &internal.Error{
				Code: internal.ErrInternal,
				Op:   "RedisRepository: FindAll",
				Err:  err,
			}
		}
	}

	// Return result
	return result, nil
}
