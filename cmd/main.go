package main

import (
	"log"
	"strconv"
	"time"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/persistence/redis"
	"github.com/emanuelefalzone/bitly/internal/adapter/service/grpc"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/internal/service"
)

func main() {
	log.Println(`Starting`)

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

	keyGenerator := service.NewRandomKeyGenerator(time.Now().Unix())

	app := application.New(redirectionRepository, keyGenerator)

	intGrpcPort, err := strconv.Atoi(grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(app)
	go grpcServer.Start(intGrpcPort)

	grpc.StartGateway()
}
