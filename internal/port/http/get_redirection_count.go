package http

import (
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/gofiber/fiber/v2"
)

type redirectionCountRequest struct {
	Key string `validate:"min=1"`
}

// @Summary      Get the redirection count
// @Accept       json
// @Produce      json
// @Param        key  path      string  true  "Key"
// @Success      200 {object}  		redirectionCountRepresentation
// @Failure      404      {object}  errorMessage
// @Failure      500      {object}  errorMessage
// @Router       /api/redirection/{key}/count [get]
func (s *Server) redirectionCountHandler(c *fiber.Ctx) error {
	// Parse key param to create a redirectionCountRequest
	request := redirectionCountRequest{Key: c.Params("key")}

	// Validate the request
	err := internal.Validate(request)
	if err != nil {
		return err
	}

	// Query execution
	value, err := s.app.GetRedirectionCount(c.Context(), request.Key)
	if err != nil {
		return err
	}

	// Create and return encoded response
	response := getRedirectionCountRepresentation(request.Key, value)

	return c.JSON(response)
}
