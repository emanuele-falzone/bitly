package grpc_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/memory"
	"github.com/emanuelefalzone/bitly/internal/adapter/service/grpc"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/emanuelefalzone/bitly/internal/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRedirectionCreate(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test            string
		location        string
		expectedErr     bool
		expectedErrCode codes.Code
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:        "Valid location",
			location:    "http://www.google.com",
			expectedErr: false,
		}, {
			test:            "Invalid location",
			location:        "google",
			expectedErr:     true,
			expectedErrCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			ctx := context.Background()
			repository := memory.NewRedirectionRepository()
			generator := service.NewRandomKeyGenerator(0)
			application := application.New(repository, generator)
			server := grpc.NewServer(application)

			_, err := server.CreateRedirection(ctx, &grpc.CreateRedirectionRequest{Location: tc.location})

			if tc.expectedErr {
				err, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedErrCode, err.Code())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestRedirectionDelete(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test            string
		location        string
		key             string
		alreadyExists   bool
		expectedErr     bool
		expectedErrCode codes.Code
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:          "Valid location",
			location:      "http://www.google.com",
			key:           "short",
			alreadyExists: true,
			expectedErr:   false,
		}, {
			test:            "Invalid location",
			location:        "short",
			alreadyExists:   false,
			expectedErr:     true,
			expectedErrCode: codes.NotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			ctx := context.Background()
			repository := memory.NewRedirectionRepository()
			generator := service.NewRandomKeyGenerator(0)
			application := application.New(repository, generator)
			server := grpc.NewServer(application)

			if tc.alreadyExists {
				repository.Create(ctx, redirection.Redirection{Key: tc.key, Location: tc.location})
			}

			_, err := server.DeleteRedirection(ctx, &grpc.DeleteRedirectionRequest{Key: tc.key})

			if tc.expectedErr {
				err, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedErrCode, err.Code())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetRedirectionLocation(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test            string
		location        string
		key             string
		alreadyExists   bool
		expectedErr     bool
		expectedErrCode codes.Code
	}
	// Create new test cases
	testCases := []testCase{
		{
			test:          "Valid location",
			location:      "http://www.google.com",
			key:           "short",
			alreadyExists: true,
			expectedErr:   false,
		}, {
			test:            "Invalid location",
			location:        "short",
			alreadyExists:   false,
			expectedErr:     true,
			expectedErrCode: codes.NotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			ctx := context.Background()
			repository := memory.NewRedirectionRepository()
			generator := service.NewRandomKeyGenerator(0)
			application := application.New(repository, generator)
			server := grpc.NewServer(application)

			if tc.alreadyExists {
				repository.Create(ctx, redirection.Redirection{Key: tc.key, Location: tc.location})
			}

			_, err := server.GetRedirectionLocation(ctx, &grpc.GetRedirectionLocationRequest{Key: tc.key})

			if tc.expectedErr {
				err, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedErrCode, err.Code())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
