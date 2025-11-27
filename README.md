# Sketchnote Artist Agent

The **Sketchnote Artist Agent** is an intelligent CLI application that turns YouTube videos into beautiful, hand-drawn style visual summaries (sketchnotes). 

Built with Go and the [Google GenAI Agent Development Kit (ADK)](https://github.com/googleapis/genai-agent-adk-go), it demonstrates the power of sequential multi-agent workflows.

## 🚀 How It Works

The application employs a chain of two specialized AI agents:

1.  **The Summarizer Agent**:
    *   **Role**: Content Strategist.
    *   **Task**: Watches the YouTube video, analyzes the content, and synthesizes a structured "Visual Brief" containing the core thesis, main takeaways, and memorable quotes.
    *   **Model**: Gemini 2.5 Flash.

2.  **The Artist Agent**:
    *   **Role**: Master Sketchnote Artist.
    *   **Task**: Interprets the Visual Brief and orchestrates the generation of a high-quality image that mimics alcohol markers and ink on paper.
    *   **Model**: Gemini 3.0 Pro Image (Imagen 3).

## 🛠️ Prerequisites

*   [Go](https://go.dev/dl/) (version 1.25.3 or later)
*   A Google Cloud Project with the **Gemini API** enabled.
*   A valid **Google Cloud API Key**.

## 📦 Installation & Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/danielvaughan/sketchnote-artist.git
    cd sketchnote-artist
    ```

2.  **Configure Environment Variables:**
    Create a `.env` file in the root directory:
    ```bash
    touch .env
    ```
    Add your API key to the file:
    ```env
    GOOGLE_API_KEY=your_google_api_key_here
    ```

3.  **Install Dependencies:**
    ```bash
    go mod download
    ```

## 🎨 Usage

Run the agent directly using `go run` and provide a prompt with a YouTube video URL.

```bash
go run . "Create a sketchnote for this video: https://www.youtube.com/watch?v=dQw4w9WgXcQ"
```

The agent will:
1.  Process the video.
2.  Print the progress of the agents.
3.  Save the resulting image as `generated_result_<timestamp>.png` in the current directory.

## 🏗️ Architecture

*   **`agent.go`**: The main entry point. Initializes the Gemini models, tools, and constructs the `SequentialAgent` workflow.
*   **`youtube_tool.go`**: A custom tool that sends video content to the Gemini model for analysis.
*   **`image_tool.go`**: A custom tool that interfaces with the Imagen 3 model to generate and save images locally.
*   **`prompts.go`**: Contains the system instructions that define the personas for the Summarizer and Artist agents.

## 📄 License

[MIT License](LICENSE)
