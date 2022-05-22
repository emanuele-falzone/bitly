//go:build unit

package memory_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/memory"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/stretchr/testify/assert"
)

func TestMemory_RedirectionRepository_Create(t *testing.T) {
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

			// Create a new (clean) memory repository
			repository := memory.NewRedirectionRepository()

			// Create the redirection if it should already exist in the repository
			if tc.alreadyExists {
				repository.Create(ctx, value)
			}

			// Insert the redirection into the repository
			err := repository.Create(ctx, value)

			// Check expected error and expected error code
			if tc.expectedErr {
				assert.Equal(t, tc.expectedErrCode, internal.ErrorCode(err))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestMemory_RedirectionRepository_Delete(t *testing.T) {
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

			// Create a new (clean) memory repository
			repository := memory.NewRedirectionRepository()

			// Create the redirection if it should already exist in the repository
			if tc.alreadyExists {
				repository.Create(ctx, value)
			}

			// Insert the redirection into the repository
			err := repository.Delete(ctx, value)

			// Check expected error and expected error code
			if tc.expectedErr {
				assert.Equal(t, tc.expectedErrCode, internal.ErrorCode(err))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestMemory_RedirectionRepository_FindByKey(t *testing.T) {
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

			// Create a new (clean) memory repository
			repository := memory.NewRedirectionRepository()

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

func TestMemory_RedirectionRepository_FindAll(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test          string
		redirections  []redirection.Redirection // Redirection stored inside the repository
		alreadyExists bool                      // True if the redirection should already exist in the repository
		expectedCount int                       // Expected number of results returned
	}

	// Create new test cases
	testCases := []testCase{
		{
			test: "Existing Redirection",
			redirections: []redirection.Redirection{
				{Key: "a", Location: "http://www.google.com"},
				{Key: "b", Location: "http://www.google.com"},
				{Key: "c", Location: "http://www.google.com"},
				{Key: "d", Location: "http://www.google.com"},
			},
			expectedCount: 4,
		}, {
			test:          "Redirection Does Not Exists",
			redirections:  []redirection.Redirection{},
			expectedCount: 0,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			// Create a new context
			ctx := context.Background()

			// Create a new (clean) memory repository
			repository := memory.NewRedirectionRepository()

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
