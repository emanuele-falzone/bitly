//go:build unit

package event_test

import (
	"testing"

	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/stretchr/testify/assert"
)

func TestDomainEvent_Type(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test              string
		event             event.Event
		expectedEventType event.Type
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:              "Read",
			event:             event.Read(redirection.Redirection{}),
			expectedEventType: event.TypeRead,
		}, {
			test:              "Create",
			event:             event.Created(redirection.Redirection{}),
			expectedEventType: event.TypeCreate,
		}, {
			test:              "Delete",
			event:             event.Deleted(redirection.Redirection{}),
			expectedEventType: event.TypeDelete,
		},
	}

	for _, tc := range testCases {
		// Run Tests
		t.Run(tc.test, func(t *testing.T) {
			assert.Equal(t, tc.expectedEventType, tc.event.Type)
		})
	}
}

func TestDomainEvent_DateTime(t *testing.T) {
	// Build our needed testcase data struct
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
			event := event.New(tc.dateTime, event.TypeRead, redirection.Redirection{})
			assert.Equal(t, tc.dateTime, event.DateTime)
		})
	}
}
