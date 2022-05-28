package application

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/event"
)

func (app *Application) GetRedirectionCount(ctx context.Context, key string) (int, error) {
	// Find the redirection inside the repository
	val, err := app.redirectionRepository.FindByKey(ctx, key)

	// If the find operation fails return error
	if err != nil {
		return 0, &internal.Error{
			Op:  "Application: GetRedirectionCount",
			Err: err,
		}
	}

	// Find the redirection inside the repository
	events, err := app.eventRepository.FindByRedirection(ctx, val)

	// If the find operation fails return error
	if err != nil {
		return 0, &internal.Error{
			Op:  "Application: GetRedirectionCount",
			Err: err,
		}
	}

	// Computer the number of read events by iterating the events
	readEventCount := 0

	for _, e := range events {
		if e.Type == event.TypeRead {
			readEventCount++
		}
	}

	// Return the number of times a specific redirection has been read
	return readEventCount, nil
}
