//go:build integration

package mongo_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/mongo"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/stretchr/testify/assert"

	_mongo "go.mongodb.org/mongo-driver/mongo"
	_mongo_options "go.mongodb.org/mongo-driver/mongo/options"
)

func TestIntegration_Mongo_EventRepository_Create(t *testing.T) {
	// Create a redirection that we are going to use in our test cases
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}

	// Build our needed testcase data struct
	type testCase struct {
		test  string
		event event.Event // Event to be inserted into the repository
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

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			// Create a new context
			ctx := context.Background()

			// Create a new (clean) mongo repository
			repository, err := newMongoRepository(ctx)
			assert.Nil(t, err)

			// Insert the event into the repository
			err = repository.Create(ctx, tc.event)

			// Check if the event has been inserted without error
			assert.Nil(t, err)
		})
	}
}

func TestIntegration_Mongo_EventRepository_FindByRedirection(t *testing.T) {
	// Create a redirection that we are going to use in our test cases
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}

	// Build our needed testcase data struct
	type testCase struct {
		test            string
		count           int    // Number of times the event has to be inserted into the repository
		expectedErr     bool   // Expecting error
		expectedErrCode string // Expected error code (if expectedErr is true)
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

			// Create a new (clean) mongo repository
			repository, err := newMongoRepository(ctx)
			assert.Nil(t, err)

			// Insert the event into the repository tc.count times
			for i := 0; i < tc.count; i++ {
				// Insert the event into the repository
				err := repository.Create(ctx, event.Read(value))
				assert.Nil(t, err)
			}

			// Find the events in the repository
			events, err := repository.FindByRedirection(ctx, value)

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

func newMongoRepository(ctx context.Context) (event.Repository, error) {
	// Read Mongo connection string from env
	connectionString, err := internal.GetEnv("INTEGRATION_MONGO_CONNECTION_STRING")
	if err != nil {
		return nil, err
	}

	// Create new mongo client with the given connection string
	client, err := _mongo.NewClient(_mongo_options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	// Connect with the mongo instance
	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	// Select database
	db := client.Database(mongo.DB)

	// Drop db
	err = db.Drop(context.Background())
	if err != nil {
		return nil, err
	}

	return mongo.NewEventRepository(connectionString)
}
