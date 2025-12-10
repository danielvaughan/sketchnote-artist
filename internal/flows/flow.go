package flows

import (
	"iter"
	"log/slog"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/session"
)

// NewSketchnoteFlow creates the sequential flow for the sketchnote artist.
func NewSketchnoteFlow(curator agent.Agent, artist agent.Agent) (agent.Agent, error) {
	seqAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "sketchnote-artist",
			Description: "A flow that summarizes a video and then creates a sketchnote.",
			SubAgents:   []agent.Agent{curator, artist},
		},
	})
	if err != nil {
		return nil, err
	}

	return &loggingAgent{Agent: seqAgent}, nil
}

// loggingAgent wraps an agent to log errors.
type loggingAgent struct {
	agent.Agent
}

// Run executes the agent and logs any error.
func (l *loggingAgent) Run(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
	return func(yield func(*session.Event, error) bool) {
		next := l.Agent.Run(ctx)
		for event, err := range next {
			if err != nil {
				slog.Error("Flow execution failed", "error", err)
			}
			if !yield(event, err) {
				return
			}
		}
	}
}
