package application

import (
	"github.com/emanuelefalzone/bitly/internal/application/event"
	"github.com/emanuelefalzone/bitly/internal/application/redirection"
	"github.com/emanuelefalzone/bitly/internal/application/service"
)

type Application struct {
	redirectionRepository redirection.Repository
	eventRepository       event.Repository
	keyGenerator          service.KeyGenerator
}

func New(redirectionRepository redirection.Repository,
	eventRepository event.Repository,
	keyGenerator service.KeyGenerator) *Application {
	return &Application{
		redirectionRepository: redirectionRepository,
		eventRepository:       eventRepository,
		keyGenerator:          keyGenerator,
	}
}
