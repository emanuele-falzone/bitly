//-go:build unit

package event_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/event"
	"github.com/emanuelefalzone/bitly/internal/application/redirection"

	"github.com/stretchr/testify/assert"
)

func TestMemory_EventRepository_Create(t *testing.T) {
	// Build our needed test case data struct
	type testCase struct {
		test      string
		eventType event.Type // Event to be inserted into the repository
	}

	// Create new test cases
	testCases := []testCase{
		{
			test:      "Created Event",
			eventType: event.TypeCreate,
		}, {
			test:      "Read Event",
			eventType: event.TypeRead,
		}, {
			test:      "Deleted Event",
			eventType: event.TypeDelete,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			// Create a new context
			ctx := context.Background()

			// Create a new (clean) memory repository
			repository := event.NewInMemoryRepository()

			// Create a redirection that we are going to use in our test cases
			redirectionValue := &redirection.Redirection{Key: "short", Location: "http://www.google.com"}

			// Create an event
			eventValue := event.Now(tc.eventType, redirectionValue)

			// Insert the event into the repository
			err := repository.Create(ctx, eventValue)

			// Check if the event has been inserted without error
			assert.Nil(t, err)
		})
	}
}

func TestMemory_EventRepository_FindByRedirection(t *testing.T) {
	// Build our needed test case data struct
	type testCase struct {
		test            string
		count           int    // Number of times the event has to be inserted into the repository
		expectedErr     bool   // True if expecting error
		expectedErrCode string // Expected error code
	}

	// Create new test cases
	testCases := []testCase{
		{
			test:            "Zero",
			count:           0,
			expectedErr:     true,
			expectedErrCode: internal.ErrNotFound,
		}, {
			test:        "Ten",
			count:       10,
			expectedErr: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			// Create a new context
			ctx := context.Background()

			// Create a new (clean) memory repository
			repository := event.NewInMemoryRepository()

			// Create a redirection that we are going to use in our test cases
			redirectionValue := &redirection.Redirection{Key: "short", Location: "http://www.google.com"}

			// Insert the event into the repository tc.count times
			for i := 0; i < tc.count; i++ {
				// Insert the event into the repository
				err := repository.Create(ctx, event.Now(event.TypeRead, redirectionValue))
				assert.Nil(t, err)
			}

			// Find the events in the repository
			events, err := repository.FindByRedirection(ctx, redirectionValue)

			// Check events length
			assert.Len(t, events, tc.count)

			// Check expected error and expected error code
			if tc.expectedErr {
				assert.Equal(t, tc.expectedErrCode, internal.ErrorCode(err))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
