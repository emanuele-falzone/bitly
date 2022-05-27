package redirection

import "context"

type Repository interface {
	Create(context.Context, Redirection) error
	Delete(context.Context, Redirection) error
	FindByKey(context.Context, string) (Redirection, error)
	FindAll(context.Context) ([]Redirection, error)
}
