package http

import (
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/gofiber/fiber/v2"
)

type RedirectionLocationRequest struct {
	Key string `validate:"min=1"`
}

// RedirectionLocationHandler godoc
// @Summary      Get redirection to a specific location
// @Accept       json
// @Produce      json
// @Param        key  path      string  true  "Key"
// @Success      302
// @Header       302      {string}  Location  "http://www.google.com"
// @Failure      404      {object}  ErrorMessage
// @Failure      500      {object}  ErrorMessage
// @Router       /{key} [get]
func (s Server) RedirectionLocationHandler(c *fiber.Ctx) error {
	// Parse key param to create a DeleteRedirectionRequest
	request := RedirectionLocationRequest{Key: c.Params("key")}

	// Validate the request
	err := internal.Validate(request)
	if err != nil {
		return err
	}

	// Query execution
	value, err := s.application.GetRedirectionLocation(c.Context(), request.Key)
	if err != nil {
		return err
	}

	// Redirect to location
	return c.Redirect(value)
}
