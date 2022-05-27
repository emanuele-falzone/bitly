//go:build unit

package application_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/emanuelefalzone/bitly/test/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestApplicationQuery_RedirectionCount(t *testing.T) {
	// Build our needed testcase data struct
	type testCaseRedirectionRepository struct {
		findByKeyMethodCall              bool   // True if we expect a call to the method
		findByKeyMethodCallReturnErr     bool   // True if we expect the method to return an error
		findByKeyMethodCallReturnErrCode string // Expected error code
	}
	type testCaseEventRepository struct {
		findByRedirectionMethodCall                       bool   // True if we expect a call to the method
		findByRedirectionMethodCallReturnCreateEventCount int    // Number of create event returned
		findByRedirectionMethodCallReturnReadEventCount   int    // Number of read event returned
		findByRedirectionMethodCallReturnDeleteEventCount int    // Number of delete event returned
		findByRedirectionMethodCallReturnErr              bool   // True if we expect the method to return an error
		findByRedirectionMethodCallReturnErrCode          string // Expected error code
	}
	type testCase struct {
		test                        string
		location                    string // Location URL to be shortened
		key                         string // key associated to the redirection location
		expectRedirectionRepository testCaseRedirectionRepository
		expectEventRepository       testCaseEventRepository
		expectErr                   bool   // True if expecting error after query execution
		expectErrCode               string // Expected error code
		expectReadEventCount        int    // Expected number of read event
	}

	// Create new test cases
	testCases := []testCase{
		{
			test:     "Ten",
			location: "http://www.google.com",
			key:      "short",
			expectRedirectionRepository: testCaseRedirectionRepository{
				findByKeyMethodCall:          true,
				findByKeyMethodCallReturnErr: false,
			},
			expectEventRepository: testCaseEventRepository{
				findByRedirectionMethodCall:                       true,
				findByRedirectionMethodCallReturnCreateEventCount: 1,
				findByRedirectionMethodCallReturnReadEventCount:   10,
				findByRedirectionMethodCallReturnDeleteEventCount: 1,
				findByRedirectionMethodCallReturnErr:              false,
			},
			expectErr:            false,
			expectReadEventCount: 10,
		}, {
			test:     "Zero",
			location: "http://www.google.com",
			key:      "short",
			expectRedirectionRepository: testCaseRedirectionRepository{
				findByKeyMethodCall:          true,
				findByKeyMethodCallReturnErr: false,
			},
			expectEventRepository: testCaseEventRepository{
				findByRedirectionMethodCall:                       true,
				findByRedirectionMethodCallReturnCreateEventCount: 1,
				findByRedirectionMethodCallReturnReadEventCount:   0,
				findByRedirectionMethodCallReturnDeleteEventCount: 0,
				findByRedirectionMethodCallReturnErr:              false,
			},
			expectErr:            false,
			expectReadEventCount: 0,
		}, {
			test:     "ErrNotFound",
			location: "http://www.google.com",
			key:      "short",
			expectRedirectionRepository: testCaseRedirectionRepository{
				findByKeyMethodCall:              true,
				findByKeyMethodCallReturnErr:     true,
				findByKeyMethodCallReturnErrCode: internal.ErrNotFound,
			},
			expectEventRepository: testCaseEventRepository{
				findByRedirectionMethodCall: false,
			},
			expectErr:     true,
			expectErrCode: internal.ErrNotFound,
		}, {
			test:     "ErrNotFound2",
			location: "http://www.google.com",
			key:      "short",
			expectRedirectionRepository: testCaseRedirectionRepository{
				findByKeyMethodCall:          true,
				findByKeyMethodCallReturnErr: false,
			},
			expectEventRepository: testCaseEventRepository{
				findByRedirectionMethodCall:              true,
				findByRedirectionMethodCallReturnErr:     true,
				findByRedirectionMethodCallReturnErrCode: internal.ErrNotFound,
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
					redirectionRepository.EXPECT().FindByKey(gomock.Any(), gomock.Any()).Return(redirection.Redirection{}, err)
				} else {
					redirectionRepository.EXPECT().FindByKey(gomock.Any(), gomock.Any()).Return(redirection.Redirection{
						Key:      tc.key,
						Location: tc.location,
					}, nil)
				}
			}

			// Create new mock repository
			eventRepository := mock.NewMockEventRepository(ctrl)

			// Expect find by key method call
			if tc.expectEventRepository.findByRedirectionMethodCall {
				// Expect error
				if tc.expectEventRepository.findByRedirectionMethodCallReturnErr {
					err := &internal.Error{Code: tc.expectEventRepository.findByRedirectionMethodCallReturnErrCode}
					eventRepository.EXPECT().FindByRedirection(gomock.Any(), gomock.Any()).Return([]event.Event{}, err)
				} else {
					events := []event.Event{}
					for i := 0; i < tc.expectEventRepository.findByRedirectionMethodCallReturnCreateEventCount; i++ {
						events = append(events, event.Created(redirection.Redirection{Key: tc.key, Location: tc.location}))
					}
					for i := 0; i < tc.expectEventRepository.findByRedirectionMethodCallReturnReadEventCount; i++ {
						events = append(events, event.Read(redirection.Redirection{Key: tc.key, Location: tc.location}))
					}
					for i := 0; i < tc.expectEventRepository.findByRedirectionMethodCallReturnDeleteEventCount; i++ {
						events = append(events, event.Deleted(redirection.Redirection{Key: tc.key, Location: tc.location}))
					}
					eventRepository.EXPECT().FindByRedirection(gomock.Any(), gomock.Any()).Return(events, nil)
				}
			}

			// Create new mock key generator
			keyGenerator := mock.NewMockKeyGenerator(ctrl)

			// Create new CreateRedirectionHandlerCommand handler
			app := application.New(redirectionRepository, eventRepository, keyGenerator)

			// Execute query and save result
			result, err := app.GetRedirectionCount(ctx, tc.key)

			// Check expected error
			if tc.expectErr {
				assert.Equal(t, tc.expectErrCode, internal.ErrorCode(err))
			} else {
				// CHeck result content
				assert.Nil(t, err)
				assert.Equal(t, tc.expectReadEventCount, result)
			}
		})
	}
}
