package internal_test

import (
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	// Define struct to validate
	type Sample struct {
		Location string `validate:"required,url"`
		Count    int    `validate:"min=3"`
	}
	// Build our needed testcase data struct
	type testCase struct {
		test         string
		location     string
		count        int
		expectedErr  bool
		expectedCode string
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:        "Internal",
			location:    "http://www.google.com",
			count:       10,
			expectedErr: false,
		}, {
			test:         "Internal",
			location:     "",
			count:        10,
			expectedErr:  true,
			expectedCode: internal.ErrInvalid,
		}, {
			test:         "Internal",
			location:     "google",
			count:        10,
			expectedErr:  true,
			expectedCode: internal.ErrInvalid,
		}, {
			test:         "Internal",
			location:     "http://www.google.com",
			count:        1,
			expectedErr:  true,
			expectedCode: internal.ErrInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {

			sample := &Sample{Location: tc.location, Count: tc.count}
			err := internal.Validate(sample)
			if tc.expectedErr {
				assert.Equal(t, tc.expectedCode, internal.ErrorCode(err))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
