package agents

import (
	"context"
	"strings"
	"testing"

	"github.com/danielvaughan/sketchnote-artist/internal/storage"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

// mockInvocationContext implements agent.InvocationContext for testing
type mockInvocationContext struct {
	agent.InvocationContext
	userContent *genai.Content
}

func (m *mockInvocationContext) UserContent() *genai.Content {
	return m.userContent
}

func (m *mockInvocationContext) Artifacts() agent.Artifacts {
	return nil
}

func (m *mockInvocationContext) Session() session.Session {
	return nil
}

func (m *mockInvocationContext) Memory() agent.Memory {
	return nil
}

func (m *mockInvocationContext) InvocationID() string {
	return "test-invocation"
}

func (m *mockInvocationContext) Branch() string {
	return "test-branch"
}

func (m *mockInvocationContext) RunConfig() *agent.RunConfig {
	return nil
}

func (m *mockInvocationContext) EndInvocation() {}

func (m *mockInvocationContext) Ended() bool {
	return false
}

func TestValidateYouTubeURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "Valid YouTube URL",
			input: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			want:  true,
		},
		{
			name:  "Valid Short YouTube URL",
			input: "https://youtu.be/dQw4w9WgXcQ",
			want:  true,
		},
		{
			name:  "Invalid Input",
			input: "not a url",
			want:  false,
		},
		{
			name:  "Invalid URL Domain",
			input: "https://www.google.com",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateYouTubeURL(tt.input); got != tt.want {
				t.Errorf("ValidateYouTubeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurator_Run_GracefulFailure(t *testing.T) {
	ctx := context.Background()
	// Dummy API key (must be non-empty)
	apiKey := "dummy_key"

	// Create agent
	cAgent, err := NewCurator(ctx, apiKey, &storage.DiskStore{})
	if err != nil {
		t.Fatalf("Failed to create Curator: %v", err)
	}

	// Mock invocation context with invalid input
	mockCtx := &mockInvocationContext{
		userContent: &genai.Content{
			Parts: []*genai.Part{
				{Text: "this is not a valid url"},
			},
		},
	}

	// Run agent
	iter := cAgent.Run(mockCtx)

	// Collect events
	var events []*session.Event
	for event, err := range iter {
		if err != nil {
			t.Errorf("Unexpected error from Run: %v", err)
		}
		if event != nil {
			events = append(events, event)
		}
	}

	// Verify we got the helpful message
	found := false
	expected := "I can't process that. Please provide a valid YouTube URL"
	for _, e := range events {
		if e.Content != nil && len(e.Content.Parts) > 0 {
			text := e.Content.Parts[0].Text
			if strings.Contains(text, expected) {
				found = true
				break
			}
		}
	}

	if !found {
		t.Errorf("Expected event with message containing %q, but not found", expected)
	}
}
func TestNormalizeYouTubeURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Standard Long URL",
			input: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			want:  "https://youtu.be/dQw4w9WgXcQ",
		},
		{
			name:  "Long URL with extra params",
			input: "https://www.youtube.com/watch?v=dQw4w9WgXcQ&feature=emb_logo",
			want:  "https://youtu.be/dQw4w9WgXcQ",
		},
		{
			name:  "Short URL",
			input: "https://youtu.be/dQw4w9WgXcQ",
			want:  "https://youtu.be/dQw4w9WgXcQ",
		},
		{
			name:  "Embed URL",
			input: "https://www.youtube.com/embed/dQw4w9WgXcQ",
			want:  "https://youtu.be/dQw4w9WgXcQ",
		},
		{
			name:  "Invalid URL",
			input: "https://www.google.com",
			want:  "https://www.google.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeYouTubeURL(tt.input); got != tt.want {
				t.Errorf("NormalizeYouTubeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractVideoID(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"Standard", "https://www.youtube.com/watch?v=dQw4w9WgXcQ", "dQw4w9WgXcQ"},
		{"Short", "https://youtu.be/dQw4w9WgXcQ", "dQw4w9WgXcQ"},
		{"Embed", "https://www.youtube.com/embed/dQw4w9WgXcQ", "dQw4w9WgXcQ"},
		{"Extra Params", "https://www.youtube.com/watch?v=dQw4w9WgXcQ&t=10s", "dQw4w9WgXcQ"},
		{"Invalid", "https://www.google.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractVideoID(tt.url); got != tt.want {
				t.Errorf("ExtractVideoID() = %v, want %v", got, tt.want)
			}
		})
	}
}
