package command

import (
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
)

type DeleteRedirectionCommand struct {
	Key string
}

type DeleteRedirectionHandler struct {
	redirections redirection.Repository
}

func NewDeleteRedirectionHandler(redirections redirection.Repository) DeleteRedirectionHandler {
	return DeleteRedirectionHandler{redirections: redirections}
}

func (h DeleteRedirectionHandler) Handle(cmd DeleteRedirectionCommand) error {
	// Find the redirection inside the repository
	val, err := h.redirections.FindByKey(cmd.Key)

	// If the find operation fails return error
	if err != nil {
		return &internal.Error{Op: "DeleteRedirectionHandler: Handle", Err: err}
	}

	// Save the redirection insire the repository
	err = h.redirections.Delete(val)

	// If the delete operation fails return error
	if err != nil {
		return &internal.Error{Op: "DeleteRedirectionHandler: Handle", Err: err}
	}

	// Return nil to indicate that the command was succesfully executed
	return nil
}
