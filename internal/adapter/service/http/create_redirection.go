package http

import (
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/gofiber/fiber/v2"
)

type CreateRedirectionRequest struct {
	Location string `json:"location" validate:"required,url" example:"https://youtu.be/yhC-361QGJw"`
} //@name Location

// CreateRedirectionHandler godoc
// @Summary      Create a new redirection
// @Accept       json
// @Produce      json
// @Param        location body      CreateRedirectionRequest  true  "Location"
// @Success      202	  {object}  RedirectionRepresentation
// @Header       202      {string}  Location  "/key"
// @Failure      400      {object}  ErrorMessage
// @Failure      500      {object}  ErrorMessage
// @Router       /api/redirection [post]
func (s Server) CreateRedirectionHandler(c *fiber.Ctx) error {
	// Create a new CreateRedirectionRequest
	request := CreateRedirectionRequest{}

	// Parse the http request body into the request
	c.BodyParser(&request)

	// Validate the request
	err := internal.Validate(request)
	if err != nil {
		return err
	}

	// Command execution
	value, err := s.application.CreateRedirection(c.Context(), request.Location)
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
