package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"

	"github.com/kevinsnydercodes/grpc-rest-service/pkg/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	fAddress = flag.String("address", "localhost:8080", "Address of the grpc-rest-service to call.")
	fTLS     = flag.Bool("tls", false, "Whether to use a secure connection to the grpc-rest-service.")
	fName    = flag.String("name", "", "Name to send to the greeter service.")
)

func run() error {
	flag.Parse()

	// Create dial options
	options := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	if *fTLS {
		options = []grpc.DialOption{
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		}
	}

	// Create connection to gRPC service
	connection, err := grpc.Dial(*fAddress, options...)
	if err != nil {
		return errors.Wrap(err, "error dialing gRPC service")
	}

	// Create GreeterV1 client
	greeterV1Client := api.NewGreeterV1Client(connection)

	// Create and send request to GreeterV1 service
	ctx := context.Background()
	req := api.SayHelloRequest{
		Name: *fName,
	}
	resp, err := greeterV1Client.SayHello(ctx, &req)
	if err != nil {
		return err
	}

	// Print the greeting
	fmt.Println(resp.Greeting)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
