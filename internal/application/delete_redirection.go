package application

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
)

type DeleteRedirectionCommand struct {
	Key string
}

type DeleteRedirectionHandler struct {
	redirections redirection.Repository
	dispatcher   *event.Dispatcher
}

func NewDeleteRedirectionHandler(redirections redirection.Repository, dispatcher *event.Dispatcher) DeleteRedirectionHandler {
	return DeleteRedirectionHandler{redirections: redirections, dispatcher: dispatcher}
}

func (h DeleteRedirectionHandler) Handle(ctx context.Context, cmd DeleteRedirectionCommand) error {
	// Find the redirection inside the repository
	val, err := h.redirections.FindByKey(ctx, cmd.Key)

	// If the find operation fails return error
	if err != nil {
		return &internal.Error{Op: "DeleteRedirectionHandler: Handle", Err: err}
	}

	// Save the redirection inside the repository
	err = h.redirections.Delete(ctx, val)

	// If the delete operation fails return error
	if err != nil {
		return &internal.Error{Op: "DeleteRedirectionHandler: Handle", Err: err}
	}

	// Dispatch deleted event
	h.dispatcher.Dispatch(ctx, event.Deleted(val))

	// Return nil to indicate that the command was successfully executed
	return nil
}
