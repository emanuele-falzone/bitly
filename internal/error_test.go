//--go:build unit

package internal_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/emanuelefalzone/bitly/internal"
)

func TestError_NilError(t *testing.T) {
	assert.Equal(t, "", internal.ErrorCode(nil))
	assert.Equal(t, "", internal.ErrorMessage(nil))
}

func TestError_EmptyError(t *testing.T) {
	err := &internal.Error{Err: errors.New("something bad")}
	assert.Equal(t, "An internal error has occurred. Please contact technical support.", internal.ErrorMessage(err))
	assert.Equal(t, internal.ErrInternal, internal.ErrorCode(err))
}

func TestError(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test            string
		op              string
		wraps           int // Number of times the error has been wrapped
		code            string
		message         string
		err             error // Starting err
		expectedCode    string
		expectedMessage string
		expectedError   string
	}
	// Create new test cases
	testCases := []testCase{
		// Testing different error code
		{
			test:            "Internal",
			code:            internal.ErrInternal,
			message:         "something bad",
			expectedCode:    internal.ErrInternal,
			expectedMessage: "something bad",
			expectedError:   fmt.Sprintf("<%s> something bad", internal.ErrInternal),
		}, {
			test:            "Conflict",
			code:            internal.ErrConflict,
			message:         "something bad",
			expectedCode:    internal.ErrConflict,
			expectedMessage: "something bad",
			expectedError:   fmt.Sprintf("<%s> something bad", internal.ErrConflict),
		}, {
			test:            "NotFound",
			code:            internal.ErrNotFound,
			message:         "something bad",
			expectedCode:    internal.ErrNotFound,
			expectedMessage: "something bad",
			expectedError:   fmt.Sprintf("<%s> something bad", internal.ErrNotFound),
		}, {
			test:            "Invalid",
			code:            internal.ErrInvalid,
			message:         "something bad",
			expectedCode:    internal.ErrInvalid,
			expectedMessage: "something bad",
			expectedError:   fmt.Sprintf("<%s> something bad", internal.ErrInvalid),
		}, {
			test:            "Default",
			code:            "",
			message:         "something bad",
			expectedCode:    internal.ErrInternal,
			expectedMessage: "something bad",
			expectedError:   "something bad",
		},

		// Testing different error code with error wrapping
		{
			test:            "Internal Wrap 10",
			wraps:           10,
			code:            internal.ErrInternal,
			message:         "something bad",
			expectedCode:    internal.ErrInternal,
			expectedMessage: "something bad",
			expectedError:   fmt.Sprintf("<%s> something bad", internal.ErrInternal),
		}, {
			test:            "Conflict Wrap 10",
			wraps:           10,
			code:            internal.ErrConflict,
			message:         "something bad",
			expectedCode:    internal.ErrConflict,
			expectedMessage: "something bad",
			expectedError:   fmt.Sprintf("<%s> something bad", internal.ErrConflict),
		}, {
			test:            "NotFound Wrap 10",
			wraps:           10,
			code:            internal.ErrNotFound,
			message:         "something bad",
			expectedCode:    internal.ErrNotFound,
			expectedMessage: "something bad",
			expectedError:   fmt.Sprintf("<%s> something bad", internal.ErrNotFound),
		}, {
			test:            "Invalid Wrap 10",
			wraps:           10,
			code:            internal.ErrInvalid,
			message:         "something bad",
			expectedCode:    internal.ErrInvalid,
			expectedMessage: "something bad",
			expectedError:   fmt.Sprintf("<%s> something bad", internal.ErrInvalid),
		}, {
			test:            "Default Wrap 10",
			wraps:           10,
			code:            "",
			message:         "something bad",
			expectedCode:    internal.ErrInternal,
			expectedMessage: "something bad",
			expectedError:   "something bad",
		},

		// Testing with previously existing Error
		{
			test:            "Internal Existing Error",
			code:            internal.ErrInternal,
			op:              "Create item in db",
			message:         "something bad",
			err:             &internal.Error{Message: "cannot connect to db"},
			expectedCode:    internal.ErrInternal,
			expectedMessage: "something bad",
			expectedError:   "Create item in db: cannot connect to db",
		}, {
			test:            "Conflict Existing Error",
			code:            internal.ErrConflict,
			op:              "Create item in db",
			message:         "something bad",
			err:             &internal.Error{Message: "cannot connect to db"},
			expectedCode:    internal.ErrConflict,
			expectedMessage: "something bad",
			expectedError:   "Create item in db: cannot connect to db",
		}, {
			test:            "NotFound Existing Error",
			code:            internal.ErrNotFound,
			op:              "Create item in db",
			message:         "something bad",
			err:             &internal.Error{Message: "cannot connect to db"},
			expectedCode:    internal.ErrNotFound,
			expectedMessage: "something bad",
			expectedError:   "Create item in db: cannot connect to db",
		}, {
			test:            "Invalid Existing Error",
			code:            internal.ErrInvalid,
			op:              "Create item in db",
			message:         "something bad",
			err:             &internal.Error{Message: "cannot connect to db"},
			expectedCode:    internal.ErrInvalid,
			expectedMessage: "something bad",
			expectedError:   "Create item in db: cannot connect to db",
		}, {
			test:            "Default Existing Error",
			code:            "",
			op:              "Create item in db",
			message:         "something bad",
			err:             &internal.Error{Message: "cannot connect to db"},
			expectedCode:    internal.ErrInternal,
			expectedMessage: "something bad",
			expectedError:   "Create item in db: cannot connect to db",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			// Create starting error
			err := &internal.Error{Code: tc.code, Message: tc.message, Err: tc.err, Op: tc.op}

			// Wrap error wraps times
			for i := 0; i < tc.wraps; i++ {
				err = &internal.Error{Err: err}
			}

			// Assert expected error code
			assert.Equal(t, tc.expectedCode, internal.ErrorCode(err))

			// Assert expected error message
			assert.Equal(t, tc.expectedMessage, internal.ErrorMessage(err))

			// Assert expected error string
			assert.Equal(t, tc.expectedError, err.Error())
		})
	}
}
