# Use the official Golang image to create a build artifact.
FROM golang:1.25-alpine as builder
RUN apk add protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN mkdir -p gen/go

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to leverage Docker cache
COPY go.mod go.sum ./
COPY gen/go/users/go.mod ./gen/go/users/
COPY gen/go/messages/go.mod ./gen/go/messages/
COPY proto ./proto

RUN protoc -I /usr/include -I proto --go_out=. --go-grpc_out=. proto/*.proto
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app for the user service
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service ./users

# Start a new stage from scratch for a smaller final image
FROM alpine:latest

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /user-service .

# Expose port 8001 to the outside world
EXPOSE 8001

# Command to run the executable
CMD ["./user-service"]