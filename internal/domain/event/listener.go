package event

import (
	"context"
)

type Listener interface {
	Consume(context.Context, Event)
}
