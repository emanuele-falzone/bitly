package command

import (
	"github.com/emanuelefalzone/bitly/internal"
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
}

func NewCreateRedirectionHandler(redirections redirection.Repository, generator service.KeyGenerator) CreateRedirectionHandler {
	return CreateRedirectionHandler{redirections: redirections, generator: generator}
}

func (h CreateRedirectionHandler) Handle(cmd CreateRedirectionCommand) (*CreateRedirectionCommandResult, error) {
	// Get a new key from the key generator
	key := h.generator.NextKey(cmd.Location)

	// Create a new redirection given the generated key and specified location
	val, err := redirection.New(key, cmd.Location)

	// If the create operation fails return error
	if err != nil {
		return nil, &internal.Error{Op: "CreateRedirectionHandler: Handle", Err: err}
	}

	// Save the redirection insire the repository
	err = h.redirections.Create(val)

	// If the save operation fails return error
	if err != nil {
		return nil, &internal.Error{Op: "CreateRedirectionHandler: Handle", Err: err}
	}

	// Return the key of the newly created redirection
	return &CreateRedirectionCommandResult{Key: key}, nil
}
