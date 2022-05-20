package integration_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/memory"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/mongo"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/emanuelefalzone/bitly/test/util"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryEventRepository(t *testing.T) {
	RunTestEventRepository(t, func() (event.Repository, error) {
		return memory.NewEventRepository(), nil
	})
}

func TestRedisEventRepository(t *testing.T) {
	RunTestEventRepository(t, func() (event.Repository, error) {
		// Create new context
		ctx := context.Background()

		// Read redis connection string from env
		connectionString, err := internal.GetEnv("INTEGRATION_MONGO_CONNECTION_STRING")
		if err != nil {
			panic(err)
		}

		// Parse connection string and check for errors
		err = util.ClearMongo(ctx, connectionString, mongo.DB)
		if err != nil {
			return nil, err
		}

		return mongo.NewEventRepository(connectionString)
	})
}

func RunTestEventRepository(t *testing.T, newRepository func() (event.Repository, error)) {
	// Build our needed testcase data struct
	type testCase struct {
		test string
		fn   func(*testing.T, func() (event.Repository, error))
	}
	// Create new test cases
	testCases := []testCase{
		{
			test: "TestCreate",
			fn:   _TestEventCreate,
		}, {
			test: "TestFindByRedirection",
			fn:   _TestEventFindByRedirection,
		},
	}

	for _, tc := range testCases {
		// Run Tests
		t.Run(tc.test, func(t *testing.T) {
			tc.fn(t, newRepository)
		})
	}

}

func _TestEventCreate(t *testing.T, newRepository func() (event.Repository, error)) {
	// Create a redirection that we are going to use in our test cases
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}

	// Build our needed testcase data struct
	type testCase struct {
		test  string
		event event.Event
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:  "Created Event",
			event: event.Created(value),
		}, {
			test:  "Read Event",
			event: event.Read(value),
		}, {
			test:  "Deleted Event",
			event: event.Deleted(value),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			ctx := context.Background()
			repository, _ := newRepository()

			err := repository.Create(ctx, tc.event)

			assert.Nil(t, err)
		})
	}
}

func _TestEventFindByRedirection(t *testing.T, newRepository func() (event.Repository, error)) {
	// Create a redirection that we are going to use in our test cases
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}

	// Build our needed testcase data struct
	type testCase struct {
		test        string
		count       int
		expectedErr bool
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:        "Zero",
			count:       0,
			expectedErr: true,
		}, {
			test:        "Ten",
			count:       10,
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			ctx := context.Background()
			repository, _ := newRepository()

			for i := 0; i < tc.count; i++ {
				err := repository.Create(ctx, event.Read(value))
				assert.Nil(t, err)
			}

			events, err := repository.FindByRedirection(ctx, value)
			assert.Len(t, events, tc.count)
			if tc.expectedErr {
				assert.Equal(t, internal.ErrNotFound, internal.ErrorCode(err))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
