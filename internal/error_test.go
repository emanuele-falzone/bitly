package internal_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/emanuelefalzone/bitly/internal"
)

func TestError_ErrorCodeNil(t *testing.T) {
	assert.Equal(t, "", internal.ErrorCode(nil))
}

func TestError_ErrorCode(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test         string
		code         string
		expectedCode string
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:         "Internal",
			code:         internal.ErrInternal,
			expectedCode: internal.ErrInternal,
		}, {
			test:         "Conflict",
			code:         internal.ErrConflict,
			expectedCode: internal.ErrConflict,
		}, {
			test:         "NotFound",
			code:         internal.ErrNotFound,
			expectedCode: internal.ErrNotFound,
		}, {
			test:         "Invalid",
			code:         internal.ErrInvalid,
			expectedCode: internal.ErrInvalid,
		}, {
			test:         "Internal default",
			code:         "",
			expectedCode: internal.ErrInternal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := &internal.Error{Code: tc.code}
			assert.Equal(t, tc.expectedCode, internal.ErrorCode(err))
		})
	}
}

func TestError_ErrorCodeWrap(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test         string
		wrapCode     string
		expectedCode string
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:         "Internal",
			wrapCode:     internal.ErrInternal,
			expectedCode: internal.ErrInternal,
		}, {
			test:         "Conflict",
			wrapCode:     internal.ErrConflict,
			expectedCode: internal.ErrConflict,
		}, {
			test:         "NotFound",
			wrapCode:     internal.ErrNotFound,
			expectedCode: internal.ErrNotFound,
		}, {
			test:         "Invalid",
			wrapCode:     internal.ErrInvalid,
			expectedCode: internal.ErrInvalid,
		}, {
			test:         "Internal default",
			wrapCode:     "",
			expectedCode: internal.ErrInternal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := &internal.Error{Err: &internal.Error{Code: tc.wrapCode}}
			assert.Equal(t, tc.expectedCode, internal.ErrorCode(err))
		})
	}
}

func TestError_ErrorCodeWrapWrap(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test         string
		wrapWrapCode string
		expectedCode string
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:         "Internal",
			wrapWrapCode: internal.ErrInternal,
			expectedCode: internal.ErrInternal,
		}, {
			test:         "Conflict",
			wrapWrapCode: internal.ErrConflict,
			expectedCode: internal.ErrConflict,
		}, {
			test:         "NotFound",
			wrapWrapCode: internal.ErrNotFound,
			expectedCode: internal.ErrNotFound,
		}, {
			test:         "Invalid",
			wrapWrapCode: internal.ErrInvalid,
			expectedCode: internal.ErrInvalid,
		}, {
			test:         "Internal default",
			wrapWrapCode: "",
			expectedCode: internal.ErrInternal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := &internal.Error{Err: &internal.Error{Err: &internal.Error{Code: tc.wrapWrapCode}}}
			assert.Equal(t, tc.expectedCode, internal.ErrorCode(err))
		})
	}
}

func TestError_ErrorMessageNil(t *testing.T) {
	assert.Equal(t, "", internal.ErrorMessage(nil))
}

func TestError_ErrorMessageEmpty(t *testing.T) {
	err := &internal.Error{Err: errors.New("something bad")}
	assert.Equal(t, "An internal error has occurred. Please contact technical support.", internal.ErrorMessage(err))
}

func TestError_ErrorMessage(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test            string
		message         string
		expectedMessage string
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:            "Internal",
			message:         "internal error",
			expectedMessage: "internal error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := &internal.Error{Message: tc.message}
			assert.Equal(t, tc.expectedMessage, internal.ErrorMessage(err))
		})
	}
}

func TestError_ErrorMessageWrap(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test            string
		wrapMessage     string
		expectedMessage string
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:            "Internal",
			wrapMessage:     "internal error",
			expectedMessage: "internal error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := &internal.Error{Err: &internal.Error{Message: tc.wrapMessage}}
			assert.Equal(t, tc.expectedMessage, internal.ErrorMessage(err))
		})
	}
}

func TestError_ErrorError(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test          string
		err           error
		expectedError string
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:          "Internal",
			err:           errors.New("something bad"),
			expectedError: "something bad",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := &internal.Error{Err: tc.err}
			assert.Equal(t, tc.expectedError, err.Error())
		})
	}
}

func TestError_ErrorErrorWrap(t *testing.T) {

	first := errors.New("something bad happened")
	err := internal.Error{Op: "test", Err: first}

	assert.Equal(t, "test: something bad happened", err.Error())
}

func TestError_ErrorErrorWrapWrap(t *testing.T) {

	first := errors.New("something bad happened")
	second := &internal.Error{Op: "write", Err: first}
	err := &internal.Error{Op: "test", Err: second}

	assert.Equal(t, "test: write: something bad happened", err.Error())
}

func TestError_ErrorErrorWrapWrapWithCode(t *testing.T) {

	first := &internal.Error{Code: internal.ErrInternal, Op: "something bad happened"}
	second := &internal.Error{Op: "write", Err: first}
	err := internal.Error{Op: "test", Err: second}

	assert.Equal(t, "test: write: something bad happened: <internal> ", err.Error())
}
