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

func TestRedirectionLocation(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// GIVEN
	redirectionRepository := mock.NewMockRedirectionRepository(ctrl)
	redirectionRepository.EXPECT().FindByKey(gomock.Any(), gomock.Any()).Return(redirection.New("abcdef", "http:/www.google.com"))

	dispatcher := event.NewDispatcher(ctx)

	handler := query.NewRedirectionLocationHandler(redirectionRepository, dispatcher)

	// WHEN
	query := query.RedirectionLocationQuery{Key: "abcdef"}
	result, err := handler.Handle(ctx, query)

	// THEN
	assert.Equal(t, nil, err)
	assert.Equal(t, "http:/www.google.com", result.Location)
}

func TestRedirectionLocation_NotFoundErr(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// GIVEN
	redirectionRepository := mock.NewMockRedirectionRepository(ctrl)
	redirectionRepository.EXPECT().FindByKey(gomock.Any(), gomock.Any()).Return(redirection.Redirection{}, &internal.Error{Code: internal.ErrNotFound})

	dispatcher := event.NewDispatcher(ctx)

	handler := query.NewRedirectionLocationHandler(redirectionRepository, dispatcher)

	// WHEN
	query := query.RedirectionLocationQuery{Key: "abcdef"}
	_, err := handler.Handle(ctx, query)

	// THEN
	assert.Equal(t, internal.ErrNotFound, internal.ErrorCode(err))
}
