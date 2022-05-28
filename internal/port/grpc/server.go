package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application"
	"github.com/emanuelefalzone/bitly/internal/port/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// Server associates a grpc server with an application
type Server struct {
	pb.UnimplementedBitlyServiceServer
	app        *application.Application
	grpcServer *grpc.Server
}

// NewServer create a new server with the given application
func NewServer(app *application.Application) *Server {
	return &Server{app: app}
}

// Start tries to start a new grpc server on the given port
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

// Stop gracefully stops the server
func (s *Server) Stop() {
	// Gracefully stop server
	s.grpcServer.GracefulStop()
}

func (s *Server) CreateRedirection(ctx context.Context,
	in *pb.CreateRedirectionRequest) (*pb.CreateRedirectionResponse, error) {
	// Command execution
	value, err := s.app.CreateRedirection(ctx, in.Location)
	if err != nil {
		return nil, mapErrorToGrpcError(err)
	}

	// Send redirection key
	return &pb.CreateRedirectionResponse{Key: value}, nil
}

func (s *Server) DeleteRedirection(ctx context.Context,
	in *pb.DeleteRedirectionRequest) (*pb.DeleteRedirectionResponse, error) {
	// Command execution
	err := s.app.DeleteRedirection(ctx, in.Key)
	if err != nil {
		return nil, mapErrorToGrpcError(err)
	}

	// Send DeleteRedirectionResponse to signal that the operation was successfully executed
	return &pb.DeleteRedirectionResponse{}, nil
}

func (s *Server) GetRedirectionLocation(ctx context.Context,
	in *pb.GetRedirectionLocationRequest) (*pb.GetRedirectionLocationResponse, error) {
	// Query execution
	value, err := s.app.GetRedirectionLocation(ctx, in.Key)
	if err != nil {
		return nil, mapErrorToGrpcError(err)
	}

	// Return redirection location
	return &pb.GetRedirectionLocationResponse{Location: value}, nil
}

func (s *Server) GetRedirectionCount(ctx context.Context,
	in *pb.GetRedirectionCountRequest) (*pb.GetRedirectionCountResponse, error) {
	// Query execution
	value, err := s.app.GetRedirectionCount(ctx, in.Key)
	if err != nil {
		return nil, mapErrorToGrpcError(err)
	}

	// Return redirection count
	return &pb.GetRedirectionCountResponse{Count: int64(value)}, nil
}

func (s *Server) GetRedirectionList(ctx context.Context,
	in *pb.GetRedirectionListRequest) (*pb.GetRedirectionListResponse, error) {
	// Query execution
	value, err := s.app.GetRedirectionList(ctx)
	if err != nil {
		return nil, mapErrorToGrpcError(err)
	}

	// Return redirection location
	return &pb.GetRedirectionListResponse{Keys: value}, nil
}

// Map internal errors to grpc error.
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
