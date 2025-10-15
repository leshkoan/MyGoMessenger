# Use the official Golang image to create a build artifact.
# This is a multi-stage build, so we use a temporary build stage.
FROM golang:1.23-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
# -o user-service: specifies the output file name
# CGO_ENABLED=0: disables CGO to create a statically linked binary
# GOOS=linux: specifies the target operating system
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service .

# Start a new stage from scratch for a smaller final image
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/user-service .

# Expose port 8001 to the outside world
EXPOSE 8001

# Command to run the executable
CMD ["./user-service"]
