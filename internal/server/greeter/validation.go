package greeter

import (
	"errors"

	"github.com/kevinsnydercodes/grpc-rest-service/pkg/api"
)

// SayHelloRequest is a type alias with validation.
type SayHelloRequest api.SayHelloRequest

// Validate validates a SayHelloRequest.
func (o *SayHelloRequest) Validate() error {
	if o.Name == "" {
		return errors.New("must provide name")
	}

	return nil
}
