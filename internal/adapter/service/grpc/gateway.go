package grpc

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/emanuelefalzone/bitly/internal/adapter/service/grpc/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	grpcServerEndpoint = "localhost:%s"
)

func StartGateway() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithForwardResponseOption(responseMessageMatcher))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterBitlyServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(grpcServerEndpoint, "4000"), opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway service endpoint: %v", err)
	}
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Could not setup HTTP endpoint: %v", err)
	}
}

func responseMessageMatcher(ctx context.Context, w http.ResponseWriter, resp protoreflect.ProtoMessage) error {
	t, ok := resp.(*pb.GetRedirectionLocationResponse)
	if ok {
		w.Header().Set("Location", t.Location)
		w.WriteHeader(http.StatusFound)
	}
	return nil
}
