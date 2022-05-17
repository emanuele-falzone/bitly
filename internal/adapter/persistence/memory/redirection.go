package memory

import (
	"fmt"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
)

type InMemoryRedirectionRepository struct {
	redirections map[string]redirection.Redirection
}

func NewRedirectionRepository() redirection.Repository {
	return InMemoryRedirectionRepository{redirections: make(map[string]redirection.Redirection)}
}

func (r InMemoryRedirectionRepository) Create(a redirection.Redirection) error {
	// Check if the Key already exists
	if _, alreadyExists := r.redirections[a.Key]; alreadyExists {
		// Cannot create a redirection with the same Key return error
		return &internal.Error{Code: internal.ErrConflict, Message: fmt.Sprintf("a redirection with key %s already exists", a.Key)}
	}
	// Store the new redirect
	r.redirections[a.Key] = a
	return nil
}

func (r InMemoryRedirectionRepository) Delete(a redirection.Redirection) error {
	// Check if the Key already exists
	if _, alreadyExists := r.redirections[a.Key]; !alreadyExists {
		// Cannot delete a redirection that does not exists, return error
		return &internal.Error{Code: internal.ErrNotFound, Message: fmt.Sprintf("cannot find a redirection with key %s", a.Key)}
	}
	// Delete the specified redirection
	delete(r.redirections, a.Key)
	return nil
}

func (r InMemoryRedirectionRepository) FindByKey(key string) (redirection.Redirection, error) {
	// Check if the key exists
	if redirect, exists := r.redirections[key]; exists {
		return redirect, nil
	}
	// Cannot find a redirection with the given key return error
	return redirection.Redirection{}, &internal.Error{Code: internal.ErrNotFound, Message: fmt.Sprintf("cannot find a redirection with key %s", key)}
}
