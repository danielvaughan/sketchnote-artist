// Package observability provides tools for reporting agent status and metrics.
package observability

import (
	"context"
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
// It is safe to call even if no reporter is configured (no-op).
func Report(ctx context.Context, message string, details ...interface{}) {
	if reporter, ok := ctx.Value(statusKey{}).(StatusReporter); ok {
		reporter(message, details...)
	}
}
