// Package config defines the model constants used throughout the application.
package config

// Model constants for Gemini models used in the application.
// These are defined by their usage usage to ensure visibility into what each component uses.
const (
	// CuratorModel is the model used by the Curator agent.
	CuratorModel = "gemini-3-flash-preview"

	// ArtistModel is the model used by the Artist agent.
	ArtistModel = "gemini-3-flash-preview"

	// YouTubeSummarizerToolModel is the model used by the YouTube summarizer tool.
	YouTubeSummarizerToolModel = "gemini-3-flash-preview"

	// ImageGeneratorToolModel is the model used by the Image Generation tool.
	ImageGeneratorToolModel = "gemini-3-pro-image-preview"
)
