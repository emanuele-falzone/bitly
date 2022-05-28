package application

import (
	"github.com/emanuelefalzone/bitly/internal/application/event"
	"github.com/emanuelefalzone/bitly/internal/application/redirection"
	"github.com/emanuelefalzone/bitly/internal/application/service"
)

// Application exposes method to provides user specified features
type Application struct {
	redirectionRepository redirection.Repository
	eventRepository       event.Repository
	keyGenerator          service.KeyGenerator
}

// New creates an application with the given repositories and key generator
func New(redirectionRepository redirection.Repository,
	eventRepository event.Repository,
	keyGenerator service.KeyGenerator) *Application {
	return &Application{
		redirectionRepository: redirectionRepository,
		eventRepository:       eventRepository,
		keyGenerator:          keyGenerator,
	}
}
