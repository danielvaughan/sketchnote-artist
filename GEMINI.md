# Sketchnote Artist Agent

## Overview

The **Sketchnote Artist Agent** is a Go-based application built using the [Google GenAI Agent Development Kit (ADK)](https://github.com/googleapis/genai-agent-adk-go). It employs a sequential workflow of two specialized agents to convert YouTube video content into a visual sketchnote.

## Architecture

The system uses a **Sequential Agent** pattern (`sketchnote-artist`) composed of:

1. **Curator Agent (`Curator`)**:
    * **Model**: `gemini-3-flash-preview`
    * **Tool**: `summarize_youtube_video`
    * **Function**: Watches a YouTube video (via URL/URI), extracts key insights, and structures them into a "Visual Brief".
2. **Artist Agent (`Artist`)**:
    * **Model**: `gemini-3-flash-preview`
    * **Tool**: `generate_image`
    * **Function**: Takes the visual brief and uses the `gemini-3-pro-image-preview` (Imagen 3) model to generate a high-quality sketchnote image simulating hand-drawn styles.

## Prerequisites

* **Go**: Version 1.25.3 or later.
* **Google Cloud Project**: With Gemini API access enabled.
* **API Key**: A Google Cloud API Key with access to Gemini and Imagen models.

## Setup

1. Ensure a `.env` file exists in the project root:

    ```env
    GOOGLE_API_KEY=your_api_key_here
    ```

2. Install dependencies:

    ```bash
    go mod download
    ```

## Usage

You can run the agent in console mode:

```bash
go run ./cmd/sketchnote console
```

Or build the binary:

```bash
go build -o sketchnote ./cmd/sketchnote
./sketchnote console
```

### Example Interaction

The agent is designed to accept natural language commands.

> "Create a sketchnote for this video: <https://www.youtube.com/watch?v=>..."

*Note: The YouTube tool currently expects a valid File URI or supported URL format compatible with the Gemini `FileURI` parameter.*

## Key Files

* `cmd/sketchnote/main.go`: CLI entry point. Initializes logs, models, and runs the launcher.
* `cmd/server/main.go`: REST API entry point.
* `internal/app/agent.go`: High-level agent creation and composition.
* `internal/flows/flow.go`: Defines the sequential multi-agent flow.
* `internal/agents/`: Agent definitions for Curator and Artist.
* `internal/tools/`: Tool implementations (YouTube summarizer, Image generator, File saver).
* `internal/prompts/prompts.go`: System instructions for the agents.
* `internal/config/config.go`: Central model and environment configuration.

## Development

* **Framework**: `google.golang.org/adk`
* **GenAI Client**: `google.golang.org/genai`
* **Pattern**: Sequential Multi-Agent
* **Commits**: Always use [Git Commit Guidelines](knowledge/processes/commits.md) when making commits.
* **Pull Requests**: Always use [Pull Request Procedure](knowledge/processes/pull_requests.md) when creating PRs.
* **Testing**: Always use [Testing Processes](knowledge/processes/testing.md) when asked to run tests.
* **Local Running**: Always use [Local Development](knowledge/processes/local-running.md) when asked to run locally.

## Project Rules

* **Temporary Files**: All temporary files created by an agent (e.g., screenshots, intermediate task outputs) must be placed in the `/temp/` directory at the project root. This directory is excluded from version control.
* **Versioning**: Always consult [Versioning Guidelines](knowledge/processes/versioning.md) when making commits to ensure the `VERSION` file is correctly incremented.

## Knowledge Base

A collection of resources and documentation for the project.

### Project Resources

* [Build Log](knowledge/build_log.md): Tracks the build process, decisions, and architecture.

### Go Development

* [Go Project Structure](knowledge/go/project-structure.md): Best practices for organizing Go code and tests.
* [Go Best Practices](knowledge/go/go-best-practices.md): Guidelines for writing idiomatic and maintainable Go code.
* [Go Linting and Reviews](knowledge/go/go_linting_and_reviews.md): Standards for linting and code reviews.
* [Logging Best Practices](knowledge/go/logging_best_practices.md): Recommendations for structured logging.

### Agent Development Kit (ADK)

* [ADK Go Quickstart](knowledge/adk/adk_go_quickstart.md): Quickstart guide for the Google GenAI ADK in Go.
* [ADK Multi-Agent Patterns](knowledge/adk/adk_multi_agent_patterns.md): Common patterns for building multi-agent systems.
* [Context Reporting Pattern](knowledge/adk/context_reporting_pattern.md): Implementation details for context reporting.
* [HTTP SSE Streaming Pattern](knowledge/adk/http_sse_streaming_pattern.md): Guide on implementing Server-Sent Events for streaming agent responses.

### Testing

* [Testing Processes](knowledge/processes/testing.md): Detailed testing procedures for Unit and End-to-End tests. **Consult this document when performing any test operations.**
* [Terraform Deployment](knowledge/processes/terraform.md): Infrastructure management and deployment steps.
* [Linting Processes](knowledge/processes/linting.md): Code quality and static analysis tools.
* [Local Development](knowledge/processes/local-running.md): Instructions for running the project on your local machine.
* [Git Commit Guidelines](knowledge/processes/commits.md): Standardized commit message templates and best practices.
