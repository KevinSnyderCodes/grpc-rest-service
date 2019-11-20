package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/kevinsnydercodes/grpc-rest-service/internal/openapi"
	"github.com/kevinsnydercodes/grpc-rest-service/internal/server"
	"github.com/kevinsnydercodes/grpc-rest-service/internal/server/greeter"
	"github.com/kevinsnydercodes/grpc-rest-service/internal/server/health"

	"github.com/kevinsnydercodes/grpc-rest-service/pkg/api"
)

const title = "grpc-rest-service"

var (
	fPort             = flag.Uint("port", 8080, "Port for the server to listen on.")
	fOpenAPIFilesGlob = flag.String("openapi-files-glob", "api/*.swagger.json", "Glob to use for reading and merging OpenAPI files.")
)

func run() error {
	flag.Parse()

	if *fOpenAPIFilesGlob == "" {
		return errors.New("must provide OpenAPI files glob")
	}

	address := fmt.Sprintf(":%d", *fPort)

	// Merge OpenAPI files
	spec, err := openapi.MergeGlob(title, *fOpenAPIFilesGlob)
	if err != nil {
		return errors.Wrap(err, "error merging OpenAPI files")
	}

	// Create listener on port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrapf(err, "error listening on port %d", *fPort)
	}

	// Create gRPC server, register servers
	grpcServer := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewHealthServer())
	api.RegisterHealthV1Server(grpcServer, health.NewHealthV1Server())
	api.RegisterGreeterV1Server(grpcServer, greeter.NewGreeterV1Server())

	// Register grpc-gateway handlers
	ctx := context.Background()
	gwmux := runtime.NewServeMux()
	dialOptions := []grpc.DialOption{grpc.WithInsecure()}
	if err := api.RegisterHealthV1HandlerFromEndpoint(ctx, gwmux, address, dialOptions); err != nil {
		return errors.Wrap(err, "error registering HealthV1 grpc-gateway handler")
	}
	if err := api.RegisterGreeterV1HandlerFromEndpoint(ctx, gwmux, address, dialOptions); err != nil {
		return errors.Wrap(err, "error registering GreeterV1 grpc-gateway handler")
	}

	// Create and serve HTTP server
	httpServer := http.Server{
		Addr:    address,
		Handler: server.NewHandler(grpcServer, gwmux, spec),
	}
	fmt.Printf("Listening on port %d\n", *fPort)
	if err := httpServer.Serve(listener); err != nil {
		return errors.Wrap(err, "error serving HTTP server")
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
