package http

import (
	"github.com/emanuelefalzone/bitly/internal/application/query"
	"github.com/gofiber/fiber/v2"
)

type RedirectionListResponse struct {
	Keys []string `json:"keys"`
} //@name Keys

// RedirectionListHandler godoc
// @Summary      Get the redirection list
// @Accept       json
// @Produce      json
// @Success      200 {object}  		RedirectionListResponse
// @Router       /api/redirections [get]
func (s Server) RedirectionListHandler(c *fiber.Ctx) error {
	// Create a new RedirectionListQuery
	q := query.RedirectionListQuery{}

	// Query execution
	value, err := s.application.Queries.RedirectionList.Handle(c.Context(), q)
	if err != nil {
		return err
	}

	// Create and return encoded response
	response := RedirectionListResponse{Keys: value.Keys}
	return c.JSON(response)
}
