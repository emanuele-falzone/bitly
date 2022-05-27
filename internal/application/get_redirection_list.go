package application

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
)

func (app *Application) GetRedirectionList(ctx context.Context) ([]string, error) {
	// Find the redirections inside the repository
	values, err := app.redirectionRepository.FindAll(ctx)

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
	return keys, nil
}
