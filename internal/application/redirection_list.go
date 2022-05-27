package application

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
)

type RedirectionListQuery struct{}

type RedirectionListQueryResult struct {
	Keys []string
}

type RedirectionListHandler struct {
	redirections redirection.Repository
}

func NewRedirectionListHandler(redirections redirection.Repository) RedirectionListHandler {
	return RedirectionListHandler{redirections: redirections}
}

func (h RedirectionListHandler) Handle(ctx context.Context, query RedirectionListQuery) (*RedirectionListQueryResult, error) {
	// Find the redirections inside the repository
	values, err := h.redirections.FindAll(ctx)

	// If the find operation fails return error
	if err != nil {
		return nil, &internal.Error{Op: "RedirectionListHandler: Handle", Err: err}
	}

	// Create keys slice
	keys := []string{}

	// Extract keys from values
	for _, value := range values {
		keys = append(keys, value.Key)
	}

	// Return the redirection keys
	return &RedirectionListQueryResult{Keys: keys}, nil
}
