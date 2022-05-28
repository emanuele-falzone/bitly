package http

import (
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/gofiber/fiber/v2"
)

type deleteRedirectionRequest struct {
	Key string `validate:"min=1"`
}

// @Summary      Delete the redirection associated with a specific key
// @Accept       json
// @Produce      json
// @Param        key  path      string  true  "Key"
// @Success      204
// @Failure      404      {object}  errorMessage
// @Failure      500      {object}  errorMessage
// @Router       /api/redirection/{key} [delete]
func (s *Server) deleteRedirectionHandler(c *fiber.Ctx) error {
	// Parse key param to create a deleteRedirectionRequest
	request := deleteRedirectionRequest{Key: c.Params("key")}

	// Validate the request
	if err := internal.Validate(request); err != nil {
		return err
	}

	// Command execution
	if err := s.app.DeleteRedirection(c.Context(), request.Key); err != nil {
		return err
	}

	// Send status no content to signal that the operation was successfully executed
	return c.SendStatus(fiber.StatusNoContent)
}
