//go:build unit

package event_test

import (
	"testing"

	"github.com/emanuelefalzone/bitly/internal/application/event"
	"github.com/emanuelefalzone/bitly/internal/application/redirection"
	"github.com/stretchr/testify/assert"
)

func TestDomainEvent_Type(t *testing.T) {
	// Build our needed test case data struct
	type testCase struct {
		test              string
		eventType         event.Type
		expectedEventType event.Type
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:              "Read",
			eventType:         event.TypeRead,
			expectedEventType: event.TypeRead,
		}, {
			test:              "Create",
			eventType:         event.TypeCreate,
			expectedEventType: event.TypeCreate,
		}, {
			test:              "Delete",
			eventType:         event.TypeDelete,
			expectedEventType: event.TypeDelete,
		},
	}

	for _, tc := range testCases {
		// Run Tests
		t.Run(tc.test, func(t *testing.T) {
			value := &redirection.Redirection{Key: "short", Location: "http://www.google.com"}
			event := event.Now(tc.eventType, value)
			assert.Equal(t, tc.expectedEventType, event.Type)
		})
	}
}

func TestDomainEvent_DateTime(t *testing.T) {
	// Build our needed test case data struct
	type testCase struct {
		test     string
		dateTime string
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:     "Read",
			dateTime: "2022-01-01T12:00:00Z",
		}, {
			test:     "Create",
			dateTime: "2022-01-02T12:00:00Z",
		}, {
			test:     "Delete",
			dateTime: "2022-01-03T12:00:00Z",
		},
	}

	for _, tc := range testCases {
		// Run Tests
		t.Run(tc.test, func(t *testing.T) {
			value := &redirection.Redirection{Key: "short", Location: "http://www.google.com"}
			event := event.New(tc.dateTime, event.TypeRead, value)
			assert.Equal(t, tc.dateTime, event.DateTime)
		})
	}
}
