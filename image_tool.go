package main

import (
	"fmt"
	"os"
	"time"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

// NewImageGenerationTool creates a new tool for generating images.
func NewImageGenerationTool(client *genai.Client) (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "generate_image",
			Description: "Generates an image based on a text prompt and saves it to disk. Returns the file path.",
		},
		func(ctx tool.Context, args struct {
			Prompt string `json:"prompt" doc:"The detailed visual description of the image to generate."`
		}) (string, error) {
			prompt := args.Prompt
			fmt.Printf(" [Tool] Generating image for: %q...\n", prompt)

			// Call Imagen 3 model
			resp, err := client.Models.GenerateContent(ctx, "gemini-3-pro-image-preview", genai.Text(prompt), nil)
			if err != nil {
				return "", fmt.Errorf("generation failed: %w", err)
			}

			// Save the image bytes to a file
			for _, candidate := range resp.Candidates {
				for _, part := range candidate.Content.Parts {
					if part.InlineData != nil {
						filename := fmt.Sprintf("generated_result_%d.png", time.Now().UnixNano())
						if err := os.WriteFile(filename, part.InlineData.Data, 0644); err != nil {
							return "", err
						}
						return fmt.Sprintf("Image successfully saved to %s", filename), nil
					}
				}
			}
			return "No image data returned by model.", nil
		},
	)
}
