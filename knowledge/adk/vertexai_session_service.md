# Vertex AI Session Service in Go ADK

By default, the Go Agent Development Kit (ADK) uses an `InMemoryService`, which
is volatile. For production deployments on Google Cloud, particularly when using
Cloud Run, you should use the **Vertex AI Session Service** to ensure persistent
and scalable session management.

## Persistence with Vertex AI Agent Engine

The Vertex AI session service leverages the **Vertex AI Agent Engine** to store
and manage conversation history. This ensures that:

- Conversational state is preserved across restarts (ideal for serverless like Cloud Run).
- Multiple server instances can share the same session data.
- User history is persisted in a managed Google Cloud service.

## Implementation Details

### Required Imports

You will need the following packages:

```go
import (
    "google.golang.org/adk/session"
    "google.golang.org/adk/session/vertexai"
)
```

### Initialization

The session service is initialized using a context and your Google Cloud project
details. Internally, it uses the authenticated environment (ADC) to connect to
Vertex AI.

```go
ctx := context.Background()

// Initialize the Vertex AI Session Service
// It typically requires the Project ID and Location (e.g., "us-central1")
sessionService, err := vertexai.NewService(ctx, vertexai.Config{
    Project:  os.Getenv("GOOGLE_CLOUD_PROJECT"),
    Location: os.Getenv("GOOGLE_CLOUD_LOCATION"),
})
if err != nil {
    log.Fatalf("failed to create vertex ai session service: %v", err)
}
```

### Integration with Launcher

Once created, replace the `InMemoryService` in your `launcher.Config`:

```go
config := &launcher.Config{
    AgentLoader:     agent.NewSingleLoader(agentInstance),
    SessionService:  sessionService, // Replaces session.InMemoryService()
    ArtifactService: artifact.InMemoryService(),
    MemoryService:   memory.InMemoryService(),
}
```

## Setup Requirements

1. **Vertex AI API**: Ensure the Vertex AI API is enabled in your Project.
2. **IAM Permissions**: The service account running the application must have the `roles/aiplatform.user` role.
3. **Reasoning Engine**: In some ADK versions, you may need a Reasoning Engine ID.

> [!NOTE]
> When using `VertexAISessionService`, the ADK will automatically handle saving
> events to the cloud after each turn and resuming them when the same
> `sessionId` is provided.
