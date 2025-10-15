# Use the official Golang image to create a build artifact.
# This is a multi-stage build, so we use a temporary build stage.
FROM golang:1.25-alpine as builder
RUN apk add protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN mkdir -p gen/go

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
COPY proto ./proto

RUN protoc -I /usr/include -I proto --go_out=. --go-grpc_out=. proto/*.proto
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
# -o notification-service: specifies the output file name
# CGO_ENABLED=0: disables CGO to create a statically linked binary
# GOOS=linux: specifies the target operating system
RUN CGO_ENABLED=0 GOOS=linux go build -o /notification-service ./notifications

# Start a new stage from scratch for a smaller final image
FROM alpine:latest

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /notification-service .

# Command to run the executable
CMD ["./notification-service"]