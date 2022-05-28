package event

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal/application/redirection"
)

type Repository interface {
	Create(context.Context, Event) error
	FindByRedirection(context.Context, redirection.Redirection) ([]Event, error)
}
