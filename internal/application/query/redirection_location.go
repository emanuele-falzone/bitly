package query

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
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
}

func NewRedirectionLocationHandler(redirections redirection.Repository) RedirectionLocationHandler {
	return RedirectionLocationHandler{redirections: redirections}
}

func (h RedirectionLocationHandler) Handle(ctx context.Context, query RedirectionLocationQuery) (*RedirectionLocationQueryResult, error) {
	// Find the redirection inside the repository
	val, err := h.redirections.FindByKey(ctx, query.Key)

	// If the find operation fails return error
	if err != nil {
		return nil, &internal.Error{Op: "RedirectionLocationHandler: Handle", Err: err}
	}

	// Return the location the specified redirection
	return &RedirectionLocationQueryResult{Location: val.Location}, nil
}
