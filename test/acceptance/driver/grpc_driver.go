package driver

import (
	"context"
	"time"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/service/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// The GrpcDriver interacts with the Grpc server
type GrpcDriver struct {
	client pb.BitlyServiceClient
}

func NewGrpcDriver(serverAddress string) (Driver, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, &internal.Error{Code: internal.ErrInternal, Err: err}
	}

	client := pb.NewBitlyServiceClient(conn)

	return &GrpcDriver{client: client}, nil
}

func (d *GrpcDriver) CreateRedirection(ctx context.Context, location string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := d.client.CreateRedirection(ctx, &pb.CreateRedirectionRequest{Location: location})
	if err != nil {
		return "", err
	}
	return response.Key, nil
}
func (d *GrpcDriver) DeleteRedirection(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := d.client.DeleteRedirection(ctx, &pb.DeleteRedirectionRequest{Key: key})
	return err
}

func (d *GrpcDriver) GetRedirectionLocation(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := d.client.GetRedirectionLocation(ctx, &pb.GetRedirectionLocationRequest{Key: key})
	if err != nil {
		return "", err
	}
	return response.Location, nil
}

func (d *GrpcDriver) GetRedirectionCount(ctx context.Context, key string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := d.client.GetRedirectionCount(ctx, &pb.GetRedirectionCountRequest{Key: key})
	if err != nil {
		return 0, err
	}
	return int(response.Count), nil
}
