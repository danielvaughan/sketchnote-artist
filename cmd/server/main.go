package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/artifact"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/memory"
	"google.golang.org/adk/server/adkrest"
	"google.golang.org/adk/session"

	"github.com/joho/godotenv"

	"github.com/danielvaughan/sketchnote-artist/internal/app"
)

func main() {
	ctx := context.Background()

	// Initialize structured logging to stdout for server
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load .env file
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../../.env"); err != nil {
			slog.Warn("No .env file found in current directory or project root")
		}
	}

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		slog.Error("GOOGLE_API_KEY not set in environment or .env file")
		os.Exit(1)
	}

	// Create the Sketchnote Agent
	agentInstance, err := app.NewSketchnoteAgent(ctx, apiKey)
	if err != nil {
		slog.Error("Failed to create agent", "error", err)
		os.Exit(1)
	}

	// Configure the launcher with in-memory services
	config := &launcher.Config{
		AgentLoader:     agent.NewSingleLoader(agentInstance),
		SessionService:  session.InMemoryService(),
		ArtifactService: artifact.InMemoryService(),
		MemoryService:   memory.InMemoryService(),
	}

	// Create the REST handler
	handler := adkrest.NewHandler(config)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Starting REST server", "port", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
