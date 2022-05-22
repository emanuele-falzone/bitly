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
	"github.com/emanuelefalzone/bitly/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// This serves as an end to end test for testing user requirements
func TestEndToEnd_GrpcServer(t *testing.T) {
	// Create a new context
	ctx := context.Background()

	// Define godog options
	var opts = godog.Options{
		Format:   "pretty",
		Output:   colors.Colored(os.Stdout),
		Paths:    []string{"../../../../test/feature"},
		TestingT: t,
	}

	// Read E2E_GRPC_SERVER environment variable
	serverAddress, err := internal.GetEnv("E2E_GRPC_SERVER") // ex: localhost:6060
	if err != nil {
		panic(err)
	}

	// Create new grpc driver
	driver_, err := NewGrpcDriver(serverAddress)
	if err != nil {
		panic(err)
	}

	// Run godog test suite
	status := godog.TestSuite{
		Name: "End to end tests using the grpc driver",
		ScenarioInitializer: test.Initialize(func() *test.Client {
			// Create a new client for each scenario (this allows to keep the client simple)
			return test.NewClient(driver_, ctx)
		}),
		Options: &opts,
	}.Run()

	// Check exit status
	if status != 0 {
		os.Exit(status)
	}
}

// The GrpcDriver interacts with the Grpc server
type GrpcDriver struct {
	client pb.BitlyServiceClient
}

func NewGrpcDriver(serverAddress string) (test.Driver, error) {
	// Connect to grpc server
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, &internal.Error{Code: internal.ErrInternal, Err: err}
	}

	// Create new grpc client
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

func (d *GrpcDriver) GetRedirectionList(ctx context.Context) ([]string, error) {
	response, err := d.client.GetRedirectionList(ctx, &pb.GetRedirectionListRequest{})
	if err != nil {
		return nil, err
	}
	return response.Keys, nil
}
