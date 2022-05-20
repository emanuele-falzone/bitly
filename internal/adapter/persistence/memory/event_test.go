//go:build unit

package memory_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/memory"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryEventRepository_Create(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	repository := memory.NewEventRepository()
	aRedirection := redirection.Redirection{Key: "short", Location: "http://www.google.com"}
	event := event.Created(aRedirection)

	// WHEN
	err := repository.Create(ctx, event)

	// THEN
	assert.Equal(t, nil, err)
}

func TestInMemoryEventRepository_FindByRedirection(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	repository := memory.NewEventRepository()
	aRedirection := redirection.Redirection{Key: "short", Location: "http://www.google.com"}
	CreateEvent := event.Created(aRedirection)
	readEvent1 := event.Read(aRedirection)
	readEvent2 := event.Read(aRedirection)
	readEvent3 := event.Read(aRedirection)
	deleteEvent := event.Deleted(aRedirection)

	err := repository.Create(ctx, CreateEvent)
	assert.Equal(t, nil, err)
	err = repository.Create(ctx, readEvent1)
	assert.Equal(t, nil, err)
	err = repository.Create(ctx, readEvent2)
	assert.Equal(t, nil, err)
	err = repository.Create(ctx, readEvent3)
	assert.Equal(t, nil, err)
	err = repository.Create(ctx, deleteEvent)
	assert.Equal(t, nil, err)

	// WHEN
	events, err := repository.FindByRedirection(ctx, aRedirection)

	// THEN
	assert.Equal(t, nil, err)
	assert.Equal(t, 5, len(events))
}

func TestInMemoryEventRepository_FindByRedirectionFailure(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	repository := memory.NewEventRepository()
	aRedirection := redirection.Redirection{Key: "short", Location: "http://www.google.com"}

	// WHEN
	events, err := repository.FindByRedirection(ctx, aRedirection)

	// THEN
	assert.Equal(t, 0, len(events))
	assert.Equal(t, internal.ErrNotFound, internal.ErrorCode(err))
}
