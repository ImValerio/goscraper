# Use the official Golang image to create a build artifact.
FROM golang:1.21 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
RUN go mod tidy

# Copy all source code except main.go and miner.go into the container
COPY . .
COPY main.go main.go
COPY miner.go miner.go

# Build the Go app
RUN go build -o server .

# Start a new stage from scratch
FROM alpine:latest  

# Add ca-certificates and libc6-compat which is required for Go binaries that were built with CGO disabled
RUN apk --no-cache add ca-certificates libc6-compat

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/server .

# Expose port 5000 to the outside world
EXPOSE 5000

# Command to run the executable
CMD ["./server"]

