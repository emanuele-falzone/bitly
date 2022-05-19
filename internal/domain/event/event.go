package event

import (
	"time"

	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
)

type Type string

const (
	TypeRead   Type = "read"
	TypeCreate      = "created"
	TypeDelete      = "deleted"
)

type Event struct {
	DateTime    string
	Type        Type
	Redirection redirection.Redirection
}

func New(dt string, type_ Type, a redirection.Redirection) Event {
	return Event{
		DateTime:    dt,
		Type:        type_,
		Redirection: a,
	}
}

func Now(type_ Type, a redirection.Redirection) Event {
	dt := time.Now().UTC().Format(time.RFC3339)
	return New(dt, type_, a)
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
