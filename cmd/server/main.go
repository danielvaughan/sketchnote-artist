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

	slog.Info("GOOGLE_API_KEY loaded successfully", "length", len(apiKey))

	// Create the Sketchnote Agent
	agentInstance, err := app.NewSketchnoteAgent(ctx, app.Config{
		APIKey: apiKey,
	})
	if err != nil {
		slog.Error("Failed to create agent", "error", err)
		os.Exit(1)
	}
	slog.Info("Agent created successfully")

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

	// Wrap the ADK handler with custom routing for UI and images
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve UI at root or /ui/
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			http.ServeFile(w, r, "web/index.html")
			return
		}

		// Serve static assets (css, js)
		if r.URL.Path == "/style.css" || r.URL.Path == "/app.js" {
			http.ServeFile(w, r, "web"+r.URL.Path)
			return
		}

		// Serve generated images
		if len(r.URL.Path) > 8 && r.URL.Path[:8] == "/images/" {
			// Serve from current directory, strip /images/ prefix
			filename := r.URL.Path[8:]
			// Basic security: prevent directory traversal
			if filename == "" || filename == "." || filename == ".." {
				http.NotFound(w, r)
				return
			}
			http.ServeFile(w, r, "sketchnotes/"+filename)
			return
		}

		// Fallback to ADK API handler
		handler.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(":"+port, finalHandler); err != nil {
		log.Fatal(err)
	}
}
