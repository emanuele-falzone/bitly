package http

import (
	"fmt"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application"
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
	application *application.Application
	fiberApp    *fiber.App
}

func NewServer(application *application.Application) *Server {
	return &Server{application: application}
}

func (s *Server) Start(port int) error {
	// Create new fiber app with custom error handler
	s.fiberApp = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})

	// Middleware
	s.fiberApp.Use(recover.New())
	s.fiberApp.Use(cors.New())

	// Serve Swagger UI
	s.fiberApp.Get("/swagger/*", swagger.HandlerDefault)

	// Handle use cases
	s.fiberApp.Post("/api/redirection", s.CreateRedirectionHandler)
	s.fiberApp.Delete("/api/redirection/:key", s.DeleteRedirectionHandler)
	s.fiberApp.Get("/api/redirection/:key/count", s.RedirectionCountHandler)
	s.fiberApp.Get("/:key", s.RedirectionLocationHandler)

	return s.fiberApp.Listen(fmt.Sprintf(":%d", port))
}

func (s *Server) Stop() {
	// Gracefully stop server
	s.fiberApp.Shutdown()
}

type ErrorMessage struct {
	Message string `json:"error"`
} //@name Error

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Compute error message
	msg := internal.ErrorMessage(err)

	// Switch over error code
	switch internal.ErrorCode(err) {
	case internal.ErrConflict:
		c.Status(fiber.StatusConflict).JSON(msg)
	case internal.ErrNotFound:
		c.SendStatus(fiber.StatusNotFound)
	case internal.ErrInvalid:
		c.Status(fiber.StatusBadRequest).JSON(msg)
	default:
		// Fallback to internal error
		c.SendStatus(fiber.StatusInternalServerError)
	}
	return nil
}
