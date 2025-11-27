package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize the Gemini model for summarization
	summarizerModel, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create summarizer model: %v", err)
	}

	// Initialize the Gemini 3.0 Pro Image model for art
	artistModel, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create artist model: %v", err)
	}

	// Initialize genai client for the tool
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create genai client: %v", err)
	}

	// Initialize the tool
	ytTool, err := NewYouTubeSummarizer(client)
	if err != nil {
		log.Fatalf("Failed to create YouTube summarizer tool: %v", err)
	}

	imageTool, err := NewImageGenerationTool(client)
	if err != nil {
		log.Fatalf("Failed to create image generation tool: %v", err)
	}

	// Create the Summarizer Agent
	summarizerAgent, err := llmagent.New(llmagent.Config{
		Name:        "Summarizer",
		Model:       summarizerModel,
		Description: "Extracts summaries from YouTube videos.",
		Instruction: SummarizerInstruction,
		OutputKey:   "visual_brief",
		Tools: []tool.Tool{
			ytTool,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create summarizer agent: %v", err)
	}

	// Create the Artist Agent
	artistAgent, err := llmagent.New(llmagent.Config{
		Name:        "Artist",
		Model:       artistModel,
		Description: "Creates sketchnotes from visual breifs.",
		Instruction: ArtistInstruction,
		OutputKey:   "sketchnote",
		Tools: []tool.Tool{
			imageTool,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create artist agent: %v", err)
	}

	// Create the Sequential Agent
	seqAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "SketchnoteFlow",
			Description: "A flow that summarizes a video and then creates a sketchnote.",
			SubAgents:   []agent.Agent{summarizerAgent, artistAgent},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create sequential agent: %v", err)
	}

	// Configure the launcher
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(seqAgent),
	}

	// Run the agent using the full launcher
	l := full.NewLauncher()
	if err := l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
