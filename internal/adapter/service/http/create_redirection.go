package http

import (
	"fmt"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/command"
	"github.com/gofiber/fiber/v2"
)

type CreateRedirectionRequest struct {
	Location string `json:"location" validate:"required,url"`
}

// CreateRedirectionHandler godoc
// @Summary      Create a new redirection
// @Accept       json
// @Param        location  body      CreateRedirectionRequest  true  "Location"
// @Success      202
// @Header       202      {string}  Location  "/key"
// @Failure      400      {object}  ErrorMessage
// @Failure      500      {object}  ErrorMessage
// @Router       /api [post]
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

	// Create a new CreateRedirectionCommand
	cmd := command.CreateRedirectionCommand{Location: request.Location}

	// Command execution
	value, err := s.app.Commands.CreateRedirection.Handle(c.Context(), cmd)
	if err != nil {
		return err
	}

	// Set location header
	c.Location(fmt.Sprintf("/%s", value.Key))

	// Send status created
	return c.SendStatus(fiber.StatusCreated)
}
