package http

import (
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/gofiber/fiber/v2"
)

type RedirectionCountRequest struct {
	Key string `validate:"min=1"`
}

// RedirectionCountHandler godoc
// @Summary      Get the redirection count
// @Accept       json
// @Produce      json
// @Param        key  path      string  true  "Key"
// @Success      200 {object}  		RedirectionCountRepresentation
// @Failure      404      {object}  ErrorMessage
// @Failure      500      {object}  ErrorMessage
// @Router       /api/redirection/{key}/count [get]
func (s Server) RedirectionCountHandler(c *fiber.Ctx) error {
	// Parse key param to create a DeleteRedirectionRequest
	request := RedirectionCountRequest{Key: c.Params("key")}

	// Validate the request
	err := internal.Validate(request)
	if err != nil {
		return err
	}

	// Query execution
	value, err := s.application.GetRedirectionCount(c.Context(), request.Key)
	if err != nil {
		return err
	}

	// Create and return encoded response
	response := getRedirectionCountRepresentation(request.Key, value)
	return c.JSON(response)
}
