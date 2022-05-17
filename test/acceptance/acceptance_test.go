package acceptance_test

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/memory"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/internal/service"
	"github.com/emanuelefalzone/bitly/test/acceptance/client"
	"github.com/emanuelefalzone/bitly/test/acceptance/driver"
)

func TestMain(m *testing.M) {

	var opts = godog.Options{
		Output: colors.Colored(os.Stdout),
		Paths:  []string{"./feature"},
	}

	godog.BindCommandLineFlags("godog.", &opts)

	status := godog.TestSuite{
		Name: "Acceptance tests using go driver and in memory repositories",
		ScenarioInitializer: InitializeScenario(func() *client.Client {
			redirectionRepository := memory.NewRedirectionRepository()
			keyGenerator := service.NewRandomKeyGenerator(0)
			application := application.New(redirectionRepository, keyGenerator)
			driver := driver.NewGoDriver(application)
			client := client.NewClient(driver)
			return client
		}),
		Options: &opts,
	}.Run()

	if status != 0 {
		os.Exit(status)
	}
}

func InitializeScenario(fn func() *client.Client) func(*godog.ScenarioContext) {
	return func(ctx *godog.ScenarioContext) {
		client := fn()

		// GIVEN
		ctx.Step(`^that I got a short link for (.*)$`, client.CreateRedirection)
		ctx.Step(`^that I got a short link that does not exist$`, client.GetNonExistingKey)
		ctx.Step(`^that the link has been visited (\d+) times$`, client.GetRedirectionLocationTimes)

		// WHEN
		ctx.Step(`^I command the system to shorten the link (.*)$`, client.CreateRedirection)
		ctx.Step(`^I command the system to delete the short link$`, client.DeleteRedirection)
		ctx.Step(`^I navigate to the short link$`, client.GetRedirectionLocation)

		// THEN
		ctx.Step(`^the system redirects me to (.*)$`, client.ConfirmLocationToBe)
		ctx.Step(`^the system returns a short link$`, client.ConfirmHasKey)
		ctx.Step(`^the system confirms that the operation was succesfully executed$`, client.ConfirmNoError)
		ctx.Step(`^the system signals that the short link does not exist$`, client.ConfirmError)
		ctx.Step(`^the system signals that the link is not valid$`, client.ConfirmError)
	}
}
