package acceptance_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/memory"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/redis"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/service"
	"github.com/emanuelefalzone/bitly/test/acceptance/client"
	"github.com/emanuelefalzone/bitly/test/acceptance/driver"
	"github.com/emanuelefalzone/bitly/test/acceptance/scenario"
)

/*
This is the most important acceptance test.
It allows for testing the application directly using golang application code
In this way we ensure that no business logic has leaked out to any service adapter
*/

func TestAcceptance_GoDriver_InMemoryRepository(t *testing.T) {

	ctx := context.Background()

	var opts = godog.Options{
		Format:   "pretty",
		Output:   colors.Colored(os.Stdout),
		Paths:    []string{"./feature"},
		TestingT: t,
	}

	redirectionRepository := memory.NewRedirectionRepository()
	eventRepository := memory.NewEventRepository()

	keyGenerator := service.NewRandomKeyGenerator(time.Now().Unix())
	dispatcher := event.NewDispatcher(ctx)
	dispatcher.Register(service.NewEventStore(eventRepository))
	application_ := application.New(redirectionRepository, eventRepository, keyGenerator, dispatcher)
	driver_ := driver.NewGoDriver(application_)

	status := godog.TestSuite{
		Name: "Acceptance tests using go driver and in memory repository",
		ScenarioInitializer: scenario.Initialize(func() *client.Client {
			return client.NewClient(driver_, ctx)
		}),
		Options: &opts,
	}.Run()

	if status != 0 {
		os.Exit(status)
	}
}

func TestAcceptance_GoDriver_RedisRepository(t *testing.T) {

	ctx := context.Background()

	var opts = godog.Options{
		Format:   "pretty",
		Output:   colors.Colored(os.Stdout),
		Paths:    []string{"./feature"},
		TestingT: t,
	}

	connectionString, err := internal.GetEnv("ACCEPTANCE_REDIS_CONNECTION_STRING")
	if err != nil {
		panic(err)
	}

	redirectionRepository, err := redis.NewRedirectionRepository(connectionString)
	if err != nil {
		panic(err)
	}

	eventRepository := memory.NewEventRepository()

	keyGenerator := service.NewRandomKeyGenerator(time.Now().Unix())
	dispatcher := event.NewDispatcher(ctx)
	dispatcher.Register(service.NewEventStore(eventRepository))
	application_ := application.New(redirectionRepository, eventRepository, keyGenerator, dispatcher)
	driver_ := driver.NewGoDriver(application_)

	status := godog.TestSuite{
		Name: "Acceptance tests using go driver and redis repository",
		ScenarioInitializer: scenario.Initialize(func() *client.Client {
			return client.NewClient(driver_, ctx)
		}),
		Options: &opts,
	}.Run()

	if status != 0 {
		os.Exit(status)
	}
}
