package event

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal/application/redirection"
)

// Repository allows to store, delete and find events
type Repository interface {
	// Create stores the given event inside the repository
	Create(context.Context, *Event) error
	// FindByRedirection retrieves all the events related to given redirection
	FindByRedirection(context.Context, *redirection.Redirection) ([]*Event, error)
}
