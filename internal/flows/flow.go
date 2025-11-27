package flows

import (
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
)

// NewSketchnoteFlow creates the sequential flow for the sketchnote artist.
func NewSketchnoteFlow(summarizer agent.Agent, artist agent.Agent) (agent.Agent, error) {
	return sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "SketchnoteFlow",
			Description: "A flow that summarizes a video and then creates a sketchnote.",
			SubAgents:   []agent.Agent{summarizer, artist},
		},
	})
}
