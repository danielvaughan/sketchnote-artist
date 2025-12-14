research# Project Enhancements

1. [ ] Support structured responses from prompts
    - **Context**: In Python (Pydantic), one defines a class and the library handles schema generation and validation.
    - **Go Implementation Strategy**:
        - **Library**: `google.golang.org/genai` supports `ResponseSchema` and `ResponseMIMEType` in `GenerationConfig`.
        - **Pattern**:
            1.  Define a Go `struct` with `json` tags representing the desired output.
            2.  Construct a corresponding `genai.Schema` object (manually or via reflection helpers if available).
            3.  Set `ResponseMIMEType` to `application/json` (or `text/x.enum`).
            4.  Pass this configuration to the model/generation call.
            5.  Unmarshal the raw JSON string response into the Go struct using `json.Unmarshal`.
    - **Action**: Research if a helper library (like `invopop/jsonschema`) can auto-generate `genai.Schema` from structs to reduce boilerplate.
2. [ ] Persist sessions
    - **Context**: Currently using `session.InMemoryService()`, which is lost on container restart.
    - **Recommendation**: Use **Google Cloud Firestore**.
    - **Reasoning**:
        - **Serverless**: Matches Cloud Run's operational model.
        - **Scalability**: Handles high concurrency and scales to zero.
        - **Persistence**: Durable storage separate from the compute instance.
    - **Implementation Strategy**:
        1.  Create a new package `internal/persistence/firestore`.
        2.  Implement a struct `FirestoreSessionService` that satisfies the `session.Service` interface (from `google.golang.org/adk/session`).
            - *Note*: Need to inspect the `session.Service` interface definition (likely `Get`, `Save`, `List`, `Delete`).
        3.  Update `cmd/server/main.go` to initialize this service when `DEPLOYMENT_MODE=cloud_run`.
        4.  Use the `cloud.google.com/go/firestore` client library.
    - **Action**: Verify the exact signature of `session.Service` during implementation.
3. [ ] Store visual briefs and sketchnotes to cloud storage bucket when deployed on Cloud Run
    - **Context**: Currently, visual brief and sketchnote files are saved to the local file system.
    - **Goal**: When deployed to Cloud Run, save these assets to Google Cloud Storage (GCS) buckets instead.
    - **Implementation Details**:
        - **Environment Detection**: Enable this behavior via an environment variable (e.g., `DEPLOYMENT_MODE=cloud_run`) set during container deployment. Default to file system for other environments.
        - **Infrastructure (Terraform)**:
            - Create two GCS buckets: one for visual briefs, one for sketchnotes.
            - Ensure the application service account has appropriate read/write permissions for these buckets.
        - **UI Integration**: The application UI should serve sketchnotes directly from the sketchnotes GCS bucket.
4. [ ] Enable multi-environment deployment (Production vs. Dev)
    - **Context**: Currently, a single Cloud Build pipeline deploys the `main` branch to Cloud Run (Production).
    - **Goal**: Enable deployment from non-main branches to a separate "Dev" environment.
    - **Requirements**:
        - **Branching Strategy**: `main` deploys to Production. Other branches deploy to Dev.
        - **Resource Separation**:
            - Use separate Cloud Storage buckets for Dev and Prod.
            - Ensure logging and other resources are separated (e.g., via distinct service names or labels).
        - **Implementation**:
            - Update `cloudbuild.yaml` or triggers to handle conditional deployments.
            - Update Terraform to provision environment-specific resources.