package service

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal/domain/event"
)

type EventStore struct {
	repository event.Repository
}

func NewEventStore(repository event.Repository) event.Listener {
	return &EventStore{repository: repository}
}

func (s *EventStore) Consume(ctx context.Context, e event.Event) {
	s.repository.Create(ctx, e)
}
