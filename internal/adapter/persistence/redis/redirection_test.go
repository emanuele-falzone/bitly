package redis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/redis"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_New(t *testing.T) {
	// WHEN
	_, err := redis.NewRedirectionRespository(fmt.Sprintf("http://localhost"))

	// THEN
	assert.Equal(t, internal.ErrInvalid, internal.ErrorCode(err))
}

func TestRedisRepository_Create(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	s, err := miniredis.Run()
	assert.Equal(t, nil, err)
	defer s.Close()
	repository, err := redis.NewRedirectionRespository(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Equal(t, nil, err)

	value := redirection.Redirection{Key: "abcdef", Location: "http://www.google.com"}

	// WHEN
	err = repository.Create(ctx, value)

	// THEN
	assert.Equal(t, nil, err)
}

func TestRedisRepository_CreateConflictErr(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	s, err := miniredis.Run()
	assert.Equal(t, nil, err)
	defer s.Close()
	repository, err := redis.NewRedirectionRespository(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Equal(t, nil, err)

	value := redirection.Redirection{Key: "abcdef", Location: "http://www.google.com"}
	value2 := redirection.Redirection{Key: "abcdef", Location: "http://www.apple.com"}

	err = repository.Create(ctx, value)
	assert.Equal(t, nil, err)

	// WHEN
	err = repository.Create(ctx, value2)

	// THEN
	assert.Equal(t, internal.ErrConflict, internal.ErrorCode(err))
}

func TestRedisRepository_CreateInternalErr(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	s, err := miniredis.Run()
	assert.Equal(t, nil, err)
	repository, err := redis.NewRedirectionRespository(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Equal(t, nil, err)

	value := redirection.Redirection{Key: "abcdef", Location: "http://www.google.com"}

	err = repository.Create(ctx, value)
	assert.Equal(t, nil, err)

	// WHEN
	s.Close()
	err = repository.Create(ctx, value)

	// THEN
	assert.Equal(t, internal.ErrInternal, internal.ErrorCode(err))
}

func TestRedisRepository_Delete(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	s, err := miniredis.Run()
	assert.Equal(t, nil, err)
	defer s.Close()
	repository, err := redis.NewRedirectionRespository(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Equal(t, nil, err)

	value := redirection.Redirection{Key: "abcdef", Location: "http://www.google.com"}

	err = repository.Create(ctx, value)
	assert.Equal(t, nil, err)

	// WHEN
	err = repository.Delete(ctx, value)

	// THEN
	assert.Equal(t, nil, err)
}

func TestRedisRepository_DeleteErrNotFound(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	s, err := miniredis.Run()
	assert.Equal(t, nil, err)
	defer s.Close()
	repository, err := redis.NewRedirectionRespository(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Equal(t, nil, err)

	value := redirection.Redirection{Key: "abcdef", Location: "http://www.google.com"}

	// WHEN
	err = repository.Delete(ctx, value)

	// THEN
	assert.Equal(t, internal.ErrNotFound, internal.ErrorCode(err))
}

func TestRedisRepository_DeleteErrInternal(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	s, err := miniredis.Run()
	assert.Equal(t, nil, err)
	repository, err := redis.NewRedirectionRespository(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Equal(t, nil, err)

	value := redirection.Redirection{Key: "abcdef", Location: "http://www.google.com"}

	// WHEN
	s.Close()
	err = repository.Delete(ctx, value)

	// THEN
	assert.Equal(t, internal.ErrInternal, internal.ErrorCode(err))
}

func TestRedisRepository_FindByKey(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	s, err := miniredis.Run()
	assert.Equal(t, nil, err)
	defer s.Close()
	repository, err := redis.NewRedirectionRespository(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Equal(t, nil, err)

	value := redirection.Redirection{Key: "abcdef", Location: "http://www.google.com"}

	err = repository.Create(ctx, value)
	assert.Equal(t, nil, err)

	// WHEN
	value, err = repository.FindByKey(ctx, value.Key)

	// THEN
	assert.Equal(t, nil, err)
	assert.Equal(t, value.Key, "abcdef")
	assert.Equal(t, value.Location, "http://www.google.com")
}

func TestRedisRepository_FindByKeyErrNotFound(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	s, err := miniredis.Run()
	assert.Equal(t, nil, err)
	defer s.Close()
	repository, err := redis.NewRedirectionRespository(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Equal(t, nil, err)

	value := redirection.Redirection{Key: "abcdef", Location: "http://www.google.com"}

	// WHEN
	value, err = repository.FindByKey(ctx, value.Key)

	// THEN
	assert.Equal(t, internal.ErrNotFound, internal.ErrorCode(err))
}

func TestRedisRepository_FindByKeyErrInternal(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	s, err := miniredis.Run()
	assert.Equal(t, nil, err)
	repository, err := redis.NewRedirectionRespository(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Equal(t, nil, err)

	value := redirection.Redirection{Key: "abcdef", Location: "http://www.google.com"}

	// WHEN
	s.Close()
	value, err = repository.FindByKey(ctx, value.Key)

	// THEN
	assert.Equal(t, internal.ErrInternal, internal.ErrorCode(err))
}
