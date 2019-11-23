# grpc-rest-service

Hello, gRPC!

grpc-rest-service is an example service that serves a gRPC API, REST API, and OpenAPI documentation. This service can be used as a starting point for developing your own service using similar patterns. It also offers an opinionated project structure that makes it easy to iterate quickly and make changes to the service (explained in detail below).

## Usage

Run the server:

```sh
# Locally
go run cmd/grpc-rest-service/main.go

# Using docker-compose
docker-compose up --build
```

Run the client:

```sh
# gRPC
go run cmd/cli/main.go --name "Your name here"

# REST
curl localhost:8080/api/v1/greeter/hello -X POST -d '{"name": "Your name here"}'
```

View the OpenAPI spec:

- [localhost:8080/swagger.json](localhost:8080/swagger.json)
- [localhost:8080/docs](localhost:8080/docs)

## Project Layout

The project layout follows the [standard Go project layout](https://github.com/golang-standards/project-layout).

### [`/api`](/tree/master/api)

The API definitions for this service. The Protobuf `*.proto` files are written by hand, and the OpenAPI `*.swagger.json` files are generated.

### [`/cmd`](/tree/master/cmd)

`main.go` files for two Go executables:

- `grpc-rest-service`: the server that serves the gRPC API, REST API, and OpenAPI documentation
- `cli`: a CLI for making requests to the server

### [`/internal`](/tree/master/internal)

Packages used internally by the server:

- `openapi`: a small package for merging the generated OpenAPI specs so they can be served by the server as a single spec
- `server`: implementations of the gRPC services described in the Protobuf files, as well as handlers for the HTTP requests that the server receives
  - `greeter`: implementation of our `GreeterV1` service
  - `health`: implementation of Google's `grpc.health.v1.Health` service and our `HealthV1` service

### [`/pkg`](/tree/master/pkg)

Exported packages that can be used by other Go projects:

- `api`: the Go client and server code generated from our Protobufs -- the client code can be imported by other projects and used to call our service

## Questions

### Why gRPC?

In short, performance and code generation.

[gRPC](https://grpc.io/) is a high-performance RPC framework developed by Google. It uses HTTP2 and a binary data format called [Protocol Buffers (Protobuf)](https://developers.google.com/protocol-buffers) to achieve small payloads that can be transferred quickly over the wire. gRPC is a great fit for services that process many requests per second, but it also offers a good way to describe and generate APIs.

Protobufs are used to describe the RPC methods, requests, and responses of your gRPC service, similar to how [OpenAPI](https://swagger.io/docs/specification/about/) describes a REST API. From this you generate client and server code that handles the networking and routing of requests -- your only responsibility is to write the handlers for your gRPC service.

### Ok, but why not use OpenAPI?

If you're just looking for code generation then OpenAPI may already be a good solution. However, there are reasons to consider gRPC instead:

- **Higher performance and smaller payloads**, as described above.
- **A cleaner, more intuitive interface description language (IDL)**. This may land on personal preference, but [compare](tree/master/api) the `*.proto` specs with the `*.swagger.json` specs for this service -- which would you prefer working with?
- **Tighter integration with Go**. Although you can generate gRPC clients and servers for [multiple languages](https://grpc.io/docs/), it is particularly well integrated with Go and generates better code than the top OpenAPI code generators -- less boilerplate, more idiomatic.
- **You can still use REST and OpenAPI with gRPC**. As this project demonstrates, it's possible and in fact easy to serve a RESTful version of your gRPC service at the same time, as well as an OpenAPI spec generated from your Protobuf files.

### What are the downsides to gRPC?

As with every technology, there are tradeoffs:

- **Harder to debug network requests.** Because gRPC uses a binary data format, it's hard to debug raw network requests using something like a network inspector. Unlike JSON, the request and response payloads are not human-readable.

