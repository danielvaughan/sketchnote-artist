# Local Development & Execution

This guide details how to set up and run the Sketchnote Artist application locally on your machine.

## Prerequisites

* [Go](https://go.dev/dl/) (version 1.25.3 or later)
* A Google Cloud Project with the **Gemini API** enabled.
* A valid **Google Cloud API Key**.

## Installation & Setup

1. **Clone the repository:**

    ```bash
    git clone https://github.com/danielvaughan/sketchnote-artist.git
    cd sketchnote-artist
    ```

2. **Configure Environment Variables:**
    Create a `.env` file in the root directory:

    ```bash
    touch .env
    ```

    Add your API key to the file:

    ```env
    GOOGLE_API_KEY=your_google_api_key_here
    ```

3. **Install Dependencies:**

    ```bash
    go mod download
    ```

## Usage Modes

### CLI Mode (Console)

Run the agent directly using `go run` pointing to the main package:

```bash
go run ./cmd/sketchnote console
```

Enter just the video URL at the User prompt:

```bash
User -> https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

The agent will process the video, report progress, and save the resulting image (e.g., `generated_result_<timestamp>.png`) in the current directory.

### REST API Mode

The application includes a REST API server for programmatic interaction.

#### Starting the Server

```bash
go run cmd/server/main.go
```

The server listens on port `8080` by default.

#### Example API Interaction

1. **List Available Apps**:

    ```bash
    curl http://localhost:8080/list-apps
    ```

2. **Create a Session**:

    ```bash
    curl -X POST http://localhost:8080/apps/sketchnote-artist/users/test-user/sessions
    ```

3. **Run the Agent**:

    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{
      "appName": "sketchnote-artist",
      "userId": "test-user",
      "sessionId": "<session-id>",
      "newMessage": {
        "role": "user",
        "parts": [
          { "text": "https://www.youtube.com/watch?v=dQw4w9WgXcQ" }
        ]
      }
    }' http://localhost:8080/run
    ```
