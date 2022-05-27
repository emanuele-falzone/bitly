package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"golang.org/x/exp/maps"
)

type InMemoryRedirectionRepository struct {
	mu           sync.Mutex
	redirections map[string]redirection.Redirection
}

func NewRedirectionRepository() *InMemoryRedirectionRepository {
	return &InMemoryRedirectionRepository{
		redirections: make(map[string]redirection.Redirection),
	}
}

func (r *InMemoryRedirectionRepository) Create(ctx context.Context, a redirection.Redirection) error {
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

func (r *InMemoryRedirectionRepository) Delete(ctx context.Context, a redirection.Redirection) error {
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

func (r *InMemoryRedirectionRepository) FindByKey(ctx context.Context, key string) (redirection.Redirection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the key exists
	if redirect, exists := r.redirections[key]; exists {
		return redirect, nil
	}

	// Cannot find a redirection with the given key return error
	return redirection.Redirection{}, &internal.Error{
		Code:    internal.ErrNotFound,
		Message: fmt.Sprintf("cannot find a redirection with key %s", key),
	}
}

func (r *InMemoryRedirectionRepository) FindAll(ctx context.Context) ([]redirection.Redirection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Return all the redirection values
	return maps.Values(r.redirections), nil
}
