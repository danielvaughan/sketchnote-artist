# Go Logging Best Practices

## Structured Logging with `log/slog`

As of Go 1.21, the standard library includes `log/slog`, a package for structured logging. It is the recommended approach for modern Go applications, replacing the legacy `log` package for most use cases.

### Why Structured Logging?

*   **Machine Readable:** Logs can be output in JSON format, making them easy to parse, filter, and analyze by log management systems (e.g., Cloud Logging, Splunk, Datadog).
*   **Contextual Data:** You can attach key-value pairs (attributes) to log entries, providing context without complex string formatting.
*   **Levels:** Built-in support for log levels (Debug, Info, Warn, Error) allows for controlling verbosity.
*   **Performance:** Designed to be high-performance with low allocation overhead.

### Implementation Pattern

#### 1. Initialization

Initialize a global logger early in your application's lifecycle (e.g., in `main`). Choose a handler based on your environment:
*   **`TextHandler`**: Human-readable, key=value format. Great for local development and CLI tools.
*   **`JSONHandler`**: JSON format. Best for production environments where logs are ingested by a system.

```go
import (
    "log/slog"
    "os"
)

func main() {
    // For CLI/Local Dev:
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    
    // For Production/Server:
    // logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    // Set as default logger so slog.Info(), slog.Error() use it
    slog.SetDefault(logger)
}
```

#### 2. Usage

Use the top-level functions `slog.Info`, `slog.Warn`, `slog.Error`, and `slog.Debug`. Pass context as alternating key-value pairs.

```go
// Simple message
slog.Info("Application started")

// With attributes
slog.Info("Processing video", "url", videoURL, "attempt", 1)

// Error handling
if err != nil {
    slog.Error("Failed to process video", "error", err, "video_id", id)
    os.Exit(1)
}
```

### Best Practices

1.  **Use Attributes, Not Formatting:** Avoid `slog.Info(fmt.Sprintf("User %s logged in", user))`. Instead, use `slog.Info("User logged in", "user", user)`. This keeps the message static and searchable, while the variable data is structured.
2.  **Consistent Keys:** Use consistent keys for attributes (e.g., always use `user_id`, not `uid` in one place and `userID` in another) to make searching easier.
3.  **Include Errors:** When logging an error, include the error object itself as an attribute (conventionally with the key "error").
4.  **Context Awareness:** `slog` supports `context.Context`. Use `slog.InfoContext` if you need to propagate trace IDs or other context-scoped values.

# Project Structure Best Practices

A well-structured Go project is easier to maintain, test, and scale. While Go doesn't enforce a strict layout, the community has converged on a set of standard patterns, often referred to as the "Standard Go Project Layout".

## 1. The Standard Layout

For production-grade applications, the following directory structure is widely accepted:

### `cmd/`
Contains the main entry points for your applications.
*   Each subdirectory represents a binary name (e.g., `cmd/server/`, `cmd/worker/`).
*   The `main.go` in these directories should be minimal: parse flags, configure the application, and call code in `internal/` or `pkg/`.
*   **Don't** put business logic here.

### `internal/`
Contains private application and library code.
*   The Go compiler enforces that code inside `internal/` cannot be imported by external repositories. This is perfect for your application's business logic.
*   **Subdirectories:** Group code by feature or layer (e.g., `internal/auth/`, `internal/database/`, `internal/api/`).

### `pkg/` (Optional)
Contains library code that is safe for external projects to use.
*   If your project is purely an application (not a library), you might not need this.
*   Modern Go advice often suggests defaulting to `internal/` unless you explicitly want to export a package.

### `configs/` (or `config/`)
Contains configuration files (e.g., `config.yaml`, `.env.example`) or configuration loading logic.

## 2. ADK & Agent-Specific Structure

When building agents with the Google GenAI ADK, you can adapt the standard layout to organize agent-specific components:

### `internal/agents/`
Define your agents here.
*   Each agent (e.g., `Summarizer`, `Artist`) can have its own package or file.
*   This keeps the agent configuration and prompt binding separate from the main application flow.

### `internal/tools/`
Implement your custom tools here.
*   Tools like `youtube_tool.go` or `image_tool.go` should reside here.
*   This promotes reusability if multiple agents need the same tool.

### `internal/flows/`
Define the orchestration logic.
*   If you have complex sequential or hierarchical flows (like `SketchnoteFlow`), define them here.
*   This separates the "wiring" of agents from the agents themselves.

### `internal/prompts/`
Store system instructions and prompts.
*   Instead of hardcoding strings in `agent.go`, use a dedicated package or text files (using `go:embed`).
*   This makes it easier to iterate on prompt engineering without touching code.

## 3. Evolution Strategy

**Start Simple:**
For a small script or proof-of-concept (like the initial Sketchnote Artist), a flat structure is fine:
```
.
├── main.go
├── agent.go
├── tools.go
└── go.mod
```

**Grow into Structure:**
As you add more agents, tools, or complex logic, refactor into the standard layout:
```
.
├── cmd/
│   └── sketchnote/
│       └── main.go
├── internal/
│   ├── agents/      # Agent definitions
│   ├── tools/       # Tool implementations
│   ├── flows/       # Orchestration
│   └── prompts/     # System instructions
├── go.mod
└── README.md
```
## 4. Testing Best Practices

### Integration Testing vs Unit Testing
In Go, it is important to distinguish between unit tests (fast, isolated, mocked) and integration tests (sclower, external dependencies, real I/O).

#### Implementation Pattern
Use build tags to separate integration tests. This allows you to run unit tests quickly by default and opt-in to integration tests.

1.  **Tagging:** Add `//go:build integration` at the very top of your integration test files.
2.  **Naming:** Naming the file `_integration_test.go` can also be helpful but the tag is the functional separator.
3.  **Running:**
    *   `go test ./...` will **skip** these files (if configured to exclude, or if you strictly require the tag). *Correction: By default, `//go:build integration` means the file is ONLY included if `-tags=integration` is passed.*
    *   `go test -tags=integration ./...` will include them.

Alternatively, use `testing.Short()`:

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }
    // ... test logic ...
}
```

**Recommendation:** For tests requiring API keys or external services (like Gemini or YouTube), use the `//go:build integration` tag to prevent accidental failures in CI environments lacking credentials.
