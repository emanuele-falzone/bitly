package application_test

import (
	"testing"

	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/test/mock"
	"github.com/golang/mock/gomock"
)

func TestApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	redirectionRepository := mock.NewMockRedirectionRepository(ctrl)
	keyGenerator := mock.NewMockKeyGenerator(ctrl)

	application.New(redirectionRepository, keyGenerator)
}
