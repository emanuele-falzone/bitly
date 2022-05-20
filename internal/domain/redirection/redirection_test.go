//go:build unit

package redirection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
)

func TestRedirection_New(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test            string
		location        string
		key             string
		expectedErr     bool
		expectedErrCode string
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:            "Malformed Location",
			location:        "google.com",
			key:             "short",
			expectedErr:     true,
			expectedErrCode: internal.ErrInvalid,
		}, {
			test:            "Empty Key",
			location:        "http://www.google.com",
			key:             "",
			expectedErr:     true,
			expectedErrCode: internal.ErrInvalid,
		}, {
			test:        "Valid location",
			location:    "http://www.google.com",
			key:         "short",
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		// Run Tests
		t.Run(tc.test, func(t *testing.T) {
			// Create a new redirection
			value, err := redirection.New(tc.key, tc.location)

			if tc.expectedErr {
				assert.Equal(t, tc.expectedErrCode, internal.ErrorCode(err))
			} else {
				assert.Equal(t, tc.key, value.Key)
				assert.Equal(t, tc.location, value.Location)
			}
		})
	}
}
