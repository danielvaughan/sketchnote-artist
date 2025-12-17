# Testing Processes

This document outlines the testing procedures for the Sketchnote Artist Agent.

## Unit Tests

To run unit tests (skipping slow integration tests):

```bash
go test ./...
```

To run all tests including integration tests (requires API key):

```bash
go test -tags=integration ./...
```

## End-to-End (E2E) Tests

Automated end-to-end tests are verified against the deployed `dev` environment using Playwright.

### Prerequisites

1. **Install Node.js dependencies:**

    ```bash
    npm install
    ```

2. **Install Playwright browsers:**

    ```bash
    npx playwright install --with-deps
    ```

### Configuration

Set the Service URL based on your testing environment:

* **For Local Testing:**

    ```bash
    export SERVICE_URL=http://localhost:8080
    ```

* **For Deployed Environment:**

    Retrieve the URL from Terraform outputs:

    ```bash
    export SERVICE_URL=$(cd terraform && terraform output -raw service_url)
    ```

### Running Tests

* **UI Test** (Simulates user interaction in the browser):

    ```bash
    npx playwright test e2e/webui.spec.ts
    ```

* **API Test** (Directly calls REST endpoints):

    ```bash
    npx playwright test e2e/api.spec.ts
    ```

* **Run All Tests:**

    ```bash
    npx playwright test
    ```

## Automation

### Pre-commit Hooks

The project uses `pre-commit` to ensure code quality. Unit tests are automatically run before every commit:

```yaml
- id: go-unit-tests
  entry: go test ./...
```

### CI/CD Pipeline

The Cloud Build pipeline (`cloudbuild.yaml`) automatically runs unit tests on every push to the `dev` branch:

1. **Unit Tests**: `go test ./...`

> [!NOTE]
> Integration tests are excluded from the CI/CD pipeline to ensure fast and reliable builds. They should be run manually during the verification phase.

#### Secret Management

Integration tests require the `GOOGLE_API_KEY`, which is stored in **Google Cloud Secret Manager**.

* **Secret ID**: `GOOGLE_API_KEY-dev` or `GOOGLE_API_KEY-prod` (environment-specific).
* **Access**: The Cloud Build service account has `roles/secretmanager.secretAccessor` permission, granted via Terraform in `terraform/secrets.tf`.
* **Usage**: The secret is passed to the build via a substitution in `cloudbuild.tf`:
  `_GOOGLE_API_KEY = "sm://projects/${var.project_id}/secrets/GOOGLE_API_KEY-${local.env}/versions/latest"`
