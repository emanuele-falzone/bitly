package acceptance_test

import (
	"context"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/test/acceptance/client"
	"github.com/emanuelefalzone/bitly/test/acceptance/driver"
	"github.com/emanuelefalzone/bitly/test/acceptance/scenario"
)

/*
This serves as an end to end test for testing user requirements
*/

func TestAcceptance_GrpcDriver_RedisRepository(t *testing.T) {

	ctx := context.Background()

	var opts = godog.Options{
		Format:   "pretty",
		Output:   colors.Colored(os.Stdout),
		Paths:    []string{"./feature"},
		TestingT: t,
	}

	serverAddress, err := internal.GetEnv("ACCEPTANCE_GRPC_SERVER")
	if err != nil {
		panic(err)
	}

	driver_, err := driver.NewGrpcDriver(serverAddress)
	if err != nil {
		panic(err)
	}

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
