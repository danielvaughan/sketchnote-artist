# Build Stage
# Use Golang 1.25 as the builder
FROM golang:1.25 AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the server binary
# CGO_ENABLED=0 ensures a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./cmd/server

# Runtime Stage
# Use Debian Slim for a small but standard runtime environment
FROM debian:bookworm-slim

# Install CA certificates to enable HTTPS calls
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server ./server

# Copy web assets
COPY web ./web

# Create directories for output artifacts
# Cloud Run is stateless, so these will be ephemeral unless mounted
RUN mkdir -p sketchnotes visual-briefs

# Set environment variables
ENV PORT=8080
ENV GO_ENV=production

# Expose the port
EXPOSE 8080

# Start the server
CMD ["./server"]
