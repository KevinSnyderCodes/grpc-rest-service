package greeter

import (
	"context"
	"fmt"

	"github.com/kevinsnydercodes/grpc-rest-service/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GreeterV1Server implements the GreeterV1 service.
type GreeterV1Server struct{}

// SayHello implements the SayHello RPC of the GreeterV1 service.
func (o *GreeterV1Server) SayHello(ctx context.Context, req *api.SayHelloRequest) (*api.SayHelloResponse, error) {
	if err := (*SayHelloRequest)(req).Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	greeting := fmt.Sprintf("Hello, %s!", req.Name)

	return &api.SayHelloResponse{
		Greeting: greeting,
	}, nil
}

// NewGreeterV1Server creates a new GreeterV1Server.
func NewGreeterV1Server() *GreeterV1Server {
	return &GreeterV1Server{}
}
