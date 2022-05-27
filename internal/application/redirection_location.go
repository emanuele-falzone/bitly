package application

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
)

type RedirectionLocationQuery struct {
	Key string
}

type RedirectionLocationQueryResult struct {
	Location string
}

type RedirectionLocationHandler struct {
	redirections redirection.Repository
	dispatcher   *event.Dispatcher
}

func NewRedirectionLocationHandler(redirections redirection.Repository, dispatcher *event.Dispatcher) RedirectionLocationHandler {
	return RedirectionLocationHandler{redirections: redirections, dispatcher: dispatcher}
}

func (h RedirectionLocationHandler) Handle(ctx context.Context, query RedirectionLocationQuery) (*RedirectionLocationQueryResult, error) {
	// Find the redirection inside the repository
	val, err := h.redirections.FindByKey(ctx, query.Key)

	// If the find operation fails return error
	if err != nil {
		return nil, &internal.Error{Op: "RedirectionLocationHandler: Handle", Err: err}
	}

	// Dispatch read event
	h.dispatcher.Dispatch(ctx, event.Read(val))

	// Return the location the specified redirection
	return &RedirectionLocationQueryResult{Location: val.Location}, nil
}
