package grpc

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./pb/bitly_service.proto

import (
	"context"
	"fmt"
	"net"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/adapter/service/grpc/pb"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/internal/application/command"
	"github.com/emanuelefalzone/bitly/internal/application/query"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedBitlyServiceServer
	app *application.Application
}

func NewServer(app *application.Application) *Server {
	return &Server{app: app}
}

func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterBitlyServiceServer(grpcServer, s)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	err = grpcServer.Serve(lis)

	return err
}

func (s *Server) CreateRedirection(ctx context.Context, in *pb.CreateRedirectionRequest) (*pb.CreateRedirectionResponse, error) {
	// Create a new CreateRedirectionCommand
	cmd := command.CreateRedirectionCommand{Location: in.Location}

	// Command execution
	value, err := s.app.Commands.CreateRedirection.Handle(ctx, cmd)
	if err != nil {
		return nil, mapErrorToGrpcError(err)
	}

	// Send redirection key
	return &pb.CreateRedirectionResponse{Key: value.Key}, nil
}

func (s *Server) DeleteRedirection(ctx context.Context, in *pb.DeleteRedirectionRequest) (*pb.DeleteRedirectionResponse, error) {
	// Create a new DeleteRedirectionCommand using the key specified in the request
	cmd := command.DeleteRedirectionCommand{Key: in.Key}

	// Command execution
	err := s.app.Commands.DeleteRedirection.Handle(ctx, cmd)
	if err != nil {
		return nil, mapErrorToGrpcError(err)
	}

	// Send DeleteRedirectionResponse to signal that the operation was succesfully executed
	return &pb.DeleteRedirectionResponse{}, nil
}

func (s *Server) GetRedirectionLocation(ctx context.Context, in *pb.GetRedirectionLocationRequest) (*pb.GetRedirectionLocationResponse, error) {
	// Create a new RedirectionLocationQuery
	q := query.RedirectionLocationQuery{Key: in.Key}

	// Query execution
	value, err := s.app.Queries.RedirectionLocation.Handle(ctx, q)
	if err != nil {
		return nil, mapErrorToGrpcError(err)
	}

	// Return redirection location
	return &pb.GetRedirectionLocationResponse{Location: value.Location}, nil
}

// Map internal errors to grpc error
func mapErrorToGrpcError(err error) error {
	msg := internal.ErrorMessage(err)
	switch internal.ErrorCode(err) {
	case internal.ErrInvalid:
		return status.Error(codes.InvalidArgument, msg)
	case internal.ErrNotFound:
		return status.Error(codes.NotFound, "")
	case internal.ErrConflict:
		return status.Error(codes.AlreadyExists, "")
	default:
		return status.Errorf(codes.Internal, msg)
	}
}
