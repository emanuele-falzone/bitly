package application

import (
	"context"
	"log"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/event"
	"github.com/emanuelefalzone/bitly/internal/application/redirection"
)

// CreateRedirection creates a new redirection for the given location and returns the associated key
func (app *Application) CreateRedirection(ctx context.Context, location string) (key string, err error) {
	// Get a new key from the key generator
	key = app.keyGenerator.NextKey()

	// Create a new redirection given the generated key and specified location
	value, err := redirection.New(key, location)

	// If the create operation fails return error
	if err != nil {
		return "", &internal.Error{
			Op:  "Application: CreateRedirection",
			Err: err,
		}
	}

	// Save the redirection inside the repository
	if err := app.redirectionRepository.Create(ctx, value); err != nil {
		return "", &internal.Error{
			Op:  "Application: CreateRedirection",
			Err: err,
		}
	}

	// Create new event
	e := event.Now(event.TypeCreate, value)

	// If the save operation fails return error
	if err := app.eventRepository.Create(ctx, e); err != nil {
		return "", &internal.Error{
			Op:  "Application: CreateRedirection",
			Err: err,
		}
	}

	// Log event to console
	log.Printf("Key: %s, Location: %s, Event: %s, DateTime: %s\n",
		e.Redirection.Key,
		e.Redirection.Location,
		e.Type,
		e.DateTime)

	// Return the key of the newly created redirection
	return key, nil
}
