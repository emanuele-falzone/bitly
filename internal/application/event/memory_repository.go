package event

import (
	"context"
	"fmt"
	"sync"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/redirection"
)

// InMemoryRepository is an event repository that store values in memory
type InMemoryRepository struct {
	mu     sync.Mutex
	events map[string][]*Event
}

// NewInMemoryRepository creates an new event repository that store values in memory
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		events: make(map[string][]*Event),
	}
}

func (repo *InMemoryRepository) Create(ctx context.Context, value *Event) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	// Append the events to the event list
	repo.events[value.Redirection.Key] = append(repo.events[value.Redirection.Key], value)

	return nil
}

func (repo *InMemoryRepository) FindByRedirection(ctx context.Context,
	value *redirection.Redirection) ([]*Event, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	// Check if the Key exists
	if events, exists := repo.events[value.Key]; exists {
		// Return the associated events
		return events, nil
	}

	// Cannot find the specified redirection, return error
	return nil, &internal.Error{
		Code:    internal.ErrNotFound,
		Message: fmt.Sprintf("cannot find a redirection with key %s", value.Key),
	}
}
