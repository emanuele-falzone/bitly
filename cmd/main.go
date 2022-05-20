package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/memory"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/mongo"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/redis"
	"github.com/emanuelefalzone/bitly/internal/adapter/service/grpc"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"github.com/emanuelefalzone/bitly/internal/service"
)

func main() {
	log.Println(`Starting`)

	ctx := context.Background()

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
	} else {
		// Otherwise use a memory repository
		log.Println(err)
		log.Println("Using in memory repository.")
		redirectionRepository = memory.NewRedirectionRepository()
	}

	// Create event repositoty
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
	} else {
		// Otherwise use a memory repository
		log.Println(err)
		log.Println("Using in memory repository.")
		eventRepository = memory.NewEventRepository()
	}

	// Create a key generator with random seed
	keyGenerator := service.NewRandomKeyGenerator(time.Now().Unix())

	// Create event logger to log event to stdout
	logger := service.NewEventLogger()

	// Create event logger to save event into event repository
	eventStore := service.NewEventStore(eventRepository)

	// Create new event dispatcher
	dispatcher := event.NewDispatcher(ctx)

	// Register logger into dispatcher
	dispatcher.Register(logger)

	// Register event store into dispatcher
	dispatcher.Register(eventStore)

	// Create a new aaplication
	app := application.New(redirectionRepository, eventRepository, keyGenerator, dispatcher)

	// Create a new grpc server
	grpcServer := grpc.NewServer(app)

	// Start grpc server on specified port
	err = grpcServer.Start(intGrpcPort)
	if err != nil {
		log.Fatal(err)
	}
}
