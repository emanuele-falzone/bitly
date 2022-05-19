package event

import (
	"context"
)

type Dispatcher struct {
	events    chan Event
	listeners []Listener
}

func NewDispatcher(ctx context.Context) *Dispatcher {
	dispatcher := &Dispatcher{
		events: make(chan Event),
	}

	// Start goroutine to let listeners consume events
	go dispatcher.Consume(ctx)
	return dispatcher
}

// Register a new listener
func (d *Dispatcher) Register(listener Listener) {
	d.listeners = append(d.listeners, listener)
}

// Send event to events channel
func (d *Dispatcher) Dispatch(ctx context.Context, e Event) {
	d.events <- e
}

// Continuously wait for an event on the channel and make listeners consume it
func (d *Dispatcher) Consume(ctx context.Context) {
	for {
		event := <-d.events
		for _, listener := range d.listeners {
			listener.Consume(ctx, event)
		}
	}
}
