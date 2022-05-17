package command_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/command"
	"github.com/emanuelefalzone/bitly/test/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateRedirection(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// GIVEN
	redirectionRepository := mock.NewMockRedirectionRepository(ctrl)
	redirectionRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	keyGenerator := mock.NewMockKeyGenerator(ctrl)
	keyGenerator.EXPECT().NextKey(gomock.Any()).Return("abcdef")

	handler := command.NewCreateRedirectionHandler(redirectionRepository, keyGenerator)

	// WHEN
	cmd := command.CreateRedirectionCommand{Location: "http://www.google.com"}
	result, err := handler.Handle(ctx, cmd)

	// THEN
	assert.Equal(t, nil, err)
	assert.Equal(t, "abcdef", result.Key)
}

func TestCreateRedirection_InvalidErr(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// GIVEN
	redirectionRepository := mock.NewMockRedirectionRepository(ctrl)

	keyGenerator := mock.NewMockKeyGenerator(ctrl)
	keyGenerator.EXPECT().NextKey(gomock.Any()).Return("abcdef")

	handler := command.NewCreateRedirectionHandler(redirectionRepository, keyGenerator)

	// WHEN
	cmd := command.CreateRedirectionCommand{Location: "google.com"}
	_, err := handler.Handle(ctx, cmd)

	// THEN
	assert.Equal(t, internal.ErrInvalid, internal.ErrorCode(err))
}

func TestCreateRedirection_ConflictErr(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// GIVEN
	redirectionRepository := mock.NewMockRedirectionRepository(ctrl)
	redirectionRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&internal.Error{Code: internal.ErrConflict})

	keyGenerator := mock.NewMockKeyGenerator(ctrl)
	keyGenerator.EXPECT().NextKey(gomock.Any()).Return("abcdef")

	handler := command.NewCreateRedirectionHandler(redirectionRepository, keyGenerator)

	// WHEN
	cmd := command.CreateRedirectionCommand{Location: "http://www.google.com"}
	_, err := handler.Handle(ctx, cmd)

	// THEN
	assert.Equal(t, internal.ErrConflict, internal.ErrorCode(err))
}
