# ADK Best Practices Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix all 10 ADK best-practice issues identified in the code review, from P0 correctness bugs through to P2 polish items.

**Architecture:** Each fix is independent and scoped to a small number of files. P0 fixes change agent/prompt logic; P1 fixes add eval scaffolding and persistent cloud services; P2 fixes tighten the HTTP server and test code.

**Tech Stack:** Go 1.25, google.golang.org/adk v0.6.0, google.golang.org/genai v1.50.0, artifact/gcsartifact (ADK built-in)

---

## File Map

| File | Change |
|------|--------|
| `internal/agents/curator_test.go` | Fix compilation: add `WithContext` + non-nil `mockSession`; annotate brittle mock |
| `internal/prompts/prompts.go` | Replace `{YouTubeURL}` string-replace with `{youtube_url}` ADK template; document both `{youtube_url}` and `{visual_brief}` substitution |
| `internal/agents/curator.go` | Pre-build `innerAgent` once; set `youtube_url` in session state before running it |
| `internal/agents/artist.go` | Return error event instead of silent empty iterator when `visual_brief` missing |
| `internal/flows/flow.go` | Remove `loggingAgent` wrapper; add flow-level error logging inline |
| `cmd/server/main.go` | GCS artifact service in cloud_run mode; `http.FileServer` for static assets; derive Content-Type from extension |
| `DESIGN_SPEC.md` | Create spec document |
| `evals/evalset.json` | Seed eval cases (documentation artefact — ADK Go v0.6.0 has no built-in eval runner) |
| `evals/test_config.json` | Eval configuration |

---

## Task 0: Fix test compilation (pre-requisite for all other tasks)

The existing tests fail to compile because `mockInvocationContext` is missing the `WithContext(context.Context) InvocationContext` method required by `agent.InvocationContext` (added in ADK v0.6.0). Additionally, `agent.Run` unconditionally calls `ctx.Session().ID()` for telemetry before entering user code, so `Session()` must return a non-nil stub.

**Files:**

- Modify: `internal/agents/curator_test.go`

- [ ] **Step 1: Verify the tests fail to compile today**

```bash
go test ./internal/agents/... 2>&1 | head -5
```

Expected: compile error mentioning `missing method WithContext`

- [ ] **Step 2: Add `WithContext` and a minimal `mockSession` stub**

In `internal/agents/curator_test.go`, add `WithContext` to `mockInvocationContext` and add a `mockSession` type. Replace the entire mock block:

```go
// mockSession is a minimal session stub satisfying session.Session.
// agent.Run calls ctx.Session().ID() unconditionally for telemetry before
// entering user code, so Session() must not return nil.
type mockSession struct{}

func (m *mockSession) ID() string                    { return "test-session" }
func (m *mockSession) AppName() string               { return "test-app" }
func (m *mockSession) UserID() string                { return "test-user" }
func (m *mockSession) State() session.State          { return nil }
func (m *mockSession) Events() session.Events        { return nil }
func (m *mockSession) LastUpdateTime() time.Time     { return time.Time{} }

// mockInvocationContext implements agent.InvocationContext for testing.
//
// NOTE: This manually implements the full agent.InvocationContext interface.
// This is unavoidable in Go without ADK providing test helpers. If the
// interface gains new methods in a future ADK version, add them here or
// the package will fail to compile.
type mockInvocationContext struct {
 context.Context
 userContent *genai.Content
}

func (m *mockInvocationContext) UserContent() *genai.Content { return m.userContent }
func (m *mockInvocationContext) Agent() agent.Agent          { return nil }
func (m *mockInvocationContext) Artifacts() agent.Artifacts  { return nil }
func (m *mockInvocationContext) Session() session.Session    { return &mockSession{} }
func (m *mockInvocationContext) Memory() agent.Memory        { return nil }
func (m *mockInvocationContext) InvocationID() string        { return "test-invocation" }
func (m *mockInvocationContext) Branch() string              { return "test-branch" }
func (m *mockInvocationContext) RunConfig() *agent.RunConfig { return nil }
func (m *mockInvocationContext) EndInvocation()              {}
func (m *mockInvocationContext) Ended() bool                 { return false }

// WithContext returns a new context with the given context.Context embedded.
// Required by agent.InvocationContext since ADK v0.6.0.
func (m *mockInvocationContext) WithContext(ctx context.Context) agent.InvocationContext {
 return &mockInvocationContext{Context: ctx, userContent: m.userContent}
}
```

