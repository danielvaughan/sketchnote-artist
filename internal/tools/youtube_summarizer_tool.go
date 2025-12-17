package tools

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"net/http"
	"time"

	"github.com/danielvaughan/sketchnote-artist/internal/config"
	"github.com/danielvaughan/sketchnote-artist/internal/prompts"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

// NewYouTubeSummarizer creates a new tool for summarizing YouTube videos.
func NewYouTubeSummarizer(ctx context.Context, apiKey string) (tool.Tool, error) {
	// Initialize genai client for the tool with extended timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Minute,
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:     apiKey,
		HTTPClient: httpClient,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client for youtube summarizer: %w", err)
	}

	return functiontool.New(
		functiontool.Config{
			Name:        "summarize_youtube_video",
			Description: "Summarizes a YouTube video given its URL. Use this tool when the user asks for a summary, overview, or details about a YouTube video.",
		},
		func(ctx tool.Context, args struct {
			URL string `json:"url" doc:"The URL of the YouTube video to summarize."`
		}) (string, error) {
			return SummarizeVideo(ctx, client, args.URL)
		},
	)
}

// SummarizeVideo performs the actual summarization logic.
// It is exported to facilitate isolation testing.
func SummarizeVideo(ctx context.Context, client *genai.Client, videoURL string) (string, error) {
	startTime := time.Now()
	slog.Info("Calling tool: summarize_youtube_video", "url", videoURL)

	// Call Gemini with the Video URI using Streaming to avoid timeouts
	iter := client.Models.GenerateContentStream(ctx, config.YouTubeSummarizerToolModel, []*genai.Content{
		{
			Parts: []*genai.Part{
				{
					FileData: &genai.FileData{
						MIMEType: "video/*",
						FileURI:  videoURL,
					},
				},
				{
					Text: prompts.YouTubeSummarizerInstruction,
				},
			},
		},
	}, &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingLevel: genai.ThinkingLevelLow,
		},
	})

	var sb strings.Builder
	for resp, err := range iter {
		if err != nil {
			slog.Error("Error processing stream chunk", "error", err)
			if strings.Contains(err.Error(), "unexpected EOF") {
				return "", fmt.Errorf("gemini server timed out while summarizing the video (unexpected EOF): %w", err)
			}
			return "", fmt.Errorf("error during generation stream: %w", err)
		}

		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			for _, part := range resp.Candidates[0].Content.Parts {
				if part.Text != "" {
					sb.WriteString(part.Text)
				}
			}
		}
	}

	summary := sb.String()
	if summary != "" {
		slog.Info("Generated summary", "summary_len", len(summary), "duration", time.Since(startTime))
		return summary, nil
	}

	slog.Warn("No summary could be generated", "summary_len", 0)
	return "No summary could be generated.", nil
}
