// Package tools contains the specific tools used by the agents (Image generation, YouTube, etc.).
package tools

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/danielvaughan/sketchnote-artist/internal/config"
	"github.com/danielvaughan/sketchnote-artist/internal/observability"
	"github.com/danielvaughan/sketchnote-artist/internal/storage"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

// NewImageGenerationTool creates a new tool for generating images.
func NewImageGenerationTool(client *genai.Client, store storage.Store, folder string) (tool.Tool, error) {
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
			// Sanitize filename to snake_case
			filename := filepath.Base(args.Filename)
			ext := filepath.Ext(filename)
			name := strings.TrimSuffix(filename, ext)

			// Convert to lowercase
			name = strings.ToLower(name)
			// Replace spaces and hyphens with underscores
			name = strings.ReplaceAll(name, " ", "_")
			name = strings.ReplaceAll(name, "-", "_")
			// Remove any other non-alphanumeric characters (simple clean)
			// For a robust implementation we might want regex, but standard lib strings map is safer/simpler for now
			// taking a simple approach: if not letter/digit/underscore, make it underscore
			var builder strings.Builder
			for _, r := range name {
				if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
					builder.WriteRune(r)
				} else {
					builder.WriteRune('_')
				}
			}
			name = builder.String()

			// Collapse multiple underscores
			for strings.Contains(name, "__") {
				name = strings.ReplaceAll(name, "__", "_")
			}
			// Trim underscores
			name = strings.Trim(name, "_")

			filename = name + ext
			observability.Report(ctx, fmt.Sprintf("\n%s The Artist is sketching...", "🎨"))

			// Log the image generation request with prompt details
			promptPreview := prompt
			if len(promptPreview) > 200 {
				promptPreview = promptPreview[:200] + "..."
			}
			slog.Info("Generating image",
				"filename", filename,
				"model", config.ImageGeneratorToolModel,
				"prompt_length", len(prompt),
				"prompt_preview", promptPreview)

			// Call Imagen 3 model
			resp, err := client.Models.GenerateContent(ctx, config.ImageGeneratorToolModel, genai.Text(prompt), nil)
			if err != nil {
				slog.Error("Image generation failed",
					"error", err,
					"model", config.ImageGeneratorToolModel,
					"prompt_length", len(prompt))
				return "", fmt.Errorf("generation failed: %w", err)
			}

			// Save the image bytes to a file
			for _, candidate := range resp.Candidates {
				for _, part := range candidate.Content.Parts {
					if part.InlineData != nil {
						// Ensure filename is unique if it already exists
						exists, err := store.Exists(ctx, folder, filename)
						if err == nil && exists {
							// If file exists, append timestamp before extension
							var ext string
							if strings.Contains(filename, ".") {
								// Simple extension check/split might be safer
								// For now assume .png as per previous code context or just simple append
								parts := strings.Split(filename, ".")
								if len(parts) > 1 {
									ext = "." + parts[len(parts)-1]
									name := strings.TrimSuffix(filename, ext)
									filename = fmt.Sprintf("%s_%d%s", name, time.Now().UnixNano(), ext)
								} else {
									filename = fmt.Sprintf("%s_%d", filename, time.Now().UnixNano())
								}
							}
						}

						if err := store.Save(ctx, folder, filename, part.InlineData.Data); err != nil {
							slog.Error("Failed to save image", "error", err)
							return "", err
						}

						observability.Report(ctx, fmt.Sprintf("\n%s The Artist has finished! View your sketchnote here: %s", "🎨", filename))
						return fmt.Sprintf("Image successfully saved to %s", filename), nil
					}
				}
			}

			// If we get here, no image was found. Construct a detailed error message.
			var errorDetails strings.Builder
			errorDetails.WriteString("No image data returned by model.")

			// Log comprehensive error details with structured logging
			logAttrs := []any{
				"model", config.ImageGeneratorToolModel,
				"prompt_length", len(prompt),
				"prompt", prompt,
				"filename", filename,
			}

			// Add PromptFeedback details if present
			if resp.PromptFeedback != nil {
				logAttrs = append(logAttrs,
					"prompt_feedback_block_reason", resp.PromptFeedback.BlockReason,
					"prompt_feedback_block_message", resp.PromptFeedback.BlockReasonMessage)

				if resp.PromptFeedback.BlockReason != "" && resp.PromptFeedback.BlockReason != genai.BlockedReasonUnspecified {
					errorDetails.WriteString(fmt.Sprintf(" Prompt blocked: Reason=%v", resp.PromptFeedback.BlockReason))
					if resp.PromptFeedback.BlockReasonMessage != "" {
						errorDetails.WriteString(fmt.Sprintf(", Message=%s", resp.PromptFeedback.BlockReasonMessage))
					}
					errorDetails.WriteString(".")
				}
			}

			// Add UsageMetadata if available
			if resp.UsageMetadata != nil {
				logAttrs = append(logAttrs,
					"prompt_token_count", resp.UsageMetadata.PromptTokenCount,
					"candidates_token_count", resp.UsageMetadata.CandidatesTokenCount,
					"total_token_count", resp.UsageMetadata.TotalTokenCount)
			}

			// Log detailed candidate information
			logAttrs = append(logAttrs, "candidate_count", len(resp.Candidates))

			for i, candidate := range resp.Candidates {
				candidatePrefix := fmt.Sprintf("candidate_%d", i)
				logAttrs = append(logAttrs,
					candidatePrefix+"_finish_reason", string(candidate.FinishReason),
					candidatePrefix+"_finish_message", candidate.FinishMessage)

				errorDetails.WriteString(fmt.Sprintf(" Candidate %d: FinishReason=%s", i, candidate.FinishReason))
				if candidate.FinishMessage != "" {
					errorDetails.WriteString(fmt.Sprintf(", Message=%s", candidate.FinishMessage))
				}

				// Add safety ratings
				if len(candidate.SafetyRatings) > 0 {
					errorDetails.WriteString(" SafetyRatings=[")
					for j, rating := range candidate.SafetyRatings {
						errorDetails.WriteString(fmt.Sprintf("%s:%s ", rating.Category, rating.Probability))
						logAttrs = append(logAttrs,
							fmt.Sprintf("%s_safety_%d_category", candidatePrefix, j), string(rating.Category),
							fmt.Sprintf("%s_safety_%d_probability", candidatePrefix, j), string(rating.Probability))
						if rating.ProbabilityScore > 0 {
							logAttrs = append(logAttrs,
								fmt.Sprintf("%s_safety_%d_score", candidatePrefix, j), rating.ProbabilityScore)
						}
						if rating.Severity != "" {
							logAttrs = append(logAttrs,
								fmt.Sprintf("%s_safety_%d_severity", candidatePrefix, j), string(rating.Severity))
						}
					}
					errorDetails.WriteString("]")
				}

				// Add grounding metadata if present
				if candidate.GroundingMetadata != nil {
					logAttrs = append(logAttrs,
						candidatePrefix+"_has_grounding_metadata", true)
				}

				errorDetails.WriteString(".")
			}

			// Log with all structured attributes
			slog.Warn(errorDetails.String(), logAttrs...)
			return errorDetails.String(), nil
		},
	)
}
