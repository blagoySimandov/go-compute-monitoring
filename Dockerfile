# Build stage
FROM golang:1.23.1 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o matrix-compute ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/matrix-compute .

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./matrix-compute"] 