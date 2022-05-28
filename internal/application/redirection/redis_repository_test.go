//go:build integration

package redirection_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/redirection"

	"github.com/stretchr/testify/assert"

	"github.com/go-redis/redis/v8"
)

func TestIntegration_Redis_RedirectionRepository_Create(t *testing.T) {
	// Build our needed test case data struct
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

			// Create a redirection
			redirectionValue := &redirection.Redirection{Key: "short", Location: "http://www.google.com"}

			// Create the redirection if it should already exist in the repository
			if tc.alreadyExists {
				repository.Create(ctx, redirectionValue)
			}

			// Insert the redirection into the repository
			err = repository.Create(ctx, redirectionValue)

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
	// Build our needed test case data struct
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

			// Create a redirection
			redirectionValue := &redirection.Redirection{Key: "short", Location: "http://www.google.com"}

			// Create the redirection if it should already exist in the repository
			if tc.alreadyExists {
				repository.Create(ctx, redirectionValue)
			}

			// Insert the redirection into the repository
			err = repository.Delete(ctx, redirectionValue)

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
	// Build our needed test case data struct
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

			// Create a redirection
			redirectionValue := &redirection.Redirection{Key: "short", Location: "http://www.google.com"}

			// Create the redirection if it should already exist in the repository
			if tc.alreadyExists {
				repository.Create(ctx, redirectionValue)
			}

			// Retrieve redirection from repository
			result, err := repository.FindByKey(ctx, redirectionValue.Key)

			// Check result and expected error
			if tc.expectedErr {
				assert.Equal(t, tc.expectedErrCode, internal.ErrorCode(err))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, redirectionValue.Key, result.Key)
				assert.Equal(t, redirectionValue.Location, result.Location)
			}
		})
	}
}

func TestIntegration_Redis_RedirectionRepository_FindAll(t *testing.T) {
	// Build our needed test case data struct
	type testCase struct {
		test          string
		redirections  []*redirection.Redirection // Redirection stored inside the repository
		expectedCount int                        // Expected number of results returned
	}

	// Create new test cases
	testCases := []testCase{
		{
			test: "Existing Redirection",
			redirections: []*redirection.Redirection{
				{Key: "a", Location: "http://www.google.com"},
				{Key: "b", Location: "http://www.google.com"},
				{Key: "c", Location: "http://www.google.com"},
				{Key: "d", Location: "http://www.google.com"},
			},
			expectedCount: 4,
		}, {
			test:          "Redirection Does Not Exists",
			redirections:  []*redirection.Redirection{},
			expectedCount: 0,
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

			// Create the redirections in the repository
			for _, value := range tc.redirections {
				repository.Create(ctx, value)
			}

			// Retrieve redirection from repository
			result, err := repository.FindAll(ctx)

			// Assert error
			assert.Nil(t, err)

			// Assert result length
			assert.Equal(t, tc.expectedCount, len(result))
		})
	}
}

func newRedisRepository(ctx context.Context) (*redirection.RedisRepository, error) {
	// Read redis connection string from env
	connectionString, err := internal.GetEnv("INTEGRATION_REDIS_CONNECTION_STRING")
	if err != nil {
		return nil, err
	}

	// Parse connection string and check for errors
	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	// Create a new redis client
	client := redis.NewClient(opt)

	// Flush all keys
	err = client.FlushAll(ctx).Err()
	if err != nil {
		return nil, err
	}

	return redirection.NewRedisRepository(connectionString)
}
