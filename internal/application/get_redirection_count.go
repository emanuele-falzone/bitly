package application

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/event"
)

// GetRedirectionCount returns the number of times the redirection with the given key ahs been read
func (app *Application) GetRedirectionCount(ctx context.Context, key string) (count int, err error) {
	// Find the redirection inside the repository
	value, err := app.redirectionRepository.FindByKey(ctx, key)

	// If the find operation fails return error
	if err != nil {
		return 0, &internal.Error{
			Op:  "Application: GetRedirectionCount",
			Err: err,
		}
	}

	// Find the redirection inside the repository
	events, err := app.eventRepository.FindByRedirection(ctx, value)

	// If the find operation fails return error
	if err != nil {
		return 0, &internal.Error{
			Op:  "Application: GetRedirectionCount",
			Err: err,
		}
	}

	// Computer the number of read events by iterating the events
	for _, e := range events {
		if e.Type == event.TypeRead {
			count++
		}
	}

	// Return the number of times a specific redirection has been read
	return count, nil
}
