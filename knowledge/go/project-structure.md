# Go Project Structure Best Practices

This document outlines standard practices for organizing Go projects, with a focus on where different types of tests should reside.

## Standard Directory Layout

While Go doesn't enforce a specific folder structure, the following layout is widely accepted in the community:

* **`/cmd`**: Main applications (entry points) for the project. For example, `cmd/sketchnote/main.go`. These should be minimal and call into internal packages.
* **`/internal`**: Private application and library code. Packages here cannot be imported by other projects. This is where most of your business logic should live.
* **`/pkg`**: Library code that is safe for use by external applications. (Use sparingly; prefer `internal` unless you explicitly want to provide a public API).
* **`/api`**: API protocol definitions (e.g., OpenAPI/Swagger specs, JSON schema files).
* **`/web`**: Web static assets, templates, and frontend code.
* **`/configs`**: Configuration file templates or default configs.
* **`/scripts`**: Scripts for build, install, analysis, etc.

## Test Placement

Testing is a first-class citizen in Go. Here is where the different types of tests should go:

### 1. Unit Tests

* **Location**: Alongside the code they test.
* **Naming**: `file_test.go` in the same package as `file.go`.
* **Best Practice**: Use the `pkg_test` naming convention for the test package (e.g., `package app_test`) to ensure you only test the public API of the package.

### 2. Integration Tests

* **Location**: Can be in the same directory as unit tests or in a dedicated `/test` folder.
* **Tagging**: Use Go build tags (e.g., `// +build integration`) at the top of the file to separate them from fast-running unit tests.
* **Execution**: `go test -tags=integration ./...`

### 3. End-to-End (E2E) Tests

* **Location**: In a dedicated **`/e2e`** or **`/tests/e2e`** directory at the project root.
* **Technology (Playwright)**: Since E2E tests often involve browser automation or high-level API flows that cross multiple boundaries, they should be kept outside the `/internal` or `/pkg` hierarchy.
* **Playwright Structure**:
  * `e2e/*.spec.ts`: Test specifications.
  * `playwright.config.ts`: Configuration file at the root.
  * `package.json`: Node.js dependencies at the root.

## Current Project Application

In `sketchnote-artist`, we follow this pattern:

* `cmd/`: Entry points for CLI and Server.
* `internal/`: Core logic (agents, tools, flows).
* `e2e/`: Playwright tests (API and Web UI).
* `knowledge/`: Documentation and project-specific guides.
