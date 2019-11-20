package server

import (
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/go-openapi/runtime/middleware"

	"google.golang.org/grpc"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const basePath = "/"

// gatewayHandler creates a http.Handler that proxies all requests through a
// grpc-gateway multiplexer.
func gatewayHandler(gwmux *runtime.ServeMux) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	return mux
}

// openapiHandler creates a http.Handler that serves an OpenAPI spec and
// documentation for that spec.
//
// The handler serves two endpoints:
//
// - `GET /swagger.json`: the OpenAPI spec as JSON
// - `GET /docs`: documentation for the OpenAPI spec
//
// If none of these endpoints are matched, the request is passed onto the next
// handler.
func openapiHandler(spec []byte, next http.Handler) http.Handler {
	return middleware.Spec(basePath, spec, middleware.Redoc(middleware.RedocOpts{}, next))
}

// isGRPCRequest returns true if a request should be handled by the gRPC server.
//
// This is the recommended logic for distinguishing gRPC and other traffic:
// https://github.com/grpc/grpc-go/blob/24f6331/server.go#L745-L753
func isGRPCRequest(r *http.Request) bool {
	return r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc")
}

// grpcHandler creates a http.Handler that delegates to grpcServer on incoming
// gRPC connections and otherHandler otherwise.
func grpcHandler(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isGRPCRequest(r) {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

// NewHandler creates a http.Handler that serves traffic for all parts of the
// service.
func NewHandler(grpcServer *grpc.Server, gwmux *runtime.ServeMux, spec []byte) http.Handler {
	return grpcHandler(grpcServer, openapiHandler(spec, gatewayHandler(gwmux)))
}
