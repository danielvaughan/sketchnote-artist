package tools

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
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
			fmt.Printf("\n🎨  The artist is sketching...\n")
			slog.Info("Generating image", "prompt", prompt)

			// Call Imagen 3 model
			resp, err := client.Models.GenerateContent(ctx, "gemini-3-pro-image-preview", genai.Text(prompt), nil)
			if err != nil {
				slog.Error("Image generation failed", "error", err)
				return "", fmt.Errorf("generation failed: %w", err)
			}

			// Save the image bytes to a file
			for _, candidate := range resp.Candidates {
				for _, part := range candidate.Content.Parts {
					if part.InlineData != nil {
						// Extract title for filename
						title := extractTitle(prompt)
						filename := fmt.Sprintf("%s.png", title)

						// Ensure filename is unique if it already exists
						if _, err := os.Stat(filename); err == nil {
							filename = fmt.Sprintf("%s_%d.png", title, time.Now().UnixNano())
						}

						if err := os.WriteFile(filename, part.InlineData.Data, 0644); err != nil {
							slog.Error("Failed to save image", "error", err)
							return "", err
						}
						slog.Info("Image saved", "filename", filename)
						fmt.Printf("\n🎨  The artist has finished! View your sketchnote here: %s\n", filename)
						return fmt.Sprintf("Image successfully saved to %s", filename), nil
					}
				}
			}

			// If we get here, no image was found. Construct a detailed error message.
			var errorDetails strings.Builder
			errorDetails.WriteString("No image data returned by model.")

			if resp.PromptFeedback != nil && resp.PromptFeedback.BlockReason != "" && resp.PromptFeedback.BlockReason != genai.BlockedReasonUnspecified {
				errorDetails.WriteString(fmt.Sprintf(" Prompt blocked: Reason=%v", resp.PromptFeedback.BlockReason))
				if resp.PromptFeedback.BlockReasonMessage != "" {
					errorDetails.WriteString(fmt.Sprintf(", Message=%s", resp.PromptFeedback.BlockReasonMessage))
				}
				errorDetails.WriteString(".")
			}

			for i, candidate := range resp.Candidates {
				errorDetails.WriteString(fmt.Sprintf(" Candidate %d: FinishReason=%s", i, candidate.FinishReason))
				if candidate.FinishMessage != "" {
					errorDetails.WriteString(fmt.Sprintf(", Message=%s", candidate.FinishMessage))
				}
				if len(candidate.SafetyRatings) > 0 {
					errorDetails.WriteString(" SafetyRatings=[")
					for _, rating := range candidate.SafetyRatings {
						errorDetails.WriteString(fmt.Sprintf("%s:%s ", rating.Category, rating.Probability))
					}
					errorDetails.WriteString("]")
				}
				errorDetails.WriteString(".")
			}

			slog.Warn(errorDetails.String())
			return errorDetails.String(), nil
		},
	)
}

func extractTitle(prompt string) string {
	// regex to find "Title: 'Some Title'" or "Title Text: Some Title"
	// We look for "Title" optionally followed by " Text", then a colon, optional whitespace and quotes
	// We capture until the next quote or newline
	re := regexp.MustCompile(`(?i)Title(?:\s+Text)?:\s*['"]?([^\n]+?)['"]?\s*(?:\n|$)`)
	match := re.FindStringSubmatch(prompt)

	if len(match) > 1 {
		title := match[1]
		// Sanitize title
		// Replace spaces with underscores
		title = strings.ReplaceAll(title, " ", "_")
		// Remove non-alphanumeric characters (except underscores and hyphens)
		reg := regexp.MustCompile(`[^a-zA-Z0-9_\-]`)
		title = reg.ReplaceAllString(title, "")
		// Truncate if too long
		if len(title) > 100 {
			title = title[:100]
		}
		return title
	}

	// Default fallback
	return fmt.Sprintf("generated_result_%d", time.Now().UnixNano())
}
