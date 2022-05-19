package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/memory"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/redis"
	"github.com/emanuelefalzone/bitly/internal/adapter/service/grpc"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/service"
)

func main() {
	log.Println(`Starting`)

	ctx := context.Background()

	redisConnectionString, err := internal.GetEnv("REDIS_CONNECTION_STRING")
	if err != nil {
		log.Fatal(err)
	}

	grpcPort, err := internal.GetEnv("GRPC_PORT")
	if err != nil {
		log.Fatal(err)
	}

	redirectionRepository, err := redis.NewRedirectionRepository(redisConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	eventRepository := memory.NewEventRepository()

	keyGenerator := service.NewRandomKeyGenerator(time.Now().Unix())

	logger := service.NewEventLogger()
	eventStore := service.NewEventStore(eventRepository)

	dispatcher := event.NewDispatcher(ctx)
	dispatcher.Register(logger)
	dispatcher.Register(eventStore)

	app := application.New(redirectionRepository, eventRepository, keyGenerator, dispatcher)

	intGrpcPort, err := strconv.Atoi(grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(app)
	err = grpcServer.Start(intGrpcPort)
	if err != nil {
		log.Fatal(err)
	}
}
