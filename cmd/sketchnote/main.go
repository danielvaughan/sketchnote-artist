// Package main is the entry point for the Sketchnote Artist CLI.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"

	"github.com/joho/godotenv"

	"github.com/danielvaughan/sketchnote-artist/internal/app"
	"github.com/danielvaughan/sketchnote-artist/internal/observability"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	// Initialize structured logging to file
	logFile, err := os.OpenFile("sketchnote-artist.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			slog.Error("Failed to close log file", "error", err)
		}
	}()

	logger := slog.New(slog.NewTextHandler(logFile, nil))
	slog.SetDefault(logger)

	// Load .env file
	// Try loading from current directory first, then fallback to root
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../../.env"); err != nil {
			slog.Warn("No .env file found in current directory or project root")
		}
	}

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("GOOGLE_API_KEY not set in environment or .env file")
	}

	// Define console reporter
	consoleReporter := func(msg string, _ ...interface{}) {
		fmt.Printf("\n%s\n", msg)
	}

	// Inject reporter into context
	ctx = observability.WithStatusReporter(ctx, consoleReporter)

	// Create the Sketchnote Agent
	agentInstance, err := app.NewSketchnoteAgent(ctx, app.Config{
		APIKey: apiKey,
	})
	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
	}

	// Configure the launcher
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(agentInstance),
	}

	// Run the agent using the full launcher
	l := full.NewLauncher()
	if err := l.Execute(ctx, config, os.Args[1:]); err != nil {
		// Log specific run error
		slog.Error("Run failed", "error", err, "syntax", l.CommandLineSyntax())
		return err
	}

	return nil
}
