package redirection

import (
	"context"
	"fmt"
	"sync"

	"github.com/emanuelefalzone/bitly/internal"

	"golang.org/x/exp/maps"
)

// InMemoryRepository is a redirection repository that store values in memory
type InMemoryRepository struct {
	mu         sync.Mutex
	collection map[string]*Redirection
}

// NewInMemoryRepository creates a new redirection repository that store values in memory
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		collection: make(map[string]*Redirection),
	}
}

func (repo *InMemoryRepository) Create(ctx context.Context, value *Redirection) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	// Check if the Key already exists
	if _, alreadyExists := repo.collection[value.Key]; alreadyExists {
		// Cannot create a redirection with the same Key return error
		return &internal.Error{
			Code:    internal.ErrConflict,
			Message: fmt.Sprintf("a redirection with key %s already exists", value.Key),
		}
	}
	// Store the new redirect
	repo.collection[value.Key] = value

	return nil
}

func (repo *InMemoryRepository) Delete(ctx context.Context, value *Redirection) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	// Check if the Key already exists
	if _, alreadyExists := repo.collection[value.Key]; !alreadyExists {
		// Cannot delete a redirection that does not exists, return error
		return &internal.Error{
			Code:    internal.ErrNotFound,
			Message: fmt.Sprintf("cannot find a redirection with key %s", value.Key),
		}
	}
	// Delete the specified redirection
	delete(repo.collection, value.Key)

	return nil
}

func (repo *InMemoryRepository) FindByKey(ctx context.Context, key string) (*Redirection, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	// Check if the key exists
	if value, exists := repo.collection[key]; exists {
		return value, nil
	}

	// Cannot find a redirection with the given key return error
	return nil, &internal.Error{
		Code:    internal.ErrNotFound,
		Message: fmt.Sprintf("cannot find a redirection with key %s", key),
	}
}

func (repo *InMemoryRepository) FindAll(ctx context.Context) ([]*Redirection, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	// Return all the redirection values
	return maps.Values(repo.collection), nil
}
