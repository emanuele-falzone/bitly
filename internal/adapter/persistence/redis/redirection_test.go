//go:build integration

package redis_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/redis"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/stretchr/testify/assert"

	_redis "github.com/go-redis/redis/v8"
)

func TestIntegration_Redis_RedirectionRepository_Create(t *testing.T) {
	// Create a redirection that we are going to use in our test cases
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}

	// Build our needed testcase data struct
	type testCase struct {
		test            string
		alreadyExists   bool   // True if the redirection should already exist in the repository
		expectedErr     bool   // True if expecting error
		expectedErrCode string // Expected error code
	}

	// Create new test cases
	testCases := []testCase{
		{
			test:          "New Redirection",
			alreadyExists: false,
			expectedErr:   false,
		}, {
			test:            "Redirection Already Exists",
			alreadyExists:   true,
			expectedErr:     true,
			expectedErrCode: internal.ErrConflict,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			// Create a new context
			ctx := context.Background()

			// Create a new (clean) redis repository
			repository, err := newRedisRepository(ctx)
			assert.Nil(t, err)

			// Create the redirection if it should already exist in the repository
			if tc.alreadyExists {
				repository.Create(ctx, value)
			}

			// Insert the redirection into the repository
			err = repository.Create(ctx, value)

			// Check expected error and expected error code
			if tc.expectedErr {
				assert.Equal(t, tc.expectedErrCode, internal.ErrorCode(err))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestIntegration_Redis_RedirectionRepository_Delete(t *testing.T) {
	// Create a redirection that we are going to use in our test cases
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}

	// Build our needed testcase data struct
	type testCase struct {
		test            string
		alreadyExists   bool   // True if the redirection should already exist in the repository
		expectedErr     bool   // True if expecting error
		expectedErrCode string // Expected error code
	}

	// Create new test cases
	testCases := []testCase{
		{
			test:          "Existing Redirection",
			alreadyExists: true,
			expectedErr:   false,
		}, {
			test:            "Redirection Does Not Exists",
			alreadyExists:   false,
			expectedErr:     true,
			expectedErrCode: internal.ErrNotFound,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			// Create a new context
			ctx := context.Background()

			// Create a new (clean) redis repository
			repository, err := newRedisRepository(ctx)
			assert.Nil(t, err)

			// Create the redirection if it should already exist in the repository
			if tc.alreadyExists {
				repository.Create(ctx, value)
			}

			// Insert the redirection into the repository
			err = repository.Delete(ctx, value)

			// Check expected error and expected error code
			if tc.expectedErr {
				assert.Equal(t, tc.expectedErrCode, internal.ErrorCode(err))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestIntegration_Redis_RedirectionRepository_FindByKey(t *testing.T) {
	// Create a redirection that we are going to use in our test cases
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}

	// Build our needed testcase data struct
	type testCase struct {
		test            string
		alreadyExists   bool   // True if the redirection should already exist in the repository
		expectedErr     bool   // True if expecting error
		expectedErrCode string // Expected error code
	}

	// Create new test cases
	testCases := []testCase{
		{
			test:          "Existing Redirection",
			alreadyExists: true,
			expectedErr:   false,
		}, {
			test:            "Redirection Does Not Exists",
			alreadyExists:   false,
			expectedErr:     true,
			expectedErrCode: internal.ErrNotFound,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			// Create a new context
			ctx := context.Background()

			// Create a new (clean) redis repository
			repository, err := newRedisRepository(ctx)
			assert.Nil(t, err)

			// Create the redirection if it should already exist in the repository
			if tc.alreadyExists {
				repository.Create(ctx, value)
			}

			// Retrieve redirection from repository
			result, err := repository.FindByKey(ctx, value.Key)

			// Check result and expected error
			if tc.expectedErr {
				assert.Equal(t, tc.expectedErrCode, internal.ErrorCode(err))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, value.Key, result.Key)
				assert.Equal(t, value.Location, result.Location)
			}
		})
	}
}

func newRedisRepository(ctx context.Context) (redirection.Repository, error) {
	// Read redis connection string from env
	connectionString, err := internal.GetEnv("INTEGRATION_REDIS_CONNECTION_STRING")
	if err != nil {
		return nil, err
	}

	// Parse connection string and check for errors
	opt, err := _redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	// Create a new redis client
	client := _redis.NewClient(opt)

	// Flush all keys
	err = client.FlushAll(ctx).Err()
	if err != nil {
		return nil, err
	}

	return redis.NewRedirectionRepository(connectionString)
}
