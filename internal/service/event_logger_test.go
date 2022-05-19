package service_test

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/emanuelefalzone/bitly/internal/service"

	"testing"
)

func TestEventLogger(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test  string
		event event.Event
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:  "Created Event",
			event: event.Created(redirection.Redirection{Key: "short", Location: "http://www.google.com"}),
		}, {
			test:  "Read Event",
			event: event.Read(redirection.Redirection{Key: "short", Location: "http://www.google.com"}),
		}, {
			test:  "Deleted Event",
			event: event.Deleted(redirection.Redirection{Key: "short", Location: "http://www.google.com"}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			ctx := context.Background()
			logger := service.NewEventLogger()
			logger.Consume(ctx, tc.event)
		})
	}
}
