package http

import (
	"github.com/emanuelefalzone/bitly/internal/application"
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
	// Create a new RedirectionListQuery
	q := application.RedirectionListQuery{}

	// Query execution
	value, err := s.application.RedirectionListHandler.Handle(c.Context(), q)
	if err != nil {
		return err
	}

	// Create and return encoded response
	response := getRedirectionListRepresentation(value.Keys)
	return c.JSON(response)
}
