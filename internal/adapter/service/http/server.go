package http

import (
	"fmt"
	"log"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	_ "github.com/emanuelefalzone/bitly/docs"
)

// @title     Bitly API
// @version   1.0.0
// @host      localhost:7070
// @BasePath  /
type Server struct {
	app       *application.Application
	validator validator.Validate
	server    *fiber.App
}

func NewServer(app *application.Application) *Server {
	return &Server{app: app, validator: *validator.New()}
}

func (s *Server) Start(port int) error {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())

	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Post("/api", s.CreateRedirectionHandler)
	app.Delete("/api/:key", s.DeleteRedirectionHandler)
	app.Get("/api/:key/count", s.RedirectionCountHandler)
	app.Get("/:key", s.RedirectionLocationHandler)

	s.server = app
	return app.Listen(fmt.Sprintf(":%d", port))
}

func (s *Server) Stop() {
	s.server.Shutdown()
}

type ErrorMessage struct {
	Message string `json:"error"`
} //@name Error

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Compute the human readable version of the error
	msg := ErrorMessage{Message: internal.ErrorMessage(err)}
	log.Printf("Error - %s", err)
	// Map the internal error code to http status code
	switch internal.ErrorCode(err) {
	case internal.ErrConflict:
		c.Status(fiber.StatusConflict).JSON(msg)
	case internal.ErrNotFound:
		c.SendStatus(fiber.StatusNotFound)
	case internal.ErrInvalid:
		c.Status(fiber.StatusBadRequest).JSON(msg)
	case internal.ErrInternal:
		c.SendStatus(fiber.StatusInternalServerError)
	}

	// Return from handler
	return nil
}
