FROM docker.io/library/golang:1.22.0-alpine as builder

# Add the source code
COPY . go/src
WORKDIR go/src

# Download dependencies
RUN go mod download

# Build the Go application
RUN go mod vendor

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -mod=vendor -o ./go-calc-demo

FROM registry.access.redhat.com/ubi9/ubi:latest

COPY --from=builder /go/go/src/go-calc-demo /go/src/go-calc-demo

# Define the entry point for the container
CMD ["/go/src/go-calc-demo"]