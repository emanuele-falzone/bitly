package redirection

import "context"

// Repository allows to store, delete and find redirections
type Repository interface {
	// Create stores the given redirection inside the repository
	Create(context.Context, *Redirection) error
	// Delete deletes the given redirection from the repository
	Delete(context.Context, *Redirection) error
	// FindByKey retrieves, if exists, a redirection by key
	FindByKey(context.Context, string) (*Redirection, error)
	// FindAll retrieves all the redirections stored inside the repository
	FindAll(context.Context) ([]*Redirection, error)
}
