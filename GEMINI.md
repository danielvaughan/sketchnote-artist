# Sketchnote Artist Agent

## Overview

The **Sketchnote Artist Agent** is a Go-based application built using the [Google GenAI Agent Development Kit (ADK)](https://github.com/googleapis/genai-agent-adk-go). It employs a sequential workflow of two specialized agents to convert YouTube video content into a visual sketchnote.

## Architecture

The system uses a **Sequential Agent** pattern (`sketchnote-artist`) composed of:

1. **Curator Agent (`Curator`)**:
    * **Model**: `gemini-3-pro-preview` (Configured with Low Thinking)
    * **Tool**: `summarize_youtube_video`
    * **Function**: Watches a YouTube video (via URL/URI), extracts key insights, and structures them into a "Visual Brief".
2. **Artist Agent (`Artist`)**:
    * **Model**: `gemini-2.5-flash`
    * **Tool**: `generate_image`
    * **Function**: Takes the visual brief and uses the `gemini-3-pro-image-preview` model to generate a high-quality sketchnote image simulating hand-drawn styles.

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

You can run the agent directly using `go run`:

```bash
go run . [prompt]
```

Or build the binary:

```bash
go build -o sketchnote-artist
./sketchnote-artist [prompt]
```

### Example Interaction

The agent is designed to accept natural language commands.

> "Create a sketchnote for this video: <https://www.youtube.com/watch?v=>..."

*Note: The YouTube tool currently expects a valid File URI or supported URL format compatible with the Gemini `FileURI` parameter.*

## Key Files

* `agent.go`: Main entry point. Configures the agents, tools, and the sequential flow.
* `prompts.go`: Contains the system instructions for the `Curator` and `Artist` agents.
* `youtube_summarizer.go`: Implements the tool to send video content to the Gemini model for summarization.
* `image_tool.go`: Implements the tool to generate images using Imagen 3 and save them locally.

## Development

* **Framework**: `google.golang.org/adk`
* **GenAI Client**: `google.golang.org/genai`
* **Pattern**: Sequential Multi-Agent
