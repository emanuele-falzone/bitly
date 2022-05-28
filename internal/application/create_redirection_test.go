//go:build unit

package application_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/test/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestApplicationCommand_CreateRedirection(t *testing.T) {
	// Build our needed testcase data struct
	type testCaseRedirectionRepository struct {
		createMethodCall              bool   // True if we expect a call to the method
		createMethodCallReturnErr     bool   // True if we expect the method to return an error
		createMethodCallReturnErrCode string // Expected error code
	}
	type testCaseEventRepository struct {
		createMethodCall              bool   // True if we expect a call to the method
		createMethodCallReturnErr     bool   // True if we expect the method to return an error
		createMethodCallReturnErrCode string // Expected error code
	}
	type testCaseKeyGenerator struct {
		nextKeyMethodCall            bool   // True if we expect a call to the method
		nextKeyMethodCallReturnValue string // Expected value returned by the method
	}
	type testCase struct {
		test                        string
		location                    string // Location URL to be shortened
		expectRedirectionRepository testCaseRedirectionRepository
		expectEventRepository       testCaseEventRepository
		expectKeyGenerator          testCaseKeyGenerator
		expectErr                   bool   // True if expecting error after command execution
		expectErrCode               string // Expected error code
	}

	// Create new test cases
	testCases := []testCase{
		{
			test:     "Success",
			location: "http://www.google.com",
			expectRedirectionRepository: testCaseRedirectionRepository{
				createMethodCall:          true,
				createMethodCallReturnErr: false,
			},
			expectEventRepository: testCaseEventRepository{
				createMethodCall:          true,
				createMethodCallReturnErr: false,
			},
			expectKeyGenerator: testCaseKeyGenerator{
				nextKeyMethodCall:            true,
				nextKeyMethodCallReturnValue: "short",
			},
			expectErr: false,
		},
		{
			test:     "ErrInvalid",
			location: "google",
			expectRedirectionRepository: testCaseRedirectionRepository{
				createMethodCall: false,
			},
			expectKeyGenerator: testCaseKeyGenerator{
				nextKeyMethodCall:            true,
				nextKeyMethodCallReturnValue: "short",
			},
			expectErr:     true,
			expectErrCode: internal.ErrInvalid,
		},
		{
			test:     "ErrConflict",
			location: "http://www.google.com",
			expectRedirectionRepository: testCaseRedirectionRepository{
				createMethodCall:              true,
				createMethodCallReturnErr:     true,
				createMethodCallReturnErrCode: internal.ErrConflict,
			},
			expectKeyGenerator: testCaseKeyGenerator{
				nextKeyMethodCall:            true,
				nextKeyMethodCallReturnValue: "short",
			},
			expectErr:     true,
			expectErrCode: internal.ErrConflict,
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

			// Expect create method call
			if tc.expectRedirectionRepository.createMethodCall {
				// Expect error
				if tc.expectRedirectionRepository.createMethodCallReturnErr {
					err := &internal.Error{Code: tc.expectRedirectionRepository.createMethodCallReturnErrCode}
					redirectionRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(err)
				} else {
					redirectionRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				}
			}

			// Create new mock key generator
			keyGenerator := mock.NewMockKeyGenerator(ctrl)

			// Expect next key method call
			if tc.expectKeyGenerator.nextKeyMethodCall {
				keyGenerator.EXPECT().NextKey().Return(tc.expectKeyGenerator.nextKeyMethodCallReturnValue)
			}

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

			// Execute command and save result
			result, err := app.CreateRedirection(ctx, tc.location)

			// Check expected error
			if tc.expectErr {
				assert.Equal(t, tc.expectErrCode, internal.ErrorCode(err))
			} else {
				// CHeck result content
				assert.Nil(t, err)
				assert.Equal(t, tc.expectKeyGenerator.nextKeyMethodCallReturnValue, result)
			}
		})
	}
}
