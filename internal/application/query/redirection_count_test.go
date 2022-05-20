//go:build unit

package query_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/query"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/emanuelefalzone/bitly/test/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRedirectionCount(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// GIVEN
	redirectionRepository := mock.NewMockRedirectionRepository(ctrl)
	redirectionRepository.EXPECT().FindByKey(gomock.Any(), gomock.Any()).Return(redirection.New("abcdef", "http:/www.google.com"))

	eventRepository := mock.NewMockEventRepository(ctrl)
	eventRepository.EXPECT().FindByRedirection(gomock.Any(), gomock.Any()).Return(
		[]event.Event{
			event.Created(redirection.Redirection{}),
			event.Read(redirection.Redirection{}),
			event.Read(redirection.Redirection{}),
			event.Read(redirection.Redirection{}),
			event.Deleted(redirection.Redirection{}),
		}, nil)

	handler := query.NewRedirectionCountHandler(redirectionRepository, eventRepository)

	// WHEN
	query := query.RedirectionCountQuery{Key: "abcdef"}
	result, err := handler.Handle(ctx, query)

	// THEN
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, result.Count)
}

func TestRedirectionCount_NotFoundErr(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// GIVEN
	redirectionRepository := mock.NewMockRedirectionRepository(ctrl)
	redirectionRepository.EXPECT().FindByKey(gomock.Any(), gomock.Any()).Return(redirection.Redirection{}, &internal.Error{Code: internal.ErrNotFound})

	eventRepository := mock.NewMockEventRepository(ctrl)

	handler := query.NewRedirectionCountHandler(redirectionRepository, eventRepository)

	// WHEN
	query := query.RedirectionCountQuery{Key: "abcdef"}
	_, err := handler.Handle(ctx, query)

	// THEN
	assert.Equal(t, internal.ErrNotFound, internal.ErrorCode(err))
}

func TestRedirectionCount_AnotherNotFoundErr(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// GIVEN
	redirectionRepository := mock.NewMockRedirectionRepository(ctrl)
	redirectionRepository.EXPECT().FindByKey(gomock.Any(), gomock.Any()).Return(redirection.New("abcdef", "http:/www.google.com"))

	eventRepository := mock.NewMockEventRepository(ctrl)
	eventRepository.EXPECT().FindByRedirection(gomock.Any(), gomock.Any()).Return([]event.Event{}, &internal.Error{Code: internal.ErrNotFound})

	handler := query.NewRedirectionCountHandler(redirectionRepository, eventRepository)

	// WHEN
	query := query.RedirectionCountQuery{Key: "abcdef"}
	_, err := handler.Handle(ctx, query)

	// THEN
	assert.Equal(t, internal.ErrNotFound, internal.ErrorCode(err))
}
