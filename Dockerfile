# Build Stage
# Use Chainguard Go as the builder (secure, minimal, signed)
FROM cgr.dev/chainguard/go:latest AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the server binary
# CGO_ENABLED=0 ensures a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./cmd/server

# Create directories for output artifacts in builder (since distroless has no shell/mkdir)
RUN mkdir -p sketchnotes visual-briefs

# Runtime Stage
# Use Chainguard Static for 0-CVE guarantee and minimal runtime
FROM cgr.dev/chainguard/static:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=nonroot:nonroot /app/server ./server

# Copy web assets
COPY --chown=nonroot:nonroot web ./web

# Copy empty directories with correct ownership for the app to write to
COPY --from=builder --chown=nonroot:nonroot /app/sketchnotes ./sketchnotes
COPY --from=builder --chown=nonroot:nonroot /app/visual-briefs ./visual-briefs

# Set environment variables
ENV PORT=8080
ENV GO_ENV=production

# Expose the port
EXPOSE 8080

# Start the server
CMD ["./server"]
