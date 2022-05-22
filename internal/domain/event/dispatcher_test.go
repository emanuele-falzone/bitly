//go:build unit

package event_test

import (
	"context"
	"sync"
	"testing"

	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/emanuelefalzone/bitly/test/mock"
	"github.com/golang/mock/gomock"
)

func TestDispatch(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test  string
		event event.Event
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:  "Read Event",
			event: event.Read(redirection.Redirection{Key: "short", Location: "http://www.google.com"}),
		}, {
			test:  "Created Event",
			event: event.Created(redirection.Redirection{Key: "short", Location: "http://www.google.com"}),
		}, {
			test:  "Deleted Event",
			event: event.Deleted(redirection.Redirection{Key: "short", Location: "http://www.google.com"}),
		},
	}

	for _, tc := range testCases {
		// Run Tests
		t.Run(tc.test, func(t *testing.T) {
			// Create new context
			ctx := context.Background()

			// Create new mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create new event dispatcher
			dispatcher := event.NewDispatcher(ctx)

			// Create mock listener
			listener := mock.NewMockListener(ctrl)

			// Register mock listener into dispatcher
			dispatcher.Register(listener)

			// Create new wait group
			var wg sync.WaitGroup

			// Add 1 delta
			wg.Add(1)

			// Expect listener consume method to be invoked once
			listener.EXPECT().Consume(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, e event.Event) {
				wg.Done()
			})

			// Dispatch the event
			dispatcher.Dispatch(ctx, tc.event)

			// Wait for wait group counter to be zero
			wg.Wait()
		})
	}
}
