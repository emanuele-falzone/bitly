package query

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
)

type RedirectionCountQuery struct {
	Key string
}

type RedirectionCountQueryResult struct {
	Count int
}

type RedirectionCountHandler struct {
	redirections redirection.Repository
	events       event.Repository
}

func NewRedirectionCountHandler(redirections redirection.Repository, events event.Repository) RedirectionCountHandler {
	return RedirectionCountHandler{redirections: redirections, events: events}
}

func (h RedirectionCountHandler) Handle(ctx context.Context, query RedirectionCountQuery) (*RedirectionCountQueryResult, error) {
	// Find the redirection inside the repository
	val, err := h.redirections.FindByKey(ctx, query.Key)

	// If the find operation fails return error
	if err != nil {
		return nil, &internal.Error{Op: "RedirectionCountHandler: Handle", Err: err}
	}

	// Find the redirection inside the repository
	events, err := h.events.FindByRedirection(ctx, val)

	// If the find operation fails return error
	if err != nil {
		return nil, &internal.Error{Op: "RedirectionCountHandler: Handle", Err: err}
	}

	// Computer the number of read events by iterating the events
	readEventCount := 0
	for _, e := range events {
		if e.Type == event.TypeRead {
			readEventCount++
		}
	}

	// Return the number of times a specific redirection has been read
	return &RedirectionCountQueryResult{Count: readEventCount}, nil
}
