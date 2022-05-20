package http

import (
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/query"
	"github.com/gofiber/fiber/v2"
)

type RedirectionCountRequest struct {
	Key string `validate:"min=1"`
}

type RedirectionCountResponse struct {
	Count int `json:"count"`
}

// RedirectionCountHandler godoc
// @Summary      Get the redirection count
// @Param        key  path      string  true  "Location"
// @Success      200 {object}  		RedirectionCountResponse
// @Failure      404      {object}  ErrorMessage
// @Router       /api/{key}/count [get]
func (s Server) RedirectionCountHandler(c *fiber.Ctx) error {
	// Parse key param to create a DeleteRedirectionRequest
	request := RedirectionCountRequest{Key: c.Params("key")}

	// Validate the request
	err := internal.Validate(request)
	if err != nil {
		return err
	}

	// Create a new RedirectionCountQuery
	q := query.RedirectionCountQuery{Key: request.Key}

	// Query execution
	value, err := s.app.Queries.RedirectionCount.Handle(c.Context(), q)
	if err != nil {
		return err
	}

	// Create and return encoded response
	response := RedirectionCountResponse{Count: value.Count}
	return c.JSON(response)
}
