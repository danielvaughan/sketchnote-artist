package main

import (
	"context"
	"log/slog"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/joho/godotenv"

	"github.com/danielvaughan/sketchnote-artist/internal/agents"
	"github.com/danielvaughan/sketchnote-artist/internal/flows"
	"github.com/danielvaughan/sketchnote-artist/internal/tools"
)

func main() {
	ctx := context.Background()

	// Initialize structured logging to file
	logFile, err := os.OpenFile("sketchnote-artist.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file", "error", err)
		os.Exit(1)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewTextHandler(logFile, nil))
	slog.SetDefault(logger)

	// Load .env file
	// Try loading from current directory first, then fallback to root
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

	// Initialize the Gemini model for summarization
	summarizerModel, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		slog.Error("Failed to create summarizer model", "error", err)
		os.Exit(1)
	}

	// Initialize the Gemini 3.0 Pro Image model for art
	artistModel, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		slog.Error("Failed to create artist model", "error", err)
		os.Exit(1)
	}

	// Initialize genai client for the tool
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		slog.Error("Failed to create genai client", "error", err)
		os.Exit(1)
	}

	// Initialize the tools
	ytTool, err := tools.NewYouTubeSummarizer(client)
	if err != nil {
		slog.Error("Failed to create YouTube summarizer tool", "error", err)
		os.Exit(1)
	}

	imageTool, err := tools.NewImageGenerationTool(client)
	if err != nil {
		slog.Error("Failed to create image generation tool", "error", err)
		os.Exit(1)
	}

	fileTool, err := tools.NewFileSaver(client)
	if err != nil {
		slog.Error("Failed to create file saver tool", "error", err)
		os.Exit(1)
	}

	// Create the Summarizer Agent
	summarizerAgent, err := agents.NewSummarizer(summarizerModel, []tool.Tool{ytTool, fileTool})
	if err != nil {
		slog.Error("Failed to create summarizer agent", "error", err)
		os.Exit(1)
	}

	// Create the Artist Agent
	artistAgent, err := agents.NewArtist(artistModel, []tool.Tool{imageTool})
	if err != nil {
		slog.Error("Failed to create artist agent", "error", err)
		os.Exit(1)
	}

	// Create the Sequential Agent
	seqAgent, err := flows.NewSketchnoteFlow(summarizerAgent, artistAgent)
	if err != nil {
		slog.Error("Failed to create sequential agent", "error", err)
		os.Exit(1)
	}

	// Configure the launcher
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(seqAgent),
	}

	// Run the agent using the full launcher
	l := full.NewLauncher()
	if err := l.Execute(ctx, config, os.Args[1:]); err != nil {
		slog.Error("Run failed", "error", err, "syntax", l.CommandLineSyntax())
		os.Exit(1)
	}
}
