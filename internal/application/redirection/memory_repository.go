package redirection

import (
	"context"
	"fmt"
	"sync"

	"github.com/emanuelefalzone/bitly/internal"

	"golang.org/x/exp/maps"
)

type InMemoryRepository struct {
	mu           sync.Mutex
	redirections map[string]Redirection
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		redirections: make(map[string]Redirection),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, a Redirection) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the Key already exists
	if _, alreadyExists := r.redirections[a.Key]; alreadyExists {
		// Cannot create a redirection with the same Key return error
		return &internal.Error{
			Code:    internal.ErrConflict,
			Message: fmt.Sprintf("a redirection with key %s already exists", a.Key),
		}
	}
	// Store the new redirect
	r.redirections[a.Key] = a

	return nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, a Redirection) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the Key already exists
	if _, alreadyExists := r.redirections[a.Key]; !alreadyExists {
		// Cannot delete a redirection that does not exists, return error
		return &internal.Error{
			Code:    internal.ErrNotFound,
			Message: fmt.Sprintf("cannot find a redirection with key %s", a.Key),
		}
	}
	// Delete the specified redirection
	delete(r.redirections, a.Key)

	return nil
}

func (r *InMemoryRepository) FindByKey(ctx context.Context, key string) (Redirection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the key exists
	if redirect, exists := r.redirections[key]; exists {
		return redirect, nil
	}

	// Cannot find a redirection with the given key return error
	return Redirection{}, &internal.Error{
		Code:    internal.ErrNotFound,
		Message: fmt.Sprintf("cannot find a redirection with key %s", key),
	}
}

func (r *InMemoryRepository) FindAll(ctx context.Context) ([]Redirection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Return all the redirection values
	return maps.Values(r.redirections), nil
}
