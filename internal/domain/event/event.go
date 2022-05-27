package event

import (
	"time"

	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
)

type Type string

const (
	TypeRead   Type = "read"
	TypeCreate Type = "created"
	TypeDelete Type = "deleted"
)

type Event struct {
	DateTime    string
	Type        Type
	Redirection redirection.Redirection
}

func New(dt string, eventType Type, a redirection.Redirection) Event {
	return Event{
		DateTime:    dt,
		Type:        eventType,
		Redirection: a,
	}
}

func Now(eventType Type, a redirection.Redirection) Event {
	// COmputer current ISO 8601 datetime string
	dt := time.Now().UTC().Format(time.RFC3339)

	return New(dt, eventType, a)
}

func Read(a redirection.Redirection) Event {
	return Now(TypeRead, a)
}

func Created(a redirection.Redirection) Event {
	return Now(TypeCreate, a)
}

func Deleted(a redirection.Redirection) Event {
	return Now(TypeDelete, a)
}
