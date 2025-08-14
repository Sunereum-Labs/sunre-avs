# SunRe AVS Performer Docker Image
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Configure private repos
ENV GOPRIVATE=github.com/Layr-Labs/*

# Download dependencies
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/

# Build the performer
RUN go build -o performer ./cmd/main.go

# Runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/performer .

# Copy config and examples
COPY config/ ./config/
COPY examples/ ./examples/

# Expose performer port
EXPOSE 8080

# Set default environment
ENV PERFORMER_PORT=8080
ENV PERFORMER_TIMEOUT=5s
ENV ENV=production

# Run performer
CMD ["./performer"]
