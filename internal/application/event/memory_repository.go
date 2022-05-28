package event

import (
	"context"
	"fmt"
	"sync"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/redirection"
)

type InMemoryRepository struct {
	mu     sync.Mutex
	events map[string][]Event
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		events: make(map[string][]Event),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, a Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Append the events to the event list
	r.events[a.Redirection.Key] = append(r.events[a.Redirection.Key], a)

	return nil
}

func (r *InMemoryRepository) FindByRedirection(ctx context.Context,
	a redirection.Redirection) ([]Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the Key exists
	if events, exists := r.events[a.Key]; exists {
		// Return the associated events
		return events, nil
	}

	// Cannot find the specified redirection, return error
	return []Event{}, &internal.Error{
		Code:    internal.ErrNotFound,
		Message: fmt.Sprintf("cannot find a redirection with key %s", a.Key),
	}
}
