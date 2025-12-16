# Go Linting and Code Review Best Practices

High-quality Go code relies on consistent style, rigorous static analysis (linting), and thorough code reviews. This document outlines the best practices for setting up linting and conducting effective code reviews for Go projects.

## 1. Linting with `golangci-lint`

The standard tool for linting in the Go ecosystem is `golangci-lint`. It is a fast, parallel linter runner that bundles dozens of useful linters.

### Installation

```bash
# Binary (Project Standard: v2.7.2)
# Note: Ensure this version matches the one defined in cloudbuild.yaml
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.7.2

# Homebrew (macOS)
brew install golangci-lint
```

### Configuration

Create a `.golangci.yml` file in your project root. This ensures every team member and CI pipeline runs the exact same checks.

**Project Configuration (.golangci.yml):**

```yaml
version: "2"

run:
  timeout: 5m
  # Helps with performance in large repos
  skip-dirs:
    - vendor

linters:
  enable:
    - errcheck      # Check for unchecked errors
    - gosec         # Security checker
    - govet         # Official Go vet tool
    - ineffassign   # Detects unused assignments
    - staticcheck   # Domain-specific checks (very powerful)
    - unused        # Checks for unused constants/vars/functions
    # Recommended additions for 2025:
    - bodyclose     # Checks that HTTP bodies are closed
    - noctx         # Checks that http requests have context
    - revive        # Drop-in replacement for golint
    - gocritic      # Provides diagnostics that check for bugs, performance and style issues
    - misspell      # Finds commonly misspelled English words

formatters:
  enable:
    - gofmt
    - goimports
```

### Best Practices (2025 Research)

1.  **Version Consistency**: Always ensure your local linter version matches your CI version (`v2.7.2` in this project). Mismatches lead to "works on my machine" but fails in CI.
2.  **Pre-commit Hooks**: Use `pre-commit` to prevent bad code from entering the repo.
    ```yaml
    - repo: https://github.com/golangci/golangci-lint
      rev: v2.7.2
      hooks:
        - id: golangci-lint
    ```
3.  **Gradual Adoption**: If adding linters to a legacy project, use `--new-from-rev=HEAD~1` to only lint new changes, preventing an avalanche of errors.
4.  **Performance**: Enable `cache` in your CI/CD (e.g., using `actions/cache` or Kaniko) to speed up linting runs.

### Running Lint Checks

```bash
# Run all enabled linters
golangci-lint run

# Fix auto-fixable issues (formatting, imports)
golangci-lint run --fix
```

---

## 2. Code Review Checklist for Go

When reviewing Go code, look beyond logic errors. Ensure the code adheres to "idiomatic Go" (The Go Way).

### General
- [ ] **Formatting**: Is the code formatted with `gofmt`? (CI should enforce this).
- [ ] **Imports**: Are imports sorted? (Use `goimports`).
- [ ] **Naming**:
    - Variables: Short names for short scopes (`i`, `r`), descriptive names for long scopes.
    - CamelCase: No underscores (`user_id` -> `userID`).
    - Exported: Capitalized (`User`), unexported lowercase (`user`).

### Error Handling
- [ ] **Don't Panic**: Use error returns instead of `panic()` for normal error conditions.
- [ ] **Wrap Errors**: Use `fmt.Errorf("...: %w", err)` to wrap errors so context is preserved but `errors.Is/As` still work.
- [ ] **Check Errors**: Never ignore errors using `_`. At minimum, log them.

### Concurrency
- [ ] **Context**: Functions doing I/O should accept `context.Context` as the first argument.
- [ ] **Data Races**: Run tests with `go test -race`.
- [ ] **Channel Closing**: Ensure channels are closed by the sender, not the receiver. Avoid closing closed channels.

### API Design
- [ ] **Interfaces**: Accept interfaces, return structs ("Accept interfaces, return concrete types").
- [ ] **Zero Values**: Make structs useful without explicit initialization/constructors where possible.

### Dependency Management
- [ ] **Modules**: Is `go.mod` and `go.sum` updated?
- [ ] **Vendoring**: If using vendoring, is the `vendor/` directory consistent?

---

## 3. References
- [Effective Go](https://go.dev/doc/effective_go)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
