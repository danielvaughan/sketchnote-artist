package tools

import (
	"fmt"
	"log/slog"

	"github.com/danielvaughan/sketchnote-artist/internal/storage"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// NewFileSaver creates a new tool for saving content to a file.
func NewFileSaver(store storage.Store, folder string) (tool.Tool, error) {
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

			slog.Info("Saving content to file", "folder", folder, "filename", filename)

			if err := store.Save(ctx, folder, filename, []byte(content)); err != nil {
				slog.Error("Failed to save file", "error", err)
				return "", fmt.Errorf("failed to save file: %w", err)
			}

			return fmt.Sprintf("Successfully saved content to %s", filename), nil
		},
	)
}
