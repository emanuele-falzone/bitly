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
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dispatcher := event.NewDispatcher(ctx)
			listener := mock.NewMockListener(ctrl)
			dispatcher.Register(listener)
			var wg sync.WaitGroup
			wg.Add(1)
			listener.EXPECT().Consume(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, e event.Event) {
				wg.Done()
			})

			dispatcher.Dispatch(ctx, tc.event)

			wg.Wait()
		})
	}
}
