package tools

import (
	"fmt"
	"log/slog"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

// NewYouTubeSummarizer creates a new tool for summarizing YouTube videos.
func NewYouTubeSummarizer(client *genai.Client) (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "summarize_youtube_video",
			Description: "Summarizes a YouTube video given its URL. Use this tool when the user asks for a summary, overview, or details about a YouTube video.",
		},
		func(ctx tool.Context, args struct {
			URL string `json:"url" doc:"The URL of the YouTube video to summarize."`
		}) (string, error) {
			videoURL := args.URL
			fmt.Printf("\n🤖 A robot is binging the video for you: %s\n", videoURL)
			slog.Info("Agent is watching the YouTube video", "url", videoURL)

			// Call Gemini with the Video URI
			resp, err := client.Models.GenerateContent(ctx, "gemini-3-pro-preview", []*genai.Content{
				{
					Parts: []*genai.Part{
						{
							FileData: &genai.FileData{
								MIMEType: "video/*",
								FileURI:  videoURL,
							},
						},
						{
							Text: "Please provide a comprehensive summary of this video, highlighting the key points, main arguments, and any important conclusions.",
						},
					},
				},
			}, nil)

			if err != nil {
				slog.Error("Error processing video", "error", err)
				return "", fmt.Errorf("failed to process video: %w", err)
			}

			// Return the Summary
			if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
				summary := resp.Candidates[0].Content.Parts[0].Text
				slog.Info("Generated summary", "summary", summary)
				return summary, nil
			}

			slog.Warn("No summary could be generated")
			return "No summary could be generated.", nil
		},
	)
}
