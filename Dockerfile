# Use the official Golang image for building the application
FROM golang:1.21 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o ssl-server-info main.go

# Create a minimal final image with the application and necessary certificates
FROM debian:bullseye-slim

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/ssl-server-info /app/ssl-server-info

# Copy SSL certificate and key (these can be replaced with real certs)
COPY server.crt server.key /app/

# Expose the HTTPS port (1443 for Docker container)
EXPOSE 1443

# Command to run the application
CMD ["/app/ssl-server-info", "-port=1443"]