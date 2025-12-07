//go:build integration

package tools

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func TestYouTubeSummarizer_Integration(t *testing.T) {
	// 1. Load .env
	if err := godotenv.Load("../../.env"); err != nil {
		t.Log("Warning: .env file not found. Assuming GOOGLE_API_KEY is set in environment.")
	}

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: GOOGLE_API_KEY not set")
	}

	// 2. Initialize Client with strict timeout control
	ctx := context.Background()

	// Create custom HTTP client with long timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Minute,
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:     apiKey,
		HTTPClient: httpClient,
	})
	if err != nil {
		t.Fatalf("Failed to create genai client: %v", err)
	}

	// 3. (Optional) Create Tool to verify it can be created
	_, err = NewYouTubeSummarizer(ctx, apiKey)
	if err != nil {
		t.Fatalf("Failed to create YouTube summarizer tool: %v", err)
	}

	// Define test cases
	tests := []struct {
		name      string
		inputURL  string
		timeout   time.Duration // 0 means no timeout (or default)
		wantError bool
		wantTitle string
	}{
		{
			name:      "Valid Short Video",
			inputURL:  "https://www.youtube.com/watch?v=jNQXAC9IVRw",
			wantError: false,
			wantTitle: "Me at the zoo",
		},
		/*
			{
				name:      "Client Timeout",
				inputURL:  "https://www.youtube.com/watch?v=jNQXAC9IVRw",
				timeout:   1 * time.Millisecond,
				wantError: true,
			},
		*/
		{
			// This video causes a server-side timeout (Unexpected EOF)
			name:      "Problematic Video",
			inputURL:  "https://www.youtube.com/watch?v=7Dtu2bilcFs",
			wantError: false,
			wantTitle: "Vibe Coding",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var runCtx context.Context
			var cancel context.CancelFunc
			if tt.timeout > 0 {
				runCtx, cancel = context.WithTimeout(ctx, tt.timeout)
			} else {
				runCtx, cancel = context.WithCancel(ctx)
			}
			defer cancel()

			// 4. Run Tool Function Directly
			result, err := SummarizeVideo(runCtx, client, tt.inputURL)

			if tt.wantError {
				// We expect an error or a specific "failed" message depending on implementation
				// The current implementation returns a string "No summary could be generated." on some failures,
				// or an error on API failures.
				if err == nil {
					if result == "No summary could be generated." {
						return // Success (it failed as expected)
					}
					// If it didn't fail, that might be unexpected for an invalid URL
					t.Logf("Expected error for invalid URL, but got nil error. Result: %v", result)
				} else {
					t.Logf("Got expected error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				if result == "" {
					t.Error("Expected non-empty summary")
				}
				if result == "No summary could be generated." {
					t.Error("Tool returned failure message for valid video")
				}

				if tt.wantTitle != "" {
					if !strings.Contains(result, tt.wantTitle) {
						t.Errorf("Summary does not contain expected title: %q", tt.wantTitle)
					} else {
						t.Logf("Confirmed title %q in summary", tt.wantTitle)
					}
				}

				t.Logf("Summary for %s:\n%s", tt.name, result)
			}
		})
	}
}
