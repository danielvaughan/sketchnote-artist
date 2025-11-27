package agents

import (
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"

	"github.com/danielvaughan/sketchnote-artist/internal/prompts"
)

// NewSummarizer creates the summarizer agent.
func NewSummarizer(model model.LLM, tools []tool.Tool) (agent.Agent, error) {
	return llmagent.New(llmagent.Config{
		Name:        "Summarizer",
		Model:       model,
		Description: "Extracts summaries from YouTube videos.",
		Instruction: prompts.SummarizerInstruction,
		OutputKey:   "visual_brief",
		Tools:       tools,
	})
}

// NewArtist creates the artist agent.
func NewArtist(model model.LLM, tools []tool.Tool) (agent.Agent, error) {
	return llmagent.New(llmagent.Config{
		Name:        "Artist",
		Model:       model,
		Description: "Creates sketchnotes from visual briefs.",
		Instruction: prompts.ArtistInstruction,
		OutputKey:   "sketchnote",
		Tools:       tools,
	})
}
