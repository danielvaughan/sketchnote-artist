package tools

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/danielvaughan/sketchnote-artist/internal/observability"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

// NewImageGenerationTool creates a new tool for generating images.
func NewImageGenerationTool(client *genai.Client, outputDir string) (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "generate_image",
			Description: "Generates an image based on a text prompt and saves it to disk. Returns the file path.",
		},
		func(ctx tool.Context, args struct {
			Prompt   string `json:"prompt" doc:"The detailed visual description of the image to generate."`
			Filename string `json:"filename" doc:"The desired filename for the generated image (e.g., visual_brief.png)."`
		}) (string, error) {
			prompt := args.Prompt
			filename := args.Filename
			observability.Report(ctx, fmt.Sprintf("\n%s The Artist is sketching...", "🎨"))
			slog.Info("Generating image", "filename", filename)

			// Ensure output directory exists
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				slog.Error("Failed to create output directory", "error", err)
				return "", fmt.Errorf("failed to create output directory: %w", err)
			}

			// Call Imagen 3 model
			resp, err := client.Models.GenerateContent(ctx, "gemini-3-pro-image-preview", genai.Text(prompt), nil)
			if err != nil {
				slog.Error("Image generation failed", "error", err)
				return "", fmt.Errorf("generation failed: %w", err)
			}

			// Helper to get full path and check existence
			getFullPath := func(fname string) string {
				return filepath.Join(outputDir, fname)
			}

			// Save the image bytes to a file
			for _, candidate := range resp.Candidates {
				for _, part := range candidate.Content.Parts {
					if part.InlineData != nil {
						fullPath := getFullPath(filename)
						// Ensure filename is unique if it already exists
						if _, err := os.Stat(fullPath); err == nil {
							// If file exists, append timestamp before extension
							ext := ".png"
							name := strings.TrimSuffix(filename, ext)
							filename = fmt.Sprintf("%s_%d%s", name, time.Now().UnixNano(), ext)
							fullPath = getFullPath(filename)
						}

						if err := os.WriteFile(fullPath, part.InlineData.Data, 0644); err != nil {
							slog.Error("Failed to save image", "error", err)
							return "", err
						}
						slog.Info("Image saved", "filename", fullPath)
						// Return relative path (filename only) or full path? Returning full relative path for clarity.
						// Wait, the agent might get confused if we change the return value format too much.
						// The web frontend expects /images/filename.
						// If we return just filename, the UI code works.
						// If we return visual-briefs/filename, existing patterns might break.
						// Let's stick to returning "Successfully saved to [filename]" where filename handles the user/agent expectation.
						// Actually, better to return just the filename so the agent knows what it saved as.
						observability.Report(ctx, fmt.Sprintf("\n%s The Artist has finished! View your sketchnote here: %s", "🎨", filename))
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
