package app

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"

	"github.com/danielvaughan/sketchnote-artist/internal/agents"
	"github.com/danielvaughan/sketchnote-artist/internal/flows"
	"github.com/danielvaughan/sketchnote-artist/internal/storage"
)

// Config holds configuration for the sketchnote agent.
type Config struct {
	APIKey string
	Store  storage.Store
}

// NewSketchnoteAgent creates the sequential sketchnote agent.
func NewSketchnoteAgent(ctx context.Context, cfg Config) (agent.Agent, error) {
	// Create the Curator Agent
	curatorAgent, err := agents.NewCurator(ctx, cfg.APIKey, cfg.Store)
	if err != nil {
		return nil, fmt.Errorf("failed to create curator agent: %w", err)
	}

	// Create the Artist Agent
	artistAgent, err := agents.NewArtist(ctx, cfg.APIKey, cfg.Store)
	if err != nil {
		return nil, fmt.Errorf("failed to create artist agent: %w", err)
	}

	// Create the Sequential Agent
	seqAgent, err := flows.NewSketchnoteFlow(curatorAgent, artistAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to create sequential agent: %w", err)
	}

	return seqAgent, nil
}
