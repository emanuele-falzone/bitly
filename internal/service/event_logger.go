package service

import (
	"context"
	"log"

	"github.com/emanuelefalzone/bitly/internal/domain/event"
)

type EventLogger struct{}

func NewEventLogger() event.Listener {
	return &EventLogger{}
}

func (l *EventLogger) Consume(ctx context.Context, e event.Event) {
	log.Printf("Key: %s, Location: %s, Event: %s, DateTime: %s\n", e.Redirection.Key, e.Redirection.Location, e.Type, e.DateTime)
}
