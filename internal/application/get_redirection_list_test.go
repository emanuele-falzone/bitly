//go:build unit

package application_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/internal/application/redirection"
	"github.com/emanuelefalzone/bitly/test/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestApplicationQuery_RedirectionList(t *testing.T) {
	// Build our needed testcase data struct
	type testCaseRedirectionRepository struct {
		findAllMethodCall              bool                      // True if we expect a call to the method
		findAllMethodCallReturnValue   []redirection.Redirection // Method return value
		findAllMethodCallReturnErr     bool                      // True if we expect the method to return an error
		findAllMethodCallReturnErrCode string                    // Expected error code
	}
	type testCase struct {
		test                        string
		expectRedirectionRepository testCaseRedirectionRepository
		expectCount                 int    // Expected number of results returned
		expectErr                   bool   // True if expecting error after query execution
		expectErrCode               string // Expected error code
	}

	// Create new test cases
	testCases := []testCase{
		{
			test: "Success",
			expectRedirectionRepository: testCaseRedirectionRepository{
				findAllMethodCall: true,
				findAllMethodCallReturnValue: []redirection.Redirection{
					{Key: "abcde", Location: "http://www.google.com"},
					{Key: "efghi", Location: "http://www.apple.com"},
				},
				findAllMethodCallReturnErr: false,
			},
			expectErr:   false,
			expectCount: 2,
		}, {
			test: "Success",
			expectRedirectionRepository: testCaseRedirectionRepository{
				findAllMethodCall:            true,
				findAllMethodCallReturnValue: []redirection.Redirection{},
				findAllMethodCallReturnErr:   false,
			},
			expectErr:   false,
			expectCount: 0,
		}, {
			test: "ErrInternal",
			expectRedirectionRepository: testCaseRedirectionRepository{
				findAllMethodCall:              true,
				findAllMethodCallReturnErr:     true,
				findAllMethodCallReturnErrCode: internal.ErrInternal,
			},
			expectErr:     true,
			expectErrCode: internal.ErrInternal,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			// Create a new context
			ctx := context.Background()

			// Create new gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create new mock repository
			redirectionRepository := mock.NewMockRedirectionRepository(ctrl)

			// Expect find by key method call
			if tc.expectRedirectionRepository.findAllMethodCall {
				// Expect error
				if tc.expectRedirectionRepository.findAllMethodCallReturnErr {
					err := &internal.Error{Code: tc.expectRedirectionRepository.findAllMethodCallReturnErrCode}
					redirectionRepository.EXPECT().FindAll(gomock.Any()).Return(nil, err)
				} else {
					redirectionRepository.EXPECT().FindAll(gomock.Any()).Return(tc.expectRedirectionRepository.findAllMethodCallReturnValue, nil)
				}
			}

			// Create new mock key generator
			keyGenerator := mock.NewMockKeyGenerator(ctrl)

			// Create new mock repository
			eventRepository := mock.NewMockEventRepository(ctrl)

			// Create new CreateRedirectionHandlerCommand handler
			app := application.New(redirectionRepository, eventRepository, keyGenerator)

			// Execute query and save result
			result, err := app.GetRedirectionList(ctx)

			// Check expected error
			if tc.expectErr {
				assert.Equal(t, tc.expectErrCode, internal.ErrorCode(err))
			} else {
				// CHeck result content
				assert.Nil(t, err)
				assert.Equal(t, tc.expectCount, len(result))
			}
		})
	}
}
