package application

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
)

// GetRedirectionList returns the key list of the available redirections
func (app *Application) GetRedirectionList(ctx context.Context) (keys []string, err error) {
	// Find the redirections inside the repository
	values, err := app.redirectionRepository.FindAll(ctx)

	// If the find operation fails return error
	if err != nil {
		return nil, &internal.Error{
			Op:  "Application: GetRedirectionList",
			Err: err,
		}
	}

	// Initialize keys
	keys = make([]string, len(values))

	// Extract keys from values
	for i, value := range values {
		keys[i] = value.Key
	}

	// Return the redirection keys
	return keys, nil
}
