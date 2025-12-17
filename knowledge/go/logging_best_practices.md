# Logging Best Practices for ADK Agents

When building agents with the Google GenAI ADK for Go, logging user input requires balancing observability with security and privacy. This document outlines best practices to ensure your logs are useful, secure, and compliant.

## 1. Security & Privacy First

The most critical rule is to **never log sensitive data in plain text**. User input to agents often contains PII (Personally Identifiable Information), secrets, or proprietary data.

### Do Not Log

* **Passwords / API Keys**: Never log these.
* **PII**: Names, email addresses, phone numbers, credit card numbers.
* **Health/Financial Data**: Any sensitive personal data.

### Mitigation Strategies

* **Masking/Redaction**: Replace sensitive parts of the string with `****`.

    ```go
    // Bad
    slog.Info("User input", "input", "My password is secret123")

    // Good
    slog.Info("User input", "input", "My password is ****")
    ```

* **Hashing**: If you need to track unique values without revealing them, log a hash (e.g., SHA-256).
* **Allow-listing**: Only log specific fields or inputs that are known to be safe.

## 2. Structured Logging with `log/slog`

Use Go's standard `log/slog` package for structured logging. This makes logs machine-readable and easier to query in tools like Cloud Logging.

### Pattern

```go
import "log/slog"

// Log with context attributes
slog.Info("Agent received input",
    "agent_name", agent.Name(),
    "session_id", sessionID,
    "input_length", len(inputStr),
)
```

### Benefits

* **Queryable**: Filter by `agent_name` or `session_id`.
* **Consistent**: Enforces a schema for your logs.
* **Contextual**: Attach metadata without messy string formatting.

## 3. ADK Implementation Patterns

To log user input effectively in an ADK agent, you can use the **Wrapper Pattern**.

### The Wrapper Pattern

Wrap your agent instance with a custom struct that intercepts the `Run` method. This allows you to log the input before it reaches the core agent logic.

```go
type LoggingAgent struct {
    agent.Agent
}

func (l *LoggingAgent) Run(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
    // 1. Sanitize/Mask Input
    safeInput := sanitize(ctx.UserContent())

    // 2. Log Structured Info
    slog.Info("Agent Run Started",
        "agent", l.Name(),
        "input_snippet", safeInput, // Log only a safe snippet
    )

    // 3. Delegate to actual agent
    return l.Agent.Run(ctx)
}
```

## 4. Contextual Information

Always include context to make logs traceable.

* **Session ID**: Essential for tracing a conversation flow.
* **Request ID**: If served via HTTP, include the request ID to correlate with server logs.
* **User ID**: (If available and safe to log) To track user-specific issues.

## 5. Input Validation

Validate and sanitize input *before* logging.

* **Truncation**: Avoid logging massive payloads. Truncate long inputs to a reasonable limit (e.g., 100 chars) to prevent log flooding.
* **Sanitization**: Remove control characters or potential injection patterns if logs are viewed in a web UI.

## Summary Checklist

* [ ] Is PII masked or redacted?
* [ ] Are secrets excluded?
* [ ] Is structured logging (`log/slog`) used?
* [ ] Is the log message concise and searchable?
* [ ] Is context (Session ID, Agent Name) included?
