package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/memory"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/mongo"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/redis"
	"github.com/emanuelefalzone/bitly/internal/adapter/service/grpc"
	"github.com/emanuelefalzone/bitly/internal/adapter/service/http"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/emanuelefalzone/bitly/internal/service"
)

var (
	Signals chan os.Signal
)

func main() {
	log.Println(`
 _     _ _   _       
| |   (_) | | |      
| |__  _| |_| |_   _ 
| '_ \| | __| | | | |
| |_) | | |_| | |_| |
|_.__/|_|\__|_|\__, |
                __/ |
               |___/ `)

	// Initialize channel
	Signals = make(chan os.Signal)

	// Subscribe to SIGINT signals
	signal.Notify(Signals, syscall.SIGINT, syscall.SIGTERM)

	// Read GRPC_PORT environment variable
	grpcPort, err := internal.GetEnv("GRPC_PORT")
	if err != nil {
		log.Fatal(err)
	}

	// Cast grpcPort to int
	intGrpcPort, err := strconv.Atoi(grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	// Read HTTP_PORT environment variable
	httpPort, err := internal.GetEnv("HTTP_PORT")
	if err != nil {
		log.Fatal(err)
	}

	// Cast httpPort to int
	intHTTPPort, err := strconv.Atoi(httpPort)
	if err != nil {
		log.Fatal(err)
	}

	// Create redirection repository
	var redirectionRepository redirection.Repository

	// Read REDIS_CONNECTION_STRING environment variable
	redisConnectionString, err := internal.GetEnv("REDIS_CONNECTION_STRING")

	// If a redis connection string has been specified
	if err == nil {
		// Create a redis repository
		redirectionRepository, err = redis.NewRedirectionRepository(redisConnectionString)
		if err != nil {
			log.Panic(err)
		}

		log.Println("Using redis repository.")
	} else {
		// Log error
		log.Println(err)

		// Otherwise use a memory repository
		redirectionRepository = memory.NewRedirectionRepository()
		log.Println("Using in memory repository.")
	}

	// Create event repository
	var eventRepository event.Repository

	// Read MONGO_CONNECTION_STRING environment variable
	mongoConnectionString, err := internal.GetEnv("MONGO_CONNECTION_STRING")

	// If a mongo connection string has been specified
	if err == nil {
		// Create a mongo repository
		eventRepository, err = mongo.NewEventRepository(mongoConnectionString)
		if err != nil {
			log.Panic(err)
		}

		log.Println("Using mongo repository.")
	} else {
		// Log error
		log.Println(err)

		// Otherwise use a memory repository
		eventRepository = memory.NewEventRepository()
		log.Println("Using in memory repository.")
	}

	// Create a key generator with random seed
	keyGenerator := service.NewRandomKeyGenerator(time.Now().Unix())

	// Create a new application
	app := application.New(redirectionRepository, eventRepository, keyGenerator)

	// Create a new grpc server
	grpcServer := grpc.NewServer(app)

	// Create a new http server
	httpServer := http.NewServer(app)

	// Start grpc server on specified port
	go func() {
		log.Printf("Starting GRPC server on port %d.\n", intGrpcPort)

		err := grpcServer.Start(intGrpcPort)
		if err != nil {
			log.Panic(err)
		}
	}()

	// Start http server on specified port
	go func() {
		log.Printf("Starting HTTP server on port %d.\n", intHTTPPort)

		err := httpServer.Start(intHTTPPort)
		if err != nil {
			log.Panic(err)
		}
	}()

	// Wait for SIGINT
	<-Signals

	log.Println()
	log.Println("Shutting down GRPC server.")

	// Gracefully shutdown grpc server
	grpcServer.Stop()

	// Gracefully shutdown http server
	log.Println("Shutting down HTTP server.")

	// Gracefully shutdown http server
	httpServer.Stop()
}
