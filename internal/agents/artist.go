package agents

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/danielvaughan/sketchnote-artist/internal/prompts"
	"github.com/danielvaughan/sketchnote-artist/internal/tools"
)

// NewArtist creates the artist agent.
func NewArtist(ctx context.Context, apiKey string) (agent.Agent, error) {
	// Initialize genai client for the tool
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	// Initialize the Gemini 3.0 Pro Image model for art
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create artist model: %w", err)
	}

	imageTool, err := tools.NewImageGenerationTool(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create image generation tool: %w", err)
	}

	return llmagent.New(llmagent.Config{
		Name:        "Artist",
		Model:       model,
		Description: "Creates sketchnotes from visual briefs.",
		Instruction: prompts.ArtistInstruction,
		OutputKey:   "sketchnote",
		Tools:       []tool.Tool{imageTool},
	})
}