Add `"time"` to imports if not already present.

- [ ] **Step 3: Run tests to verify compilation and existing tests pass**

```bash
go test ./internal/agents/...
```

Expected: all tests pass

- [ ] **Step 4: Commit**

```bash
git add internal/agents/curator_test.go
git commit -m "fix(test): add WithContext and mockSession to satisfy agent.InvocationContext v0.6.0"
```

---

## Task 1: Fix prompt templates (P0 — #3)

Add comments explaining ADK's `{key}` → session state substitution. Change `{YouTubeURL}` to `{youtube_url}` to match the state key set in Task 2.

**Important contract:** Non-optional `{key}` placeholders raise an error if the key is absent from state. The `?` suffix (e.g. `{key?}`) makes substitution optional. `{youtube_url}` is non-optional because the Curator always sets it before running the inner agent.

**Files:**

- Modify: `internal/prompts/prompts.go`

- [ ] **Step 1: Update `CuratorInstruction` placeholder and add template comments to both instructions**

In `internal/prompts/prompts.go`:

```go
// CuratorInstruction is the system prompt for the Curator agent.
//
// ADK template substitution: {youtube_url} is replaced at runtime with the
// value of session state key "youtube_url". The outer Curator agent sets this
// key before invoking this inner llmagent — see internal/agents/curator.go.
// Non-optional placeholders (no trailing ?) cause a hard error if the key is
// absent, so the Set call in the outer agent is required.
const CuratorInstruction = `You are an expert Content Strategist for visual note-taking.

Your goal is to create a "Visual Brief" for the video at the provided URL: {youtube_url}
...` // keep the rest of the prompt unchanged, only the first line and comment change
```

```go
// ArtistInstruction is the system prompt for the Artist agent.
//
// ADK template substitution: {visual_brief} is replaced at runtime with the
// value of session state key "visual_brief". The Curator agent writes its
// final output text to this key via OutputKey: "visual_brief" in its
// llmagent.Config — see internal/agents/curator.go.
const ArtistInstruction = `...` // unchanged body
```

- [ ] **Step 2: Build**

```bash
go build ./...
```

Expected: clean

- [ ] **Step 3: Commit**

```bash
git add internal/prompts/prompts.go
git commit -m "docs(prompts): document ADK template substitution; rename YouTubeURL placeholder to youtube_url"
```

---

## Task 2: Pre-build Curator's inner agent (P0 — #1)

Build `innerAgent` once in `NewCurator`. Set `youtube_url` in session state on the happy path so ADK substitutes `{youtube_url}` in the instruction at invocation time.

Note: `State()` returns `nil` in the unit test mock (Task 0). The `State().Set()` call runs only on the happy path (after URL validation passes), which is not exercised by the graceful-failure unit test. Integration tests and production always have a non-nil `State`.

**Files:**

- Modify: `internal/agents/curator.go`

- [ ] **Step 1: Refactor to pre-build innerAgent and set state**

Replace the `NewCurator` function with:

