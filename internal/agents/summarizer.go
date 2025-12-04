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

// NewSummarizer creates the summarizer agent.
func NewSummarizer(ctx context.Context, apiKey string) (agent.Agent, error) {
	// Initialize genai client for the tool
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	// Initialize the Gemini model for summarization
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create summarizer model: %w", err)
	}

	// Initialize the tools
	ytTool, err := tools.NewYouTubeSummarizer(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create YouTube summarizer tool: %w", err)
	}

	fileTool, err := tools.NewFileSaver(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create file saver tool: %w", err)
	}

	return llmagent.New(llmagent.Config{
		Name:        "Summarizer",
		Model:       model,
		Description: "Extracts summaries from YouTube videos.",
		Instruction: prompts.SummarizerInstruction,
		OutputKey:   "visual_brief",
		Tools:       []tool.Tool{ytTool, fileTool},
	})
}
