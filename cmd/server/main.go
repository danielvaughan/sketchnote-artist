// Package main is the entry point for the Sketchnote Artist server.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/artifact"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/memory"
	"google.golang.org/adk/server/adkrest"
	"google.golang.org/adk/session"

	"github.com/joho/godotenv"

	"github.com/danielvaughan/sketchnote-artist/internal/app"
	"github.com/danielvaughan/sketchnote-artist/internal/storage"

	"google.golang.org/adk/session/database"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
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

	// Initialize Storage
	var store storage.Store
	if os.Getenv("DEPLOYMENT_MODE") == "cloud_run" {
		briefsBucket := os.Getenv("GCS_BUCKET_BRIEFS")
		imagesBucket := os.Getenv("GCS_BUCKET_IMAGES")
		slog.Info("Initializing Cloud Storage", "briefsBucket", briefsBucket, "imagesBucket", imagesBucket)
		gcsStore, err := storage.NewGCSStore(ctx, briefsBucket, imagesBucket)
		if err != nil {
			slog.Error("Failed to initialize GCS store", "error", err)
			os.Exit(1)
		}
		store = gcsStore
	} else {
		slog.Info("Initializing Local Disk Storage")
		store = &storage.DiskStore{}
	}

	// Create the Sketchnote Agent
	agentInstance, err := app.NewSketchnoteAgent(ctx, app.Config{
		APIKey: apiKey,
		Store:  store,
	})
	if err != nil {
		slog.Error("Failed to create agent", "error", err)
		os.Exit(1)
	}
	slog.Info("Agent created successfully")

	// Configure the launcher services
	var sessionService session.Service
	if os.Getenv("DEPLOYMENT_MODE") == "cloud_run" {
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASS")
		dbName := os.Getenv("DB_NAME")
		dbConn := os.Getenv("DB_CONNECTION_NAME")

		dsn := fmt.Sprintf("user=%s password=%s database=%s host=/cloudsql/%s", dbUser, dbPass, dbName, dbConn)
		slog.Info("Initializing PostgreSQL Session Service", "user", dbUser, "database", dbName, "connectionName", dbConn)
		dbService, err := database.NewSessionService(postgres.Open(dsn))
		if err != nil {
			slog.Error("Failed to initialize PostgreSQL session service", "error", err)
			os.Exit(1)
		}
		sessionService = dbService
	} else {
		slog.Info("Initializing SQLite Session Service")
		dbService, err := database.NewSessionService(sqlite.Open("sketchnote.db"))
		if err != nil {
			slog.Error("Failed to initialize SQLite session service", "error", err)
			os.Exit(1)
		}
		sessionService = dbService
	}

	// Ensure the database schema is up to date
	if err := database.AutoMigrate(sessionService); err != nil {
		slog.Error("Failed to run database auto-migration", "error", err)
		os.Exit(1)
	}

	config := &launcher.Config{
		AgentLoader:     agent.NewSingleLoader(agentInstance),
		SessionService:  sessionService,
		ArtifactService: artifact.InMemoryService(),
		MemoryService:   memory.InMemoryService(),
	}

	// Cache version at startup
	versionBytes, err := os.ReadFile("VERSION")
	var version string
	if err != nil {
		slog.Error("Failed to read VERSION file at startup", "error", err)
		version = "unknown"
	} else {
		version = strings.TrimSpace(string(versionBytes))
	}

	// Start the server
	handler := adkrest.NewHandler(config, 30*time.Second)

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

			// Stream content from storage (local or GCS)
			reader, err := store.Get(r.Context(), "images", filename)
			if err != nil {
				// If error is file not found, return 404
				if os.IsNotExist(err) {
					http.NotFound(w, r)
					return
				}
				slog.Error("Failed to retrieve image", "filename", filename, "error", err)
				http.Error(w, "Failed to retrieve image", http.StatusInternalServerError)
				return
			}
			defer func() {
				if err := reader.Close(); err != nil {
					slog.Error("Failed to close image reader", "error", err)
				}
			}()

			// Basic Content-Type sniffing or default to png
			// Since we know these are generated as PNGs usually:
			w.Header().Set("Content-Type", "image/png")
			if _, err := io.Copy(w, reader); err != nil {
				slog.Error("Failed to stream image content", "error", err)
			}
			return
		}

		// Serve version endpoint
		if r.URL.Path == "/version" {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]string{"version": version}); err != nil {
				slog.Error("Failed to write version response", "error", err)
			}
			return
		}

		// Fallback to ADK API handler
		handler.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           finalHandler,
		ReadHeaderTimeout: 3 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
