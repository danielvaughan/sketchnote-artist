package agents

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	adkmodel "google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/danielvaughan/sketchnote-artist/internal/observability"
	"github.com/danielvaughan/sketchnote-artist/internal/prompts"
	"github.com/danielvaughan/sketchnote-artist/internal/tools"
)

const CuratorEmoji = "🧐"

// NewCurator creates the curator agent.
func NewCurator(ctx context.Context, apiKey string) (agent.Agent, error) {
	// Initialize the Gemini model for the Curator agent
	model, err := gemini.NewModel(ctx, "gemini-3-pro-preview", &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create curator model: %w", err)
	}

	// Initialize the tools
	ytTool, err := tools.NewYouTubeSummarizer(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create YouTube summarizer tool: %w", err)
	}

	fileTool, err := tools.NewFileSaver()
	if err != nil {
		return nil, fmt.Errorf("failed to create file saver tool: %w", err)
	}

	return agent.New(agent.Config{
		Name:        "Curator",
		Description: "Curates video content to create a visual brief for the artist.",
		Run: func(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
			content := ctx.UserContent()
			var input string
			if content != nil {
				for _, part := range content.Parts {
					if part.Text != "" {
						input += part.Text
					}
				}
			}
			input = strings.TrimSpace(input)

			// Logging
			slog.Info("Agent Input", "agent", "Curator", "input", input)

			// Validation
			if !ValidateYouTubeURL(input) {
				slog.Warn("Invalid input rejected", "agent", "Curator", "input_snippet", input)
				return func(yield func(*session.Event, error) bool) {
					yield(&session.Event{
						LLMResponse: adkmodel.LLMResponse{
							Content: &genai.Content{
								Role: "model",
								Parts: []*genai.Part{
									{Text: "I can't process that. Please provide a valid YouTube URL (e.g., https://www.youtube.com/watch?v=...)."},
								},
							},
						},
					}, nil)
				}
			}

			// Report progress via context
			observability.Report(ctx, fmt.Sprintf("%s The Curator is analyzing the video to create a visual brief: %s", CuratorEmoji, input))

			// Inject URL into instruction
			instruction := strings.ReplaceAll(prompts.CuratorInstruction, "{YouTubeURL}", input)

			// Create dynamic agent
			innerAgent, err := llmagent.New(llmagent.Config{
				Name:        "Curator",
				Model:       model,
				Description: "Curates video content to create a visual brief for the artist.",
				Instruction: instruction,
				OutputKey:   "visual_brief",
				Tools:       []tool.Tool{ytTool, fileTool},
			})
			if err != nil {
				return func(yield func(*session.Event, error) bool) {
					yield(nil, fmt.Errorf("failed to create inner agent: %w", err))
				}
			}

			startTime := time.Now()
			return func(yield func(*session.Event, error) bool) {
				for event, err := range innerAgent.Run(ctx) {
					if !yield(event, err) {
						return
					}
				}
				// Report completion
				observability.Report(ctx, fmt.Sprintf("%s The Curator has completed the visual brief in %s", CuratorEmoji, time.Since(startTime).Round(time.Second)))
			}
		},
	})
}

// ValidateYouTubeURL checks if the input string is a valid YouTube URL.
func ValidateYouTubeURL(input string) bool {
	ytRegex := regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.be)/.+$`)
	return ytRegex.MatchString(input)
}
