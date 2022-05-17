package integration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/memory"
	_redis "github.com/emanuelefalzone/bitly/internal/adapter/persistence/redis"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRedirectionRepository(t *testing.T) {
	RunTestRedirectionRepository(t, func() (redirection.Repository, error) {
		return memory.NewRedirectionRepository(), nil
	})
}

func TestRedisRedirectionRepository(t *testing.T) {
	RunTestRedirectionRepository(t, func() (redirection.Repository, error) {
		// Create new context
		ctx := context.Background()

		// Read redis connection string from env
		connectionString, err := internal.GetEnv("INTEGRATION_REDIS_CONNECTION_STRING")
		if err != nil {
			panic(err)
		}

		// Parse connection string and check for errors
		opt, err := redis.ParseURL(connectionString)
		if err != nil {
			return nil, err
		}

		// Create a new redis client
		client := redis.NewClient(opt)

		// Flush all keys
		flushErr := client.FlushAll(ctx).Err()
		if flushErr != nil {
			return nil, flushErr
		}

		return _redis.NewRedirectionRepository(connectionString)
	})
}

func RunTestRedirectionRepository(t *testing.T, newRepository func() (redirection.Repository, error)) {
	// Build our needed testcase data struct
	type testCase struct {
		test string
		fn   func(*testing.T, func() (redirection.Repository, error))
	}
	// Create new test cases
	testCases := []testCase{
		{
			test: "TestCreate",
			fn:   _TestRedirectionCreate,
		}, {
			test: "TestDelete",
			fn:   _TestRedirectionDelete,
		}, {
			test: "TestFindByKey",
			fn:   _TestRedirectionFindByKey,
		},
	}

	for _, tc := range testCases {
		// Run Tests
		t.Run(tc.test, func(t *testing.T) {
			tc.fn(t, newRepository)
		})
	}

}

func _TestRedirectionCreate(t *testing.T, newRepository func() (redirection.Repository, error)) {
	// Create a redirection that we are going to use in our test cases
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}

	// Build our needed testcase data struct
	type testCase struct {
		test            string
		alreadyExists   bool
		expectedErr     bool
		expectedErrCode string
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:          "New Redirection",
			alreadyExists: false,
			expectedErr:   false,
		}, {
			test:            "Redirection already exists",
			alreadyExists:   true,
			expectedErr:     true,
			expectedErrCode: internal.ErrConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			ctx := context.Background()
			repository, _ := newRepository()

			if tc.alreadyExists {
				repository.Create(ctx, value)
			}

			err := repository.Create(ctx, value)

			if tc.expectedErr {
				assert.Equal(t, tc.expectedErrCode, internal.ErrorCode(err))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func _TestRedirectionDelete(t *testing.T, newRepository func() (redirection.Repository, error)) {
	// Create a redirection that we are going to use in our test cases
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}

	// Build our needed testcase data struct
	type testCase struct {
		test            string
		alreadyExists   bool
		expectedErr     bool
		expectedErrCode string
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:          "Existing Redirection",
			alreadyExists: true,
			expectedErr:   false,
		}, {
			test:            "Redirection does not exists",
			alreadyExists:   false,
			expectedErr:     true,
			expectedErrCode: internal.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			ctx := context.Background()
			repository, _ := newRepository()

			if tc.alreadyExists {
				repository.Create(ctx, value)
			}

			err := repository.Delete(ctx, value)

			if tc.expectedErr {
				assert.Equal(t, tc.expectedErrCode, internal.ErrorCode(err))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func _TestRedirectionFindByKey(t *testing.T, newRepository func() (redirection.Repository, error)) {
	// Create a redirection that we are going to use in our test cases
	value := redirection.Redirection{Key: "short", Location: "http://www.google.com"}

	// Build our needed testcase data struct
	type testCase struct {
		test            string
		alreadyExists   bool
		expectedErr     bool
		expectedErrCode string
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:          "Existing Redirection",
			alreadyExists: true,
			expectedErr:   false,
		}, {
			test:            "Redirection does not exists",
			alreadyExists:   false,
			expectedErr:     true,
			expectedErrCode: internal.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			ctx := context.Background()
			repository, _ := newRepository()

			if tc.alreadyExists {
				repository.Create(ctx, value)
			}

			result, err := repository.FindByKey(ctx, value.Key)

			fmt.Println(result)
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
