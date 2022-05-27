package application

import (
	"context"
	"log"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
)

func (app *Application) DeleteRedirection(ctx context.Context, key string) error {
	// Find the redirection inside the repository
	val, err := app.redirectionRepository.FindByKey(ctx, key)

	// If the find operation fails return error
	if err != nil {
		return &internal.Error{Op: "DeleteRedirectionHandler: Handle", Err: err}
	}

	// Save the redirection inside the repository
	err = app.redirectionRepository.Delete(ctx, val)

	// If the delete operation fails return error
	if err != nil {
		return &internal.Error{Op: "DeleteRedirectionHandler: Handle", Err: err}
	}

	// Create new event
	event := event.Deleted(val)

	// Store created event in repository
	err = app.eventRepository.Create(ctx, event)

	// If the save operation fails return error
	if err != nil {
		return &internal.Error{Op: "Application: CreateRedirection", Err: err}
	}

	// Log event to console
	log.Printf("Key: %s, Location: %s, Event: %s, DateTime: %s\n", event.Redirection.Key, event.Redirection.Location, event.Type, event.DateTime)

	// Return nil to indicate that the command was successfully executed
	return nil
}
