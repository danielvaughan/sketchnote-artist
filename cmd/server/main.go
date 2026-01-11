// Package main is the entry point for the Sketchnote Artist server.
package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
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
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// gzipResponseWriter wraps http.ResponseWriter to provide gzip compression.
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// gzipWriterPool reuses gzip writers to reduce allocations.
var gzipWriterPool = sync.Pool{
	New: func() any {
		return gzip.NewWriter(nil)
	},
}

// gzipMiddleware compresses responses for clients that accept gzip encoding.
// It skips compression for SSE endpoints and already-compressed content types.
func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip compression for SSE endpoints (they need streaming without buffering)
		if strings.HasSuffix(r.URL.Path, "_sse") || strings.Contains(r.URL.Path, "/run_sse") {
			next.ServeHTTP(w, r)
			return
		}

		// Skip if client doesn't accept gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz := gzipWriterPool.Get().(*gzip.Writer)
		gz.Reset(w)
		defer func() {
			if err := gz.Close(); err != nil {
				slog.Error("Failed to close gzip writer", "error", err)
			}
			gzipWriterPool.Put(gz)
		}()

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length") // Length changes with compression

		next.ServeHTTP(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}

func main() {
	ctx := context.Background()

	// Initialize structured logging to stdout for server
	slogLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(slogLogger)

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
		dbService, err := database.NewSessionService(postgres.Open(dsn), &gorm.Config{
			Logger:      logger.Default.LogMode(logger.Silent),
			PrepareStmt: true, // Cache prepared statements for better performance
		})
		if err != nil {
			slog.Error("Failed to initialize PostgreSQL session service", "error", err)
			os.Exit(1)
		}

		// Configure connection pool for production
		if sqlDB, err := dbService.DB().DB(); err == nil {
			sqlDB.SetMaxOpenConns(25)               // Maximum open connections
			sqlDB.SetMaxIdleConns(10)               // Maximum idle connections
			sqlDB.SetConnMaxLifetime(5 * time.Minute) // Connection max lifetime
			sqlDB.SetConnMaxIdleTime(1 * time.Minute) // Idle connection max lifetime
			slog.Info("PostgreSQL connection pool configured", "maxOpen", 25, "maxIdle", 10)
		}

		sessionService = dbService
	} else {
		slog.Info("Initializing SQLite Session Service")
		dbService, err := database.NewSessionService(sqlite.Open("sketchnote.db"), &gorm.Config{
			Logger:      logger.Default.LogMode(logger.Silent),
			PrepareStmt: true, // Cache prepared statements for better performance
		})
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
	handler := adkrest.NewHandler(config, 5*time.Minute)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Starting REST server", "port", port)

	// Wrap the ADK handler with custom routing for UI and images
	routingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve UI at root or /ui/
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			// Short cache for HTML (may change frequently)
			w.Header().Set("Cache-Control", "public, max-age=60")
			http.ServeFile(w, r, "web/index.html")
			return
		}

		// Serve static assets (css, js) with caching
		if r.URL.Path == "/style.css" || r.URL.Path == "/app.js" {
			// Cache static assets for 1 hour (versioned via query params)
			w.Header().Set("Cache-Control", "public, max-age=3600")
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

			// Cache generated images for 24 hours (immutable content)
			w.Header().Set("Cache-Control", "public, max-age=86400, immutable")
			w.Header().Set("Content-Type", "image/png")
			if _, err := io.Copy(w, reader); err != nil {
				slog.Error("Failed to stream image content", "error", err)
			}
			return
		}

		// Serve version endpoint
		if r.URL.Path == "/version" {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Cache-Control", "no-cache")
			if err := json.NewEncoder(w).Encode(map[string]string{"version": version}); err != nil {
				slog.Error("Failed to write version response", "error", err)
			}
			return
		}

		// Serve user identity endpoint
		if r.URL.Path == "/me" {
			userEmail := r.Header.Get("x-goog-authenticated-user-email")
			if userEmail != "" {
				// Strip "accounts.google.com:" prefix
				userEmail = strings.TrimPrefix(userEmail, "accounts.google.com:")
			} else {
				userEmail = "local-user"
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Cache-Control", "private, no-cache")
			if err := json.NewEncoder(w).Encode(map[string]string{"email": userEmail}); err != nil {
				slog.Error("Failed to write user identity response", "error", err)
			}
			return
		}

		// Fallback to ADK API handler
		handler.ServeHTTP(w, r)
	})

	// Apply gzip compression middleware
	finalHandler := gzipMiddleware(routingHandler)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           finalHandler,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       30 * time.Second,  // Limit time to read entire request
		WriteTimeout:      5 * time.Minute,   // Match ADK handler timeout for long-running ops
		IdleTimeout:       120 * time.Second, // Keep-alive connection timeout
		MaxHeaderBytes:    1 << 20,           // 1MB max header size
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
