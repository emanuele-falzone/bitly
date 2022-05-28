package http

import (
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/gofiber/fiber/v2"
)

type redirectionLocationRequest struct {
	Key string `validate:"min=1"`
}

// @Summary      Get redirection to a specific location
// @Accept       json
// @Produce      json
// @Param        key  path      string  true  "Key"
// @Success      302
// @Header       302      {string}  Location  "http://www.google.com"
// @Failure      404      {object}  errorMessage
// @Failure      500      {object}  errorMessage
// @Router       /{key} [get]
func (s *Server) redirectionLocationHandler(c *fiber.Ctx) error {
	// Parse key param to create a redirectionLocationRequest
	request := redirectionLocationRequest{Key: c.Params("key")}

	// Validate the request
	if err := internal.Validate(request); err != nil {
		return err
	}

	// Query execution
	value, err := s.app.GetRedirectionLocation(c.Context(), request.Key)
	if err != nil {
		return err
	}

	// Redirect to location
	return c.Redirect(value)
}
