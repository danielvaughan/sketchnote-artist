Here is the raw markdown content. You can save this directly as `context_reporting_pattern.md` or copy it into your agent's system prompt or knowledge base.

````markdown
# Pattern: Context-Based Status Reporting in Go

## Overview
In agentic workflows, "deep" components (like Tools, Plugins, or Sub-agents) often need to report progress to the user without being coupled to the output mechanism (CLI, WebSocket, HTTP Response). 

The **Context Injection** pattern solves this by embedding a `StatusReporter` function directly into the `context.Context`. This allows any function receiving the context to emit status updates safely, without dependencies on global loggers or specific UI implementations.

## Implementation

### 1. Definitions (The Plumbing)
Define these helpers in a shared utility package (e.g., `pkg/utils` or `pkg/observability`).

```go
package observability

import (
 "context"
 "fmt"
)

// StatusReporter is the function signature for sending updates.
type StatusReporter func(message string, details ...interface{})

// statusKey is a private key type to prevent collisions in the context.
type statusKey struct{}

// WithStatusReporter returns a new context containing the reporter.
func WithStatusReporter(ctx context.Context, reporter StatusReporter) context.Context {
 return context.WithValue(ctx, statusKey{}, reporter)
}

// Report sends a status update if a reporter is present in the context.
// It is safe to call even if no reporter is configured (no-op or fallback).
func Report(ctx context.Context, message string, details ...interface{}) {
 if reporter, ok := ctx.Value(statusKey{}).(StatusReporter); ok {
  reporter(message, details...)
 } else {
  // Optional: Fallback to standard log if debugging
  // fmt.Printf("[Fallback Log] " + message + "\n", details...)
 }
}
````

### 2. Usage in Tools (The Producer)

Tools and business logic do not need to know *how* the message is displayed. They simply call `Report`.

```go
func (t *MyFileTool) ReadFile(ctx context.Context, filename string) (string, error) {
    // Notify the user what is happening
    observability.Report(ctx, "Accessing file system", "file", filename)

    data, err := os.ReadFile(filename)
    if err != nil {
        return "", err
    }

    observability.Report(ctx, "File read successfully", "bytes", len(data))
    return string(data), nil
}
```

### 3. Usage at Entry Point (The Consumer)

The "Consumer" determines how those updates are actually presented to the user. This is usually done in your `main.go` or API handler.

#### Example A: CLI Output

```go
func main() {
    // Define how to handle reports (print to console)
    cliReporter := func(msg string, details ...interface{}) {
        fmt.Printf(">> [AGENT UPDATE]: %s %v\n", msg, details)
    }

    // Inject the reporter
    ctx := observability.WithStatusReporter(context.Background(), cliReporter)

    // Run the agent
    agent.Run(ctx, "Analyze the data")
}
```

#### Example B: WebSocket / HTTP Stream

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    // Define how to handle reports (stream JSON to client)
    webReporter := func(msg string, details ...interface{}) {
        update := map[string]interface{}{
            "type": "progress",
            "message": msg,
            "details": details,
        }
        json.NewEncoder(w).Encode(update)
        w.(http.Flusher).Flush() // Send immediately
    }

    ctx := observability.WithStatusReporter(r.Context(), webReporter)
    agent.Run(ctx, "Analyze the data")
}
```

## Benefits for Agents

1. **Decoupling:** Your deep logic (`ReadFile`) doesn't import `fmt` or `websocket`.
2. **Testability:** In unit tests, you can inject a "mock" reporter to verify that your agent is reporting the correct progress steps.
3. **Concurrency Safe:** Because `context` is immutable and passed down the stack, this works perfectly in concurrent Go routines.
4. **UI Agnostic:** The same agent code can power a CLI tool, a Slack bot, and a Web UI without changing a single line of business logic.

<!-- end list -->

```go