```go
func NewCurator(ctx context.Context, apiKey string, store storage.Store) (agent.Agent, error) {
 model, err := gemini.NewModel(ctx, config.CuratorModel, &genai.ClientConfig{
  APIKey: apiKey,
 })
 if err != nil {
  return nil, fmt.Errorf("failed to create curator model: %w", err)
 }

 ytTool, err := tools.NewYouTubeSummarizer(ctx, apiKey)
 if err != nil {
  return nil, fmt.Errorf("failed to create YouTube summarizer tool: %w", err)
 }

 fileTool, err := tools.NewFileSaver(store, "visual-briefs")
 if err != nil {
  return nil, fmt.Errorf("failed to create file saver tool: %w", err)
 }

 // Build the inner LLM agent once. {youtube_url} in CuratorInstruction is
 // resolved by ADK from session state at invocation time (non-optional, so
 // the key must be present — see state Set call below).
 innerAgent, err := llmagent.New(llmagent.Config{
  Name:        "Curator",
  Model:       model,
  Description: "Curates video content to create a visual brief for the artist.",
  Instruction: prompts.CuratorInstruction,
  OutputKey:   "visual_brief",
  Tools:       []tool.Tool{ytTool, fileTool},
 })
 if err != nil {
  return nil, fmt.Errorf("failed to create inner curator agent: %w", err)
 }

 return agent.New(agent.Config{
  Name:        "Curator",
  Description: "Curates video content to create a visual brief for the artist.",
  Run: func(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
   content := ctx.UserContent()
   var input string
   if content != nil {
    for _, part := range content.Parts {
     if part.Text != "" {
      input += part.Text
     }
    }
   }
   input = strings.TrimSpace(input)
   input = NormalizeYouTubeURL(input)

   slog.Info("Agent Input", "agent", "Curator", "input", input)

   if !ValidateYouTubeURL(input) {
    slog.Warn("Invalid input rejected", "agent", "Curator", "input_snippet", input)
    msg := "I can't process that. Please provide a valid YouTube URL (e.g., https://www.youtube.com/watch?v=...)."
    observability.Report(ctx, msg)
    return func(yield func(*session.Event, error) bool) {
     yield(&session.Event{
      LLMResponse: adkmodel.LLMResponse{
       Content: &genai.Content{
        Role:  "model",
        Parts: []*genai.Part{{Text: msg}},
       },
      },
     }, nil)
    }
   }

   observability.Report(ctx, fmt.Sprintf("%s The Curator is analyzing the video to create a visual brief: %s", CuratorEmoji, input))

   // Set youtube_url in session state so the {youtube_url} ADK template
   // placeholder in CuratorInstruction resolves correctly at invocation time.
   if err := ctx.Session().State().Set("youtube_url", input); err != nil {
    return func(yield func(*session.Event, error) bool) {
     yield(nil, fmt.Errorf("failed to set youtube_url in state: %w", err))
    }
   }

   startTime := time.Now()
   return func(yield func(*session.Event, error) bool) {
    for event, err := range innerAgent.Run(ctx) {
     if !yield(event, err) {
      return
     }
    }
    observability.Report(ctx, fmt.Sprintf("%s The Curator has completed the visual brief in %s", CuratorEmoji, time.Since(startTime).Round(time.Second)))
   }
  },
 })
}
```

- [ ] **Step 2: Build**

```bash
go build ./...
```

Expected: clean

- [ ] **Step 3: Run tests**

```bash
go test ./internal/agents/...
```

Expected: all pass (graceful-failure test exercises the validation path only; `State().Set()` is not reached)

- [ ] **Step 4: Commit**

```bash
git add internal/agents/curator.go
git commit -m "fix(curator): pre-build inner llmagent once; inject URL via session state"
```

---

## Task 3: Artist returns error event when visual_brief missing (P0 — #2)

**Files:**

- Modify: `internal/agents/artist.go`

- [ ] **Step 1: Replace silent no-op with error event**

In `internal/agents/artist.go`, replace:

```go
_, err := ctx.Session().State().Get("visual_brief")
if err != nil {
    slog.Warn("Artist skipping execution: visual_brief missing from state", "error", err)
    return func(_ func(*session.Event, error) bool) {}
}
```

With:

```go
_, err := ctx.Session().State().Get("visual_brief")
if err != nil {
    msg := "Artist cannot proceed: visual_brief is missing from session state."
    slog.Error("Artist aborting: visual_brief missing", "error", err)
    return func(yield func(*session.Event, error) bool) {
        yield(&session.Event{
            LLMResponse: adkmodel.LLMResponse{
                Content: &genai.Content{
                    Role:  "model",
                    Parts: []*genai.Part{{Text: msg}},
                },
            },
        }, nil)
    }
}
```

Add import `adkmodel "google.golang.org/adk/model"` to `artist.go`.

- [ ] **Step 2: Build and test**

```bash
go build ./... && go test ./internal/agents/...
```

Expected: clean

- [ ] **Step 3: Commit**

```bash
git add internal/agents/artist.go
git commit -m "fix(artist): yield error event when visual_brief missing from state"
```

