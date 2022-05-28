package http

import (
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/gofiber/fiber/v2"
)

type createRedirectionRequest struct {
	Location string `json:"location" validate:"required,url" example:"https://youtu.be/yhC-361QGJw"`
} // @name Location

// @Summary      Create a new redirection
// @Accept       json
// @Produce      json
// @Param        location body      createRedirectionRequest  true  "Location"
// @Success      202	  {object}  redirectionRepresentation
// @Header       202      {string}  Location  "/key"
// @Failure      400      {object}  errorMessage
// @Failure      500      {object}  errorMessage
// @Router       /api/redirection [post]
func (s *Server) createRedirectionHandler(c *fiber.Ctx) error {
	// Create a new createRedirectionRequest
	request := createRedirectionRequest{}

	// Parse the http request body into the request
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	// Validate the request
	if err := internal.Validate(request); err != nil {
		return err
	}

	// Command execution
	value, err := s.app.CreateRedirection(c.Context(), request.Location)
	if err != nil {
		return err
	}

	// Create response
	response := getRedirectionRepresentation(value)

	// Set location header
	c.Location(response.Links.Self.Href)

	// Set response status
	c.Status(fiber.StatusCreated)

	// Send response
	return c.JSON(response)
}
