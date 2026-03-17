// Package flows defines the workflow logic for the Sketchnote Artist.
package flows

import (
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
)

// NewSketchnoteFlow creates the sequential flow for the sketchnote artist.
// Errors from sub-agents are logged at their source (curator.go, artist.go)
// and propagated to the caller via the event iterator.
func NewSketchnoteFlow(curator agent.Agent, artist agent.Agent) (agent.Agent, error) {
	return sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "sketchnote-artist",
			Description: "A flow that summarizes a video and then creates a sketchnote.",
			SubAgents:   []agent.Agent{curator, artist},
		},
	})
}
