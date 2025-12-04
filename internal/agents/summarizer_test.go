package agents

import (
	"iter"
	"testing"

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

// mockAgent implements agent.Agent for testing
type mockAgent struct {
	agent.Agent
}

func (m *mockAgent) Run(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
	return func(yield func(*session.Event, error) bool) {
		// Mock success
		yield(nil, nil)
	}
}

func (m *mockAgent) Name() string {
	return "MockAgent"
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
