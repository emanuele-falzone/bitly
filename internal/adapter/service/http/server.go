package http

import (
	"fmt"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	// Import docs package
	_ "github.com/emanuelefalzone/bitly/docs"
)

// @title     Bitly API
// @version   1.0.0
// @host      localhost:7070
// @BasePath  /
type Server struct {
	app      *application.Application
	fiberApp *fiber.App
}

func NewServer(app *application.Application) *Server {
	return &Server{app: app}
}

func (s *Server) Start(port int) error {
	// Create new fiber app with custom error handler
	s.fiberApp = fiber.New(fiber.Config{ErrorHandler: errorHandler})

	// Middleware
	s.fiberApp.Use(recover.New())
	s.fiberApp.Use(cors.New())

	// Serve Swagger UI
	s.fiberApp.Get("/swagger/*", swagger.HandlerDefault)

	// Handle use cases
	s.fiberApp.Post("/api/redirection", s.createRedirectionHandler)
	s.fiberApp.Delete("/api/redirection/:key", s.deleteRedirectionHandler)
	s.fiberApp.Get("/api/redirection/:key/count", s.redirectionCountHandler)
	s.fiberApp.Get("/api/redirections", s.redirectionListHandler)
	s.fiberApp.Get("/:key", s.redirectionLocationHandler)

	return s.fiberApp.Listen(fmt.Sprintf(":%d", port))
}

func (s *Server) Stop() {
	// Gracefully stop server
	if err := s.fiberApp.Shutdown(); err != nil {
		panic(err)
	}
}

type errorMessage struct {
	Message string `json:"error"`
} // @name Error

func errorHandler(c *fiber.Ctx, err error) error {
	// Compute error message
	msg := internal.ErrorMessage(err)

	// Switch over error code
	switch internal.ErrorCode(err) {
	case internal.ErrConflict:
		return c.Status(fiber.StatusConflict).JSON(errorMessage{Message: msg})
	case internal.ErrNotFound:
		return c.SendStatus(fiber.StatusNotFound)
	case internal.ErrInvalid:
		return c.Status(fiber.StatusBadRequest).JSON(errorMessage{Message: msg})
	default:
		// Fallback to internal error
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}
