package agents

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/danielvaughan/sketchnote-artist/internal/observability"
	"github.com/danielvaughan/sketchnote-artist/internal/prompts"
	"github.com/danielvaughan/sketchnote-artist/internal/storage"
	"github.com/danielvaughan/sketchnote-artist/internal/tools"
)

const ArtistEmoji = "🎨"

// NewArtist creates the artist agent.
func NewArtist(ctx context.Context, apiKey string, store storage.Store) (agent.Agent, error) {
	// Initialize genai client for the tool
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	// Initialize the model for the Artist agent
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create artist model: %w", err)
	}

	imageTool, err := tools.NewImageGenerationTool(client, store, "sketchnotes")
	if err != nil {
		return nil, fmt.Errorf("failed to create image generation tool: %w", err)
	}

	innerAgent, err := llmagent.New(llmagent.Config{
		Name:        "Artist",
		Model:       model,
		Description: "Creates sketchnotes from visual briefs.",
		Instruction: prompts.ArtistInstruction,
		OutputKey:   "sketchnote",
		Tools:       []tool.Tool{imageTool},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create inner artist agent: %w", err)
	}

	return agent.New(agent.Config{
		Name:        "Artist",
		Description: "Creates sketchnotes from visual briefs.",
		Run: func(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
			// Check if the visual brief exists in the state
			_, err := ctx.Session().State().Get("visual_brief")
			if err != nil {
				slog.Warn("Artist skipping execution: visual_brief missing from state", "error", err)
				return func(yield func(*session.Event, error) bool) {}
			}

			observability.Report(ctx, fmt.Sprintf("%s The Artist is reading the visual brief...", ArtistEmoji))

			return func(yield func(*session.Event, error) bool) {
				for event, err := range innerAgent.Run(ctx) {
					if !yield(event, err) {
						return
					}
				}
			}
		},
	})
}
