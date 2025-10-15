# Use the official Golang image to create a build artifact.
FROM golang:1.25-alpine as builder

# Install necessary packages for building
RUN apk add --no-cache protobuf-dev git

# Install protobuf generators
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Create directory for generated code
RUN mkdir -p gen/go

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
COPY gen ./gen
COPY proto ./proto

# Generate protobuf code
RUN protoc -I /usr/include -I proto --go_out=. --go-grpc_out=. proto/*.proto
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -o /notification-service ./notifications

# Start a new stage from scratch for a smaller final image
FROM alpine:latest

# Add metadata labels
LABEL maintainer="MyGoMessenger Team"
LABEL service="notification-service"
LABEL version="1.0"

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -s /bin/sh -u 1000 -G appuser appuser

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /notification-service .

# Change ownership to non-root user
RUN chown appuser:appuser /notification-service

# Switch to non-root user
USER appuser

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Command to run the executable
CMD ["./notification-service"]