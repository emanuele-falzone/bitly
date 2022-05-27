package http

import (
	"github.com/gofiber/fiber/v2"
)

// RedirectionListHandler godoc
// @Summary      Get the redirection list
// @Accept       json
// @Produce      json
// @Success      200 {object}  	   RedirectionListRepresentation
// @Failure      500      {object}  ErrorMessage
// @Router       /api/redirections [get]
func (s Server) RedirectionListHandler(c *fiber.Ctx) error {
	// Query execution
	value, err := s.app.GetRedirectionList(c.Context())
	if err != nil {
		return err
	}

	// Create and return encoded response
	response := getRedirectionListRepresentation(value)

	return c.JSON(response)
}
