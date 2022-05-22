package test

import (
	"github.com/cucumber/godog"
)

func Initialize(fn func() *Client) func(*godog.ScenarioContext) {
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
		ctx.Step(`^I ask the system how many times the link has been visited$`, client.GetRedirectionCount)
		ctx.Step(`^I ask the system for the redirection list$`, client.GetRedirectionList)

		// THEN
		ctx.Step(`^the system redirects me to (.*)$`, client.ConfirmLocationToBe)
		ctx.Step(`^the system returns a short link$`, client.ConfirmHasKey)
		ctx.Step(`^the system confirms that the operation was successfully executed$`, client.ConfirmNoError)
		ctx.Step(`^the system signals that the short link does not exist$`, client.ConfirmError)
		ctx.Step(`^the system says the link has been visited (\d+) times$`, client.ConfirmCountToBe)
		ctx.Step(`^the system signals that the link is not valid$`, client.ConfirmError)
		ctx.Step(`^the system returns a list containing the short links$`, client.ConfirmCorrectList)
	}
}
