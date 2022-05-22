//go:build unit

package internal_test

import (
	"os"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test            string
		setEnv          bool
		key             string
		value           string
		expectedErr     bool
		expectedErrCode string
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:        "Valid",
			setEnv:      true,
			key:         "somekey",
			value:       "somevalue",
			expectedErr: false,
		}, {
			test:            "Key not set",
			setEnv:          false,
			expectedErr:     true,
			expectedErrCode: internal.ErrNotFound,
		}, {
			test:            "Empty value",
			setEnv:          true,
			key:             "somekey",
			value:           "",
			expectedErr:     true,
			expectedErrCode: internal.ErrInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			// Set env variable if specified
			if tc.setEnv {
				os.Setenv(tc.key, tc.value)
			}

			// Read environment variable
			value, err := internal.GetEnv(tc.key)

			// Check and assert error
			if tc.expectedErr {
				assert.Equal(t, tc.expectedErrCode, internal.ErrorCode(err))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.value, value)
			}
		})
	}
}
