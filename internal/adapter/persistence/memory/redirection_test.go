package memory_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/memory"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRedirectionRepository_Create(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	repository := memory.NewRedirectionRepository()
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}

	// WHEN
	err := repository.Create(ctx, value)

	// THEN
	assert.Equal(t, nil, err)
}

func TestInMemoryRedirectionRepository_CreateFailed(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	repository := memory.NewRedirectionRepository()
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}
	repository.Create(ctx, value)

	// WHEN
	err := repository.Create(ctx, value)

	// THEN
	assert.Equal(t, internal.ErrConflict, internal.ErrorCode(err))
}

func TestInMemoryRedirectionRepository_Delete(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	repository := memory.NewRedirectionRepository()
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}
	repository.Create(ctx, value)

	// WHEN
	err := repository.Delete(ctx, value)

	// THEN
	assert.Equal(t, nil, err)
	value, err = repository.FindByKey(ctx, value.Key)
	assert.Equal(t, internal.ErrNotFound, internal.ErrorCode(err))
}

func TestInMemoryRedirectionRepository_DeleteFailed(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	repository := memory.NewRedirectionRepository()
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}

	// WHEN
	err := repository.Delete(ctx, value)

	// THEN
	assert.Equal(t, internal.ErrNotFound, internal.ErrorCode(err))

}

func TestInMemoryRedirectionRepository_FindByKey(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	repository := memory.NewRedirectionRepository()
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}
	repository.Create(ctx, value)

	// WHEN
	valueFound, err := repository.FindByKey(ctx, value.Key)

	// THEN
	assert.Equal(t, nil, err)
	assert.Equal(t, value.Key, valueFound.Key)
	assert.Equal(t, value.Location, valueFound.Location)
}

func TestInMemoryRedirectionRepository_FindByKeyFailed(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	repository := memory.NewRedirectionRepository()
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}

	// WHEN
	_, err := repository.FindByKey(ctx, value.Key)

	// THEN
	assert.Equal(t, internal.ErrNotFound, internal.ErrorCode(err))
}