---

## Task 4: Remove loggingAgent wrapper (P2 — #8)

The `loggingAgent` wrapper logs errors at the flow boundary but slog is already used at every call site where errors originate. Removing it eliminates the indirection. To preserve flow-level visibility, the sequential agent's errors are logged inside the `NewSketchnoteFlow` caller (the wrapper is removed; errors propagate naturally to the server's SSE handler which already handles them).

**Files:**

- Modify: `internal/flows/flow.go`

- [ ] **Step 1: Replace flow.go with the simplified version**

```go
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
```

- [ ] **Step 2: Build and test**

```bash
go build ./... && go test ./...
```

Expected: clean

- [ ] **Step 3: Commit**

```bash
git add internal/flows/flow.go
git commit -m "refactor(flows): remove loggingAgent wrapper; errors logged at source"
```

---

## Task 5: Fix HTTP server — static assets and Content-Type (P2 — #7, #9)

**Files:**

- Modify: `cmd/server/main.go`

- [ ] **Step 1: Replace individual static-asset path matches with `http.FileServer`**

Replace the existing `style.css`/`app.js` block with a general file server. Note: `.png` is intentionally excluded from static exts — generated images are served from the storage backend via `/images/`, not from `web/`.

```go
// Serve static web assets from web/ directory (css, js, fonts, icons).
// Generated images are served separately via the /images/ handler below.
ext := filepath.Ext(r.URL.Path)
staticExts := map[string]bool{".css": true, ".js": true, ".ico": true, ".svg": true, ".woff2": true}
if staticExts[ext] {
    http.FileServer(http.Dir("web")).ServeHTTP(w, r)
    return
}
```

Add `"path/filepath"` to imports if not already present.

- [ ] **Step 2: Derive Content-Type from filename extension for image serving**

Replace:

```go
w.Header().Set("Content-Type", "image/png")
```

With:

```go
contentType := mime.TypeByExtension(filepath.Ext(filename))
if contentType == "" {
    contentType = "image/png" // fallback for unknown or extension-less filenames
}
w.Header().Set("Content-Type", contentType)
```

Add `"mime"` to imports.

- [ ] **Step 3: Build**

```bash
go build ./cmd/server/...
```

Expected: clean

- [ ] **Step 4: Smoke-test**

```bash
pkill -f "go run cmd/server/main.go" 2>/dev/null; sleep 1
go run cmd/server/main.go > /tmp/server.log 2>&1 &
sleep 4
curl -sI http://localhost:8080/style.css | grep -i content-type
curl -sI http://localhost:8080/app.js | grep -i content-type
pkill -f "go run cmd/server/main.go"
```

Expected: `text/css` and `application/javascript` (or `text/javascript`)

- [ ] **Step 5: Commit**

```bash
git add cmd/server/main.go
git commit -m "fix(server): use FileServer for static assets; derive image Content-Type from extension"
```

---

## Task 6: Switch artifact service to GCS in cloud_run mode (P1 — #5)

Note: ADK v0.6.0 has no persistent memory service — `memory.InMemoryService()` is the only option and is acceptable since this project does not use ADK memory.

**Files:**

- Modify: `cmd/server/main.go`

- [ ] **Step 1: Add GCS artifact service for cloud_run deployments**

Replace:

```go
config := &launcher.Config{
    AgentLoader:     agent.NewSingleLoader(agentInstance),
    SessionService:  sessionService,
    ArtifactService: artifact.InMemoryService(),
    MemoryService:   memory.InMemoryService(),
}
```

With:

```go
var artifactService artifact.Service
if os.Getenv("DEPLOYMENT_MODE") == "cloud_run" {
    artifactsBucket := os.Getenv("GCS_BUCKET_ARTIFACTS")
    if artifactsBucket == "" {
        slog.Error("GCS_BUCKET_ARTIFACTS must be set in cloud_run mode")
        os.Exit(1)
    }
    slog.Info("Initializing GCS Artifact Service", "bucket", artifactsBucket)
    gcsSvc, err := gcsartifact.NewService(ctx, artifactsBucket)
    if err != nil {
        slog.Error("Failed to initialize GCS artifact service", "error", err)
        os.Exit(1)
    }
    artifactService = gcsSvc
} else {
    artifactService = artifact.InMemoryService()
}

config := &launcher.Config{
    AgentLoader:     agent.NewSingleLoader(agentInstance),
    SessionService:  sessionService,
    ArtifactService: artifactService,
    // ADK v0.6.0 has no persistent memory service; InMemoryService is acceptable
    // because this project does not use ADK memory state.
    MemoryService: memory.InMemoryService(),
}
```

