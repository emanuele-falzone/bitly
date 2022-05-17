package driver

import (
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/internal/application/command"
	"github.com/emanuelefalzone/bitly/internal/application/query"
)

// The GoDriver directly interacts with the application
// This is the deepest driver we can use to test that the application really
// fulfills user requirements
type GoDriver struct {
	application *application.Application
}

func NewGoDriver(application *application.Application) Driver {
	return &GoDriver{application: application}
}

func (d *GoDriver) CreateRedirection(location string) (string, error) {
	// Create a new CreateRedirectionCommand
	cmd := command.CreateRedirectionCommand{Location: location}

	// Command execution
	value, err := d.application.Commands.CreateRedirection.Handle(cmd)
	if err != nil {
		return "", err
	}

	// Return key value
	return value.Key, nil
}
func (d *GoDriver) DeleteRedirection(key string) error {
	// Create a new DeleteRedirectionCommand useing th ekey specified in the request
	cmd := command.DeleteRedirectionCommand{Key: key}

	// Command execution
	err := d.application.Commands.DeleteRedirection.Handle(cmd)

	// Return operation error
	return err

}

func (d *GoDriver) GetRedirectionLocation(key string) (string, error) {
	// Create a new RedirectionLocationQuery
	q := query.RedirectionLocationQuery{Key: key}

	// Query execution
	value, err := d.application.Queries.RedirectionLocation.Handle(q)
	if err != nil {
		return "", err
	}

	// Return location
	return value.Location, nil
}
