package application_test

import (
	"context"
	"testing"

	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/test/mock"
	"github.com/golang/mock/gomock"
)

func TestApplication(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	redirectionRepository := mock.NewMockRedirectionRepository(ctrl)
	eventRepository := mock.NewMockEventRepository(ctrl)
	keyGenerator := mock.NewMockKeyGenerator(ctrl)

	dispatcher := event.NewDispatcher(ctx)

	application.New(redirectionRepository, eventRepository, keyGenerator, dispatcher)
}
