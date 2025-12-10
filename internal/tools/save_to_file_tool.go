package tools

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// NewFileSaver creates a new tool for saving content to a file.
func NewFileSaver(outputDir string) (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "save_to_file",
			Description: "Saves the provided text content to a file with the given filename.",
		},
		func(ctx tool.Context, args struct {
			Filename string `json:"filename" doc:"The name of the file to save (e.g., visual_brief.txt)."`
			Content  string `json:"content" doc:"The text content to write to the file."`
		}) (string, error) {
			filename := args.Filename
			content := args.Content

			// Ensure output directory exists
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				slog.Error("Failed to create output directory", "error", err)
				return "", fmt.Errorf("failed to create output directory: %w", err)
			}

			// Prepend output directory
			fullPath := filepath.Join(outputDir, filename)

			slog.Info("Saving content to file", "filename", fullPath)

			if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
				slog.Error("Failed to save file", "error", err)
				return "", fmt.Errorf("failed to save file: %w", err)
			}

			return fmt.Sprintf("Successfully saved content to %s", filename), nil
		},
	)
}