Add import:

```go
"google.golang.org/adk/artifact/gcsartifact"
```

- [ ] **Step 2: Build**

```bash
go build ./cmd/server/...
```

Expected: clean

- [ ] **Step 3: Commit**

```bash
git add cmd/server/main.go
git commit -m "fix(server): use GCS artifact service in cloud_run mode; fail fast if bucket unset"
```

---

## Task 7: Create DESIGN_SPEC.md (P1 — #6)

**Files:**

- Create: `DESIGN_SPEC.md`

- [ ] **Step 1: Write the spec**

````markdown
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
|------|-------|---------|
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

````

- [ ] **Step 2: Commit**

```bash
git add DESIGN_SPEC.md
git commit -m "docs: add DESIGN_SPEC.md as ADK primary source of truth"
```

---

## Task 8: Create evaluation seed cases (P1 — #4)

**Note:** ADK Go v0.6.0 has no built-in eval runner. These files are documentation artefacts in the ADK Python evalset format, capturing expected agent behaviour for future use when a Go eval runner is available, or for use with the ADK Python runner against a deployed endpoint.

**Files:**

- Create: `evals/evalset.json`
- Create: `evals/test_config.json`

- [ ] **Step 1: Create eval directory and evalset**

```bash
mkdir -p evals
```

`evals/evalset.json`:

```json
[
  {
    "name": "invalid_input_graceful_failure",
    "description": "Non-URL input should return a helpful error without calling any tools",
    "conversation": [
      {
        "role": "user",
        "content": "hello world"
      }
    ],
    "expected_tool_use": [],
    "expected_intermediate_agents": [],
    "reference": "I can't process that. Please provide a valid YouTube URL"
  },
  {
    "name": "invalid_domain_graceful_failure",
    "description": "Non-YouTube URL should be rejected before any LLM tool call",
    "conversation": [
      {
        "role": "user",
        "content": "https://www.google.com/search?q=test"
      }
    ],
    "expected_tool_use": [],
    "expected_intermediate_agents": [],
    "reference": "I can't process that. Please provide a valid YouTube URL"
  }
]
```

`evals/test_config.json`:

```json
{
  "criteria": {
    "tool_trajectory_avg_score": 1.0,
    "response_match_score": 0.8
  }
}
```

- [ ] **Step 2: Validate JSON**

```bash
python3 -c "import json; json.load(open('evals/evalset.json')); json.load(open('evals/test_config.json')); print('OK')"
```

Expected: `OK`

- [ ] **Step 3: Commit**

```bash
git add evals/
git commit -m "test(evals): add ADK evaluation seed cases (graceful failure paths)"
```

---

## Task 9: Final build and test verification

- [ ] **Step 1: Full build**

```bash
go build ./...
```

Expected: clean

- [ ] **Step 2: Full unit test run**

```bash
go test ./...
```

Expected: all pass

- [ ] **Step 3: Checklist**

All issues addressed:

- ✅ P0 #1: Curator inner agent pre-built once (Task 2)
- ✅ P0 #2: Artist returns error event when visual_brief missing (Task 3)
- ✅ P0 #3: Template substitution documented in prompts (Task 1)
- ✅ P1 #4: Eval seed cases created (Task 8)
- ✅ P1 #5: GCS artifact service in cloud_run mode (Task 6)
- ✅ P1 #6: DESIGN_SPEC.md created (Task 7)
- ✅ P2 #7: Static assets via FileServer (Task 5)
- ✅ P2 #8: loggingAgent wrapper removed (Task 4)
- ✅ P2 #9: Content-Type derived from extension (Task 5)
- ✅ P2 #10: Brittle mock annotated + WithContext added (Task 0)
