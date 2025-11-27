package main

import (
	"fmt"

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
			fmt.Printf("📺 Agent is watching YouTube video: %s\n", videoURL)

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
				return "", fmt.Errorf("failed to process video: %w", err)
			}

			// Return the Summary
			if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
				return resp.Candidates[0].Content.Parts[0].Text, nil
			}

			return "No summary could be generated.", nil
		},
	)
}
