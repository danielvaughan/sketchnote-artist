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

	// The 'time' import is not duplicated in the provided code.
	// If there was a duplicate, it would be removed here.
	"time"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/artifact"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/memory"
	"google.golang.org/adk/server/adkrest"
	"google.golang.org/adk/session"
	"google.golang.org/genai"

	"github.com/joho/godotenv"

	"github.com/danielvaughan/sketchnote-artist/internal/app"
	"github.com/danielvaughan/sketchnote-artist/internal/storage"
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

	// Initialize services
	sessionSvc := session.InMemoryService()

	// Configure the launcher with in-memory services
	config := &launcher.Config{
		AgentLoader:     agent.NewSingleLoader(agentInstance),
		SessionService:  sessionSvc,
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

		// NEW: Custom Streaming Handler (Server-Sent Events)
		// This bypasses the default ADK REST handler to provide real-time updates
		// and prevent "socket hang up" on long-running requests.
		if r.URL.Path == "/stream-run" && r.Method == "POST" {
			// 1. Parse Request
			var req struct {
				AppName    string `json:"appName"`
				UserID     string `json:"userId"`
				SessionID  string `json:"sessionId"`
				NewMessage struct {
					Role  string `json:"role"`
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"newMessage"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
				return
			}
			defer func() {
				if err := r.Body.Close(); err != nil {
					slog.Warn("Failed to close request body", "error", err)
				}
			}()

			if req.SessionID == "" {
				http.Error(w, "sessionId is required", http.StatusBadRequest)
				return
			}
			if len(req.NewMessage.Parts) == 0 {
				http.Error(w, "message text is required", http.StatusBadRequest)
				return
			}

			// 2. Prepare Context
			// Retrieve the session to ensure state (like visual_brief) is accessible
			resp, err := sessionSvc.Get(r.Context(), &session.GetRequest{
				AppName:   req.AppName,
				UserID:    req.UserID,
				SessionID: req.SessionID,
			})
			if err != nil {
				slog.Error("Failed to functionality load session", "session_id", req.SessionID, "error", err)
				http.Error(w, "Session not found", http.StatusNotFound)
				return
			}
			sess := resp.Session

			// We need to implement agent.InvocationContext to run the agent.
			// Since we don't have access to the internal adkrest context logic, we'll create a minimal implementation.
			invCtx := &SimpleInvocationContext{
				Context:   r.Context(),
				sessionID: req.SessionID,
				session:   sess,
				userContent: &genai.Content{
					Parts: []*genai.Part{{Text: req.NewMessage.Parts[0].Text}},
				},
			}

			// 3. Prepare Response for SSE
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			flusher, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "Streaming not supported", http.StatusInternalServerError)
				return
			}

			// 4. Run Agent and Stream Events
			slog.Info("Starting streaming run", "session_id", req.SessionID)

			// Send initial heartbeat
			if _, err := fmt.Fprintf(w, ": heartbeat\n\n"); err != nil {
				slog.Error("Failed to write heartbeat", "error", err)
				return
			}
			flusher.Flush()

			for event, err := range agentInstance.Run(invCtx) {
				if err != nil {
					slog.Error("Error during agent run", "error", err)
					// Send error event
					data, _ := json.Marshal(map[string]string{"error": err.Error()})
					if _, err := fmt.Fprintf(w, "event: error\ndata: %s\n\n", data); err != nil {
						slog.Error("Failed to write error event", "error", err)
						return
					}
					flusher.Flush()
					return
				}

				// Serialize event to JSON
				data, err := json.Marshal(event)
				if err != nil {
					slog.Warn("Failed to marshal event", "error", err)
					continue
				}

				if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
					slog.Error("Failed to write data event", "error", err)
					return
				}
				flusher.Flush()
			}

			// Send done event
			if _, err := fmt.Fprintf(w, "event: done\ndata: {}\n\n"); err != nil {
				slog.Error("Failed to write done event", "error", err)
				return
			}
			flusher.Flush()
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

// SimpleInvocationContext implements agent.InvocationContext for streaming
type SimpleInvocationContext struct {
	context.Context
	sessionID   string
	session     session.Session
	userContent *genai.Content
}

func (s *SimpleInvocationContext) Session() session.Session    { return s.session }
func (s *SimpleInvocationContext) RunConfig() *agent.RunConfig { return nil }
func (s *SimpleInvocationContext) InvocationID() string        { return "stream-" + s.sessionID }
func (s *SimpleInvocationContext) Memory() agent.Memory        { return nil }
func (s *SimpleInvocationContext) Artifacts() agent.Artifacts  { return nil }
func (s *SimpleInvocationContext) UserContent() *genai.Content { return s.userContent }
func (s *SimpleInvocationContext) EndInvocation()              {}
func (s *SimpleInvocationContext) Ended() bool                 { return false }
func (s *SimpleInvocationContext) Agent() agent.Agent          { return nil }
func (s *SimpleInvocationContext) Branch() string              { return "" }
