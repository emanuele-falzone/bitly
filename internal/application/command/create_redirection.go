package command

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/emanuelefalzone/bitly/internal/service"
)

type CreateRedirectionCommand struct {
	Location string
}

type CreateRedirectionCommandResult struct {
	Key string
}

type CreateRedirectionHandler struct {
	redirections redirection.Repository
	generator    service.KeyGenerator
	dispatcher   *event.Dispatcher
}

func NewCreateRedirectionHandler(redirections redirection.Repository, generator service.KeyGenerator, dispatcher *event.Dispatcher) CreateRedirectionHandler {
	return CreateRedirectionHandler{redirections: redirections, generator: generator, dispatcher: dispatcher}
}

func (h CreateRedirectionHandler) Handle(ctx context.Context, cmd CreateRedirectionCommand) (*CreateRedirectionCommandResult, error) {
	// Get a new key from the key generator
	key := h.generator.NextKey(cmd.Location)

	// Create a new redirection given the generated key and specified location
	val, err := redirection.New(key, cmd.Location)

	// If the create operation fails return error
	if err != nil {
		return nil, &internal.Error{Op: "CreateRedirectionHandler: Handle", Err: err}
	}

	// Save the redirection inside the repository
	err = h.redirections.Create(ctx, val)

	// If the save operation fails return error
	if err != nil {
		return nil, &internal.Error{Op: "CreateRedirectionHandler: Handle", Err: err}
	}

	// Dispatch created event
	h.dispatcher.Dispatch(ctx, event.Created(val))

	// Return the key of the newly created redirection
	return &CreateRedirectionCommandResult{Key: key}, nil
}
