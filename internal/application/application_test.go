//go:build acceptance

package application_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/memory"
	"github.com/emanuelefalzone/bitly/internal/application/command"
	"github.com/emanuelefalzone/bitly/internal/application/query"
	"github.com/emanuelefalzone/bitly/internal/service"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/test/acceptance/client"
	"github.com/emanuelefalzone/bitly/test/acceptance/driver"
	"github.com/emanuelefalzone/bitly/test/acceptance/scenario"
)

/*
This is the most important test.
It allows for testing the application directly using golang application code.
In this way we ensure that no business logic has leaked out to any service adapter.
*/
func TestAcceptance_Application(t *testing.T) {
	// Create a new context
	ctx := context.Background()

	// Define godog options
	var opts = godog.Options{
		Format:   "pretty",
		Output:   colors.Colored(os.Stdout),
		Paths:    []string{"../../test/acceptance/feature"},
		TestingT: t,
	}

	// Create a new in memory redirection repository
	redirectionRepository := memory.NewRedirectionRepository()

	// Create a new in memory event repository
	eventRepository := memory.NewEventRepository()

	// Create a new random key generator
	keyGenerator := service.NewRandomKeyGenerator(time.Now().Unix())

	// Create a new event dispatcher
	dispatcher := event.NewDispatcher(ctx)

	// Create a new event store
	eventStore := service.NewEventStore(eventRepository)

	// Register the event store as event listener
	dispatcher.Register(eventStore)

	// Create a new application
	application_ := application.New(redirectionRepository, eventRepository, keyGenerator, dispatcher)

	// Create a new go driver
	driver_ := NewGoDriver(application_)

	// Run godog test suite
	status := godog.TestSuite{
		Name: "Acceptance tests using go driver and in memory repository",
		ScenarioInitializer: scenario.Initialize(func() *client.Client {
			// Create a new client for each scenario (this allows to keep the client simple)
			return client.NewClient(driver_, ctx)
		}),
		Options: &opts,
	}.Run()

	if status != 0 {
		os.Exit(status)
	}
}

// The GoDriver directly interacts with the application
// This is the deepest driver we can use to test that the application really
// fulfills user requirements
type GoDriver struct {
	application *application.Application
}

func NewGoDriver(application *application.Application) driver.Driver {
	return &GoDriver{application: application}
}

func (d *GoDriver) CreateRedirection(ctx context.Context, location string) (string, error) {
	// Create a new CreateRedirectionCommand
	cmd := command.CreateRedirectionCommand{Location: location}

	// Command execution
	value, err := d.application.Commands.CreateRedirection.Handle(ctx, cmd)
	if err != nil {
		return "", err
	}

	// Return key value
	return value.Key, nil
}
func (d *GoDriver) DeleteRedirection(ctx context.Context, key string) error {
	// Create a new DeleteRedirectionCommand using the key specified in the request
	cmd := command.DeleteRedirectionCommand{Key: key}

	// Command execution
	err := d.application.Commands.DeleteRedirection.Handle(ctx, cmd)

	// Return operation error
	return err

}

func (d *GoDriver) GetRedirectionLocation(ctx context.Context, key string) (string, error) {
	// Create a new RedirectionLocationQuery
	q := query.RedirectionLocationQuery{Key: key}

	// Query execution
	value, err := d.application.Queries.RedirectionLocation.Handle(ctx, q)
	if err != nil {
		return "", err
	}

	// Return location
	return value.Location, nil
}

func (d *GoDriver) GetRedirectionCount(ctx context.Context, key string) (int, error) {
	// Create a new RedirectionCountQuery
	q := query.RedirectionCountQuery{Key: key}

	// Query execution
	value, err := d.application.Queries.RedirectionCount.Handle(ctx, q)
	if err != nil {
		return 0, err
	}

	// Return Count
	return value.Count, nil
}
