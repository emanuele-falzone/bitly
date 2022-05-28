package application

import (
	"context"
	"log"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/event"
)

// GetRedirectionLocation returns the location associated with the given redirection key
func (app *Application) GetRedirectionLocation(ctx context.Context, key string) (location string, err error) {
	// Find the redirection inside the repository
	value, err := app.redirectionRepository.FindByKey(ctx, key)

	// If the find operation fails return error
	if err != nil {
		return "", &internal.Error{
			Op:  "Application: GetRedirectionLocation",
			Err: err,
		}
	}

	// Create new event
	e := event.Now(event.TypeRead, value)

	// Store created event in repository
	if err := app.eventRepository.Create(ctx, e); err != nil {
		return "", &internal.Error{
			Op:  "Application: GetRedirectionLocation",
			Err: err,
		}
	}

	// Log event to console
	log.Printf("Key: %s, Location: %s, Event: %s, DateTime: %s\n",
		e.Redirection.Key,
		e.Redirection.Location,
		e.Type,
		e.DateTime)

	// Return the location the specified redirection
	return value.Location, nil
}
