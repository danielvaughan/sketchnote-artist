# Linting Processes

This project uses several linting tools to ensure code quality, consistency, and security.

## Linting Tools

| Language/Type | Tool | Configuration |
| :--- | :--- | :--- |
| **Go** | [golangci-lint](https://golangci-lint.run/) | [.golangci.yml](file:///Users/danielvaughan/Development/git/sketchnote-artist/.golangci.yml) |
| **Dockerfile** | [hadolint](https://github.com/hadolint/hadolint) | [.hadolint.yaml](file:///Users/danielvaughan/Development/git/sketchnote-artist/.hadolint.yaml) |
| **Markdown** | [markdownlint](https://github.com/igorshubovych/markdownlint-cli) | [.markdownlint.json](file:///Users/danielvaughan/Development/git/sketchnote-artist/.markdownlint.json) |
| **Secrets** | [gitleaks](https://github.com/gitleaks/gitleaks) | Default (.git/hooks or pre-commit) |
| **Terraform** | [tflint](https://github.com/terraform-linters/tflint) | Via pre-commit |

## Local Execution

You can run each linter manually using the following commands:

### Go

```bash
golangci-lint run --timeout=5m
```

### Docker

```bash
hadolint Dockerfile
```

### Markdown

```bash
npx markdownlint-cli "**/*.md" --ignore node_modules --ignore .gemini --ignore private
```

### Secrets

```bash
gitleaks detect --source=. --verbose
```

## Pre-commit Integration

The project uses [pre-commit](https://pre-commit.com/) to run linters automatically before each commit.

### Installation

1. Install pre-commit: `brew install pre-commit` (or via pip).
2. Install the hooks:

   ```bash
   pre-commit install
   ```

### Manual Run

To run all linters against all files manually:

```bash
pre-commit run --all-files
```

## CI/CD Pipeline

Linting is enforced in the Google Cloud Build pipeline as defined in [cloudbuild.yaml](file:///Users/danielvaughan/Development/git/sketchnote-artist/cloudbuild.yaml). The build will fail if any linter reports an error.

The CI/CD pipeline runs:

1. `golangci-lint`
2. `markdownlint-cli`
3. `hadolint`
4. `gitleaks`
