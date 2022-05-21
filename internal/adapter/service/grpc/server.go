//go:build e2e

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
	application *application.Application
	grpcServer  *grpc.Server
}

func NewServer(application *application.Application) *Server {
	return &Server{application: application}
}

func (s *Server) Start(port int) error {
	// Announce on the local network
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	// Create new grpc server
	s.grpcServer = grpc.NewServer()

	// Register for bitly service
	pb.RegisterBitlyServiceServer(s.grpcServer, s)

	// Register reflection service on gRPC server.
	reflection.Register(s.grpcServer)

	// Serve
	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop() {
	// Gracefully stop server
	s.grpcServer.GracefulStop()
}

func (s *Server) CreateRedirection(ctx context.Context, in *pb.CreateRedirectionRequest) (*pb.CreateRedirectionResponse, error) {
	// Create a new CreateRedirectionCommand
	cmd := command.CreateRedirectionCommand{Location: in.Location}

	// Command execution
	value, err := s.application.Commands.CreateRedirection.Handle(ctx, cmd)
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
	err := s.application.Commands.DeleteRedirection.Handle(ctx, cmd)
	if err != nil {
		return nil, mapErrorToGrpcError(err)
	}

	// Send DeleteRedirectionResponse to signal that the operation was successfully executed
	return &pb.DeleteRedirectionResponse{}, nil
}

func (s *Server) GetRedirectionLocation(ctx context.Context, in *pb.GetRedirectionLocationRequest) (*pb.GetRedirectionLocationResponse, error) {
	// Create a new RedirectionLocationQuery
	q := query.RedirectionLocationQuery{Key: in.Key}

	// Query execution
	value, err := s.application.Queries.RedirectionLocation.Handle(ctx, q)
	if err != nil {
		return nil, mapErrorToGrpcError(err)
	}

	// Return redirection location
	return &pb.GetRedirectionLocationResponse{Location: value.Location}, nil
}

func (s *Server) GetRedirectionCount(ctx context.Context, in *pb.GetRedirectionCountRequest) (*pb.GetRedirectionCountResponse, error) {
	// Create a new RedirectionCountQuery
	q := query.RedirectionCountQuery{Key: in.Key}

	// Query execution
	value, err := s.application.Queries.RedirectionCount.Handle(ctx, q)
	if err != nil {
		return nil, mapErrorToGrpcError(err)
	}

	// Return redirection count
	return &pb.GetRedirectionCountResponse{Count: int64(value.Count)}, nil
}

// Map internal errors to grpc error
func mapErrorToGrpcError(err error) error {
	// Compute error message
	msg := internal.ErrorMessage(err)

	// Switch over error code
	switch internal.ErrorCode(err) {
	case internal.ErrInvalid:
		return status.Error(codes.InvalidArgument, msg)
	case internal.ErrNotFound:
		return status.Error(codes.NotFound, "")
	case internal.ErrConflict:
		return status.Error(codes.AlreadyExists, "")
	default:
		// Fallback to internal error
		return status.Errorf(codes.Internal, msg)
	}
}
