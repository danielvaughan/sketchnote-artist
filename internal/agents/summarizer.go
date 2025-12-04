package agents

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"regexp"
	"strings"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
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

	a, err := llmagent.New(llmagent.Config{
		Name:        "Summarizer",
		Model:       model,
		Description: "Extracts summaries from YouTube videos.",
		Instruction: prompts.SummarizerInstruction,
		OutputKey:   "visual_brief",
		Tools:       []tool.Tool{ytTool, fileTool},
	})
	if err != nil {
		return nil, err
	}

	// Chain the agents: Logging -> Validation -> Summarizer
	validation := &validationAgent{Agent: a}
	logging := &loggingAgent{Agent: validation}

	return logging, nil
}

// ValidateYouTubeURL checks if the input string is a valid YouTube URL.
func ValidateYouTubeURL(input string) bool {
	ytRegex := regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.be)/.+$`)
	return ytRegex.MatchString(input)
}

type loggingAgent struct {
	agent.Agent
}

func (l *loggingAgent) Run(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
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
	slog.Info("Agent Input", "agent", l.Name(), "input", input)
	return l.Agent.Run(ctx)
}

type validationAgent struct {
	agent.Agent
}

func (v *validationAgent) Run(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
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

	if !ValidateYouTubeURL(input) {
		slog.Warn("Invalid input rejected", "agent", v.Name(), "input_snippet", input)
		return func(yield func(*session.Event, error) bool) {
			yield(nil, fmt.Errorf("invalid input: please provide a valid YouTube URL"))
		}
	}
	return v.Agent.Run(ctx)
}
