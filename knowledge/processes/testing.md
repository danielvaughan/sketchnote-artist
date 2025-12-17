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
