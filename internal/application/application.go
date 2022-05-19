package application

import (
	"github.com/emanuelefalzone/bitly/internal/application/command"
	"github.com/emanuelefalzone/bitly/internal/application/query"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/emanuelefalzone/bitly/internal/service"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateRedirection command.CreateRedirectionHandler
	DeleteRedirection command.DeleteRedirectionHandler
}

type Queries struct {
	RedirectionLocation query.RedirectionLocationHandler
}

func New(redirections redirection.Repository, events event.Repository, generator service.KeyGenerator, dispatcher *event.Dispatcher) *Application {
	return &Application{
		Commands: Commands{
			CreateRedirection: command.NewCreateRedirectionHandler(redirections, generator, dispatcher),
			DeleteRedirection: command.NewDeleteRedirectionHandler(redirections, dispatcher),
		},
		Queries: Queries{
			RedirectionLocation: query.NewRedirectionLocationHandler(redirections, dispatcher),
		},
	}
}
