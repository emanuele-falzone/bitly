package memory

import (
	"context"
	"fmt"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
)

type InMemoryEventRepository struct {
	events map[string][]event.Event
}

func NewEventRepository() event.Repository {
	return InMemoryEventRepository{
		events: make(map[string][]event.Event),
	}
}

func (r InMemoryEventRepository) Create(ctx context.Context, a event.Event) error {
	// Append the events to the event list
	r.events[a.Redirection.Key] = append(r.events[a.Redirection.Key], a)
	return nil
}

func (r InMemoryEventRepository) FindByRedirection(ctx context.Context, a redirection.Redirection) ([]event.Event, error) {
	// Check if the Key exists
	if events, exists := r.events[a.Key]; exists {
		// Return the associated events
		return events, nil
	}
	// Cannot find the specified redirection, return error
	return []event.Event{}, &internal.Error{Code: internal.ErrNotFound, Message: fmt.Sprintf("cannot find a redirection with key %s", a.Key)}
}
