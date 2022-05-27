package application

import (
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/emanuelefalzone/bitly/internal/service"
)

type Application struct {
	CreateRedirectionHandler   CreateRedirectionHandler
	DeleteRedirectionHandler   DeleteRedirectionHandler
	RedirectionLocationHandler RedirectionLocationHandler
	RedirectionCountHandler    RedirectionCountHandler
	RedirectionListHandler     RedirectionListHandler
}

func New(redirections redirection.Repository, events event.Repository, generator service.KeyGenerator, dispatcher *event.Dispatcher) *Application {
	return &Application{
		CreateRedirectionHandler:   NewCreateRedirectionHandler(redirections, generator, dispatcher),
		DeleteRedirectionHandler:   NewDeleteRedirectionHandler(redirections, dispatcher),
		RedirectionLocationHandler: NewRedirectionLocationHandler(redirections, dispatcher),
		RedirectionCountHandler:    NewRedirectionCountHandler(redirections, events),
		RedirectionListHandler:     NewRedirectionListHandler(redirections),
	}
}
