package health

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/kevinsnydercodes/grpc-rest-service/pkg/api"
)

// HealthServer implements the gRPC health checking protocol.
type HealthServer struct{}

// Check implements the Check RPC of the gRPC health checking protocol.
func (o *HealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch implements the Watch RPC of the gRPC health checking protocol.
func (o *HealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

// NewHealthServer creates a new HealthServer.
func NewHealthServer() *HealthServer {
	return &HealthServer{}
}

// HealthV1Server implements the HealthV1 service.
type HealthV1Server struct{}

// Check implements the Check RPC of the HealthV1 service.
func (o *HealthV1Server) Check(ctx context.Context, req *api.HealthCheckRequest) (*api.HealthCheckResponse, error) {
	return &api.HealthCheckResponse{}, nil
}

// NewHealthV1Server creates a new HealthV1Server.
func NewHealthV1Server() *HealthV1Server {
	return &HealthV1Server{}
}
