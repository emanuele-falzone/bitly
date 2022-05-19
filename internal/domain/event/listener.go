package service

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal/domain/event"
)

type Listener interface {
	Consume(context.Context, event.Event)
}
