# Build Stage
# Use Chainguard Go as the builder (secure, minimal, signed)
FROM cgr.dev/chainguard/go:latest@sha256:86178b42db2e32763304e37f4cf3c6ec25b7bb83660dcb985ab603e3726a65a6 AS builder

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
FROM cgr.dev/chainguard/static:latest@sha256:1888f4db2c92e5a3e1b81952d8727e63c1b5b87ad3df374de318999beb4fd194

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
