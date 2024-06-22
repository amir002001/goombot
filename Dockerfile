# Stage 1: Build the Go binary
FROM golang:1.22.4-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp

# Stage 2: Create a minimal image with the binary and CA certificates
FROM alpine:latest AS ca-certificates

RUN apk --no-cache add ca-certificates

# Stage 3: Create a minimal image with the binary
FROM scratch

COPY --from=ca-certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from the builder stage
COPY --from=builder /app/myapp /

# Command to run the binary
CMD ["/myapp"]
