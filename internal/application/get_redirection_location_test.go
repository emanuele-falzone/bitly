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

func TestApplicationQuery_RedirectionLocation(t *testing.T) {
	// Build our needed testcase data struct
	type testCaseRedirectionRepository struct {
		findByKeyMethodCall              bool   // True if we expect a call to the method
		findByKeyMethodCallReturnErr     bool   // True if we expect the method to return an error
		findByKeyMethodCallReturnErrCode string // Expected error code
	}
	type testCaseEventRepository struct {
		createMethodCall              bool   // True if we expect a call to the method
		createMethodCallReturnErr     bool   // True if we expect the method to return an error
		createMethodCallReturnErrCode string // Expected error code
	}
	type testCase struct {
		test                        string
		location                    string // Location URL to be shortened
		key                         string // key associated to the redirection location
		expectRedirectionRepository testCaseRedirectionRepository
		expectEventRepository       testCaseEventRepository
		expectErr                   bool   // True if expecting error after query execution
		expectErrCode               string // Expected error code
	}

	// Create new test cases
	testCases := []testCase{
		{
			test:     "Success",
			location: "http://www.google.com",
			key:      "short",
			expectRedirectionRepository: testCaseRedirectionRepository{
				findByKeyMethodCall:          true,
				findByKeyMethodCallReturnErr: false,
			},
			expectEventRepository: testCaseEventRepository{
				createMethodCall:          true,
				createMethodCallReturnErr: false,
			},
			expectErr: false,
		}, {
			test:     "ErrNotFound",
			location: "http://www.google.com",
			key:      "short",
			expectRedirectionRepository: testCaseRedirectionRepository{
				findByKeyMethodCall:              true,
				findByKeyMethodCallReturnErr:     true,
				findByKeyMethodCallReturnErrCode: internal.ErrNotFound,
			},
			expectErr:     true,
			expectErrCode: internal.ErrNotFound,
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
			if tc.expectRedirectionRepository.findByKeyMethodCall {
				// Expect error
				if tc.expectRedirectionRepository.findByKeyMethodCallReturnErr {
					err := &internal.Error{Code: tc.expectRedirectionRepository.findByKeyMethodCallReturnErrCode}
					redirectionRepository.EXPECT().FindByKey(gomock.Any(), gomock.Any()).Return(nil, err)
				} else {
					redirectionRepository.EXPECT().FindByKey(gomock.Any(), gomock.Any()).Return(&redirection.Redirection{
						Key:      tc.key,
						Location: tc.location,
					}, nil)
				}
			}

			// Create new mock key generator
			keyGenerator := mock.NewMockKeyGenerator(ctrl)

			// Create new mock repository
			eventRepository := mock.NewMockEventRepository(ctrl)

			// Expect create method call
			if tc.expectEventRepository.createMethodCall {
				// Expect error
				if tc.expectEventRepository.createMethodCallReturnErr {
					err := &internal.Error{Code: tc.expectEventRepository.createMethodCallReturnErrCode}
					eventRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(err)
				} else {
					eventRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				}
			}

			// Create new CreateRedirectionHandlerCommand handler
			app := application.New(redirectionRepository, eventRepository, keyGenerator)

			// Execute query and save result
			result, err := app.GetRedirectionLocation(ctx, tc.key)

			// Check expected error
			if tc.expectErr {
				assert.Equal(t, tc.expectErrCode, internal.ErrorCode(err))
			} else {
				// CHeck result content
				assert.Nil(t, err)
				assert.Equal(t, tc.location, result)
			}
		})
	}
}
