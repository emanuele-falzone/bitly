//go:build e2e

package grpc_test

import (
	"context"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/service/grpc/pb"
	"github.com/emanuelefalzone/bitly/test/acceptance/client"
	"github.com/emanuelefalzone/bitly/test/acceptance/driver"
	"github.com/emanuelefalzone/bitly/test/acceptance/scenario"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
This serves as an end to end test for testing user requirements
*/

func TestAcceptance_GrpcDriver_RedisRepository(t *testing.T) {

	ctx := context.Background()

	var opts = godog.Options{
		Format:   "pretty",
		Output:   colors.Colored(os.Stdout),
		Paths:    []string{"./feature"},
		TestingT: t,
	}

	serverAddress, err := internal.GetEnv("E2E_GRPC_SERVER")
	if err != nil {
		panic(err)
	}

	driver_, err := NewGrpcDriver(serverAddress)
	if err != nil {
		panic(err)
	}

	status := godog.TestSuite{
		Name: "Acceptance tests using go driver and redis repository",
		ScenarioInitializer: scenario.Initialize(func() *client.Client {
			return client.NewClient(driver_, ctx)
		}),
		Options: &opts,
	}.Run()

	if status != 0 {
		os.Exit(status)
	}
}

// The GrpcDriver interacts with the Grpc server
type GrpcDriver struct {
	client pb.BitlyServiceClient
}

func NewGrpcDriver(serverAddress string) (driver.Driver, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, &internal.Error{Code: internal.ErrInternal, Err: err}
	}

	client := pb.NewBitlyServiceClient(conn)

	return &GrpcDriver{client: client}, nil
}

func (d *GrpcDriver) CreateRedirection(ctx context.Context, location string) (string, error) {
	response, err := d.client.CreateRedirection(ctx, &pb.CreateRedirectionRequest{Location: location})
	if err != nil {
		return "", err
	}
	return response.Key, nil
}
func (d *GrpcDriver) DeleteRedirection(ctx context.Context, key string) error {
	_, err := d.client.DeleteRedirection(ctx, &pb.DeleteRedirectionRequest{Key: key})
	return err
}

func (d *GrpcDriver) GetRedirectionLocation(ctx context.Context, key string) (string, error) {
	response, err := d.client.GetRedirectionLocation(ctx, &pb.GetRedirectionLocationRequest{Key: key})
	if err != nil {
		return "", err
	}
	return response.Location, nil
}

func (d *GrpcDriver) GetRedirectionCount(ctx context.Context, key string) (int, error) {
	response, err := d.client.GetRedirectionCount(ctx, &pb.GetRedirectionCountRequest{Key: key})
	if err != nil {
		return 0, err
	}
	return int(response.Count), nil
}
