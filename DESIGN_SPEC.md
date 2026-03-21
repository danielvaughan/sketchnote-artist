# Sketchnote Artist — Design Spec

## Overview

Sketchnote Artist converts YouTube videos into hand-drawn visual summaries (sketchnotes)
using a two-agent sequential workflow powered by Google ADK Go.

## Agents

### Curator Agent

- **Model:** `gemini-3-flash-preview`
- **Input:** YouTube video URL (validated before LLM call)
- **Tools:** `summarize_youtube_video`, `save_to_file`
- **Output:** A structured "Visual Brief" stored in `state["visual_brief"]` via `OutputKey`
- **Constraint:** Must validate URL format before any LLM call. Invalid inputs return a
  human-readable error message without consuming tokens.

### Artist Agent

- **Model:** `gemini-3-flash-preview`
- **Input:** `state["visual_brief"]` injected into instruction via ADK `{visual_brief}` template
- **Tools:** `generate_image`
- **Output:** A sketchnote PNG saved to storage; final response includes the filename
- **Constraint:** Must not proceed if `visual_brief` is absent from state — return an error event.

## Tools

| Tool | Model | Purpose |
| --- | --- | --- |
| `summarize_youtube_video` | `gemini-3-flash-preview` | Streams video transcript/content from YouTube |
| `generate_image` | `gemini-3-pro-image-preview` | Calls Imagen to produce PNG |
| `save_to_file` | — | Persists text content to storage backend |

## Sequential Flow

```text
User URL → [Curator] → state["visual_brief"] → [Artist] → PNG filename
```

## ADK Template Substitution

- `{youtube_url}` in `CuratorInstruction`: set in session state by the outer Curator agent
  before running the inner llmagent. Non-optional (no `?`) — key must always be present.
- `{visual_brief}` in `ArtistInstruction`: written to session state by Curator via
  `OutputKey: "visual_brief"`. Non-optional.

## Success Criteria

1. Valid YouTube URL → sketchnote PNG produced and saved within 120 seconds
2. Invalid input → human-readable error returned, no LLM tokens consumed
3. Curator Visual Brief contains all required fields: Title, Speaker, Central Metaphor,
   Visual Hierarchy, Iconography, Mood
4. Artist image prompt references the brief title and at least 3 key points
5. Generated PNG filename is returned in the final agent response

## Constraints

- All model names defined in `internal/config/config.go` — never hardcoded inline
- URL validation in the outer agent wrapper, before any LLM call
- `{youtube_url}` and `{visual_brief}` resolved via ADK session state — no `strings.ReplaceAll`
- Storage backend: local disk in development, GCS in `cloud_run` mode
- Artifact service: in-memory in development, GCS in `cloud_run` mode (`GCS_BUCKET_ARTIFACTS`)
