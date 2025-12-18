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

	maxAttempts := 3
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if attempt > 1 {
			slog.Info("Retrying YouTube summarization", "attempt", attempt, "url", videoURL)
		}

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
		streamErr := false

		for resp, err := range iter {
			if err != nil {
				slog.Error("Error processing stream chunk", "error", err, "attempt", attempt)
				lastErr = err

				// If it's a retryable error and we have attempts left, break inner loop to retry
				if strings.Contains(err.Error(), "unexpected EOF") && attempt < maxAttempts {
					streamErr = true
					break
				}
				// Otherwise return the error
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

		if streamErr {
			// Wait before retrying
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
			continue
		}

		summary := sb.String()
		if summary != "" && summary != "No summary could be generated." {
			slog.Info("Generated summary", "summary_len", len(summary), "duration", time.Since(startTime), "attempts", attempt)
			return summary, nil
		}

		if attempt == maxAttempts {
			if summary == "" {
				slog.Warn("No summary could be generated after max attempts", "attempts", attempt)
				return "No summary could be generated.", nil
			}
			return summary, nil
		}
	}

	return "No summary could be generated.", lastErr
}
