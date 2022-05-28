package event

import (
	"time"

	"github.com/emanuelefalzone/bitly/internal/application/redirection"
)

// Type represents the event type
type Type string

const (
	// TypeRead signals that a redirection has been read
	TypeRead Type = "read"
	// TypeRead signals that a redirection has been created
	TypeCreate Type = "created"
	// TypeRead signals that a redirection has been deleted
	TypeDelete Type = "deleted"
)

// Event associates a redirection to a datetime and event type
type Event struct {
	DateTime    string
	Type        Type
	Redirection *redirection.Redirection
}

// New creates a new event given the datetime, event type and redirection
func New(datetime string, eventType Type, value *redirection.Redirection) *Event {
	return &Event{
		DateTime:    datetime,
		Type:        eventType,
		Redirection: value,
	}
}

// Now creates a new event with the current datetime and specified event type and redirection
func Now(eventType Type, value *redirection.Redirection) *Event {
	// COmputer current ISO 8601 datetime string
	datetime := time.Now().UTC().Format(time.RFC3339)

	return New(datetime, eventType, value)
}
