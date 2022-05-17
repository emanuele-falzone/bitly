package command_test

import (
	"testing"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/command"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/emanuelefalzone/bitly/test/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteRedirection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// GIVEN
	redirectionRepository := mock.NewMockRedirectionRepository(ctrl)
	redirectionRepository.EXPECT().FindByKey(gomock.Any()).Return(redirection.New("abcdef", "http:/www.google.com"))
	redirectionRepository.EXPECT().Delete(gomock.Any()).Return(nil)

	handler := command.NewDeleteRedirectionHandler(redirectionRepository)

	// WHEN
	cmd := command.DeleteRedirectionCommand{Key: "abcdef"}
	err := handler.Handle(cmd)

	// THEN
	assert.Equal(t, nil, err)
}

func TestDeleteRedirection_NotFoundErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// GIVEN
	redirectionRepository := mock.NewMockRedirectionRepository(ctrl)
	redirectionRepository.EXPECT().FindByKey(gomock.Any()).Return(redirection.Redirection{}, &internal.Error{Code: internal.ErrNotFound})

	handler := command.NewDeleteRedirectionHandler(redirectionRepository)

	// WHEN
	cmd := command.DeleteRedirectionCommand{Key: "abcdef"}
	err := handler.Handle(cmd)

	// THEN
	assert.Equal(t, internal.ErrNotFound, internal.ErrorCode(err))
}

func TestDeleteRedirection_AnotherNotFoundErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// GIVEN
	redirectionRepository := mock.NewMockRedirectionRepository(ctrl)
	redirectionRepository.EXPECT().FindByKey(gomock.Any()).Return(redirection.New("abcdef", "http:/www.google.com"))
	redirectionRepository.EXPECT().Delete(gomock.Any()).Return(&internal.Error{Code: internal.ErrNotFound})

	handler := command.NewDeleteRedirectionHandler(redirectionRepository)

	// WHEN
	cmd := command.DeleteRedirectionCommand{Key: "abcdef"}
	err := handler.Handle(cmd)

	// THEN
	assert.Equal(t, internal.ErrNotFound, internal.ErrorCode(err))
}
