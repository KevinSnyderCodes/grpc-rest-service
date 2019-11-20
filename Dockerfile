FROM golang:1.13 as builder

WORKDIR /grpc-rest-service

COPY go.mod go.sum /grpc-rest-service/
RUN go mod download

# Build server
COPY cmd cmd
COPY internal internal
COPY pkg pkg

RUN CGO_ENABLED=0 GOOS=linux GOFLAGS=-ldflags=-w go build -o /go/bin/grpc-rest-service -ldflags=-s -v github.com/kevinsnydercodes/grpc-rest-service/cmd/grpc-rest-service

FROM alpine:latest

COPY --from=builder /go/bin/grpc-rest-service /usr/local/bin/grpc-rest-service
COPY api/*.swagger.json /etc/kevinsnydercodes/grpc-rest-service/

CMD ["/usr/local/bin/grpc-rest-service", "-openapi-files-glob", "/etc/kevinsnydercodes/grpc-rest-service/*.swagger.json"]
