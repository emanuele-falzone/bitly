package http

import (
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/query"
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
// @Router       /{key} [get]
func (s Server) RedirectionLocationHandler(c *fiber.Ctx) error {
	// Parse key param to create a DeleteRedirectionRequest
	request := RedirectionLocationRequest{Key: c.Params("key")}

	// Validate the request
	err := internal.Validate(request)
	if err != nil {
		return err
	}

	// Create a new RedirectionLocationQuery
	q := query.RedirectionLocationQuery{Key: request.Key}

	// Query execution
	value, err := s.application.Queries.RedirectionLocation.Handle(c.Context(), q)
	if err != nil {
		return err
	}

	// Redirect to location
	return c.Redirect(value.Location)
}
