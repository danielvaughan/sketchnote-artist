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

	"github.com/danielvaughan/sketchnote-artist/internal/config"
	"github.com/danielvaughan/sketchnote-artist/internal/observability"
	"github.com/danielvaughan/sketchnote-artist/internal/prompts"
	"github.com/danielvaughan/sketchnote-artist/internal/storage"
	"github.com/danielvaughan/sketchnote-artist/internal/tools"
)

// CuratorEmoji is the emoji used for curator log messages.
const CuratorEmoji = "🧐"

// NewCurator creates the curator agent.
func NewCurator(ctx context.Context, apiKey string, store storage.Store) (agent.Agent, error) {
	// Initialize the Gemini model for the Curator agent
	model, err := gemini.NewModel(ctx, config.CuratorModel, &genai.ClientConfig{
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

	fileTool, err := tools.NewFileSaver(store, "visual-briefs")
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
			input = NormalizeYouTubeURL(input)

			// Logging
			slog.Info("Agent Input", "agent", "Curator", "input", input)

			// Validation
			if !ValidateYouTubeURL(input) {
				slog.Warn("Invalid input rejected", "agent", "Curator", "input_snippet", input)
				msg := "I can't process that. Please provide a valid YouTube URL (e.g., https://www.youtube.com/watch?v=...)."
				observability.Report(ctx, msg)
				return func(yield func(*session.Event, error) bool) {
					yield(&session.Event{
						LLMResponse: adkmodel.LLMResponse{
							Content: &genai.Content{
								Role: "model",
								Parts: []*genai.Part{
									{Text: msg},
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

// NormalizeYouTubeURL extracts the video ID and returns a shortened youtu.be URL.
// If no ID is found, it returns the original input.
func NormalizeYouTubeURL(input string) string {
	id := ExtractVideoID(input)
	if id == "" {
		return input
	}
	return fmt.Sprintf("https://youtu.be/%s", id)
}

// ExtractVideoID returns the 11-character YouTube video ID from various URL formats.
func ExtractVideoID(url string) string {
	regExp := regexp.MustCompile(`(?i)(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})`)
	match := regExp.FindStringSubmatch(url)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
