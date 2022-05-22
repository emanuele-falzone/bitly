package http

import (
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/command"
	"github.com/gofiber/fiber/v2"
)

type DeleteRedirectionRequest struct {
	Key string `validate:"min=1"`
}

// DeleteRedirectionHandler godoc
// @Summary      Delete the redirection associated with a specific key
// @Accept       json
// @Produce      json
// @Param        key  path      string  true  "Key"
// @Success      204
// @Failure      404      {object}  ErrorMessage
// @Router       /api/redirection/{key} [delete]
func (s Server) DeleteRedirectionHandler(c *fiber.Ctx) error {
	// Parse key param to create a DeleteRedirectionRequest
	request := DeleteRedirectionRequest{Key: c.Params("key")}

	// Validate the request
	err := internal.Validate(request)
	if err != nil {
		return err
	}

	// Create a new DeleteRedirectionCommand using the key specified in the request
	cmd := command.DeleteRedirectionCommand{Key: request.Key}

	// Command execution
	err = s.application.Commands.DeleteRedirection.Handle(c.Context(), cmd)
	if err != nil {
		return err
	}

	// Send status no content to signal that the operation was successfully executed
	return c.SendStatus(fiber.StatusNoContent)
}
