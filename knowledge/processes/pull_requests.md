# Pull Request Procedure

This document outlines the standard procedure and best practices for creating and reviewing Pull Requests (PRs) in the Sketchnote Artist project.

## Creating a Pull Request

### 1. Atomic Changes

* Keep PRs focused on a single task, bug fix, or feature.
* Avoid "mega-PRs" that touch many unrelated files. Small, focused PRs are easier to review and less likely to introduce bugs.

### 2. Descriptive Titles

* Use [Conventional Commits](commits.md) style for PR titles:
  * `feat: add new summarization tool`
  * `fix: resolve api timeout issue`
  * `docs: update pull request procedure`

### 3. Comprehensive Description

* **What**: Briefly explain the changes.
* **Why**: Explain the rationale (link to issues if applicable).
* **How**: Summarize the technical approach taken.
* **Checklist**:
  * [ ] Unit tests pass
  * [ ] Integration tests pass (run manually)
  * [ ] Documentation updated
  * [ ] Linting and pre-commit hooks pass locally

### 4. Self-Review

* Before requesting a review, read through your own diff.
* Check for console logs, commented-out code, or obvious errors.

## Review Process

### 1. Reviewer Responsibilities

* **Correctness**: Does the code do what it's supposed to?
* **Readability**: Is the code easy to understand?
* **Maintainability**: Does it follow [Go Best Practices](../go/go-best-practices.md)?
* **Tests**: Are there adequate tests for the new logic?

### 2. Etiquette

* **Be Kind**: Provide constructive feedback, not criticism.
* **Ask, Don't Command**: Use "Could we...?" or "What do you think about...?" instead of "Change this."
* **Explain Why**: If requesting a change, provide a reason or a better alternative.
* **Approve if Minor**: If only minor typos are found, approve with a "nit: fix typo" comment.

### 3. Automated Checks

* PRs must pass all **Cloud Build** checks (Linting, Unit Tests, Secret Scanning) before they can be merged.

## Merging

* PRs should be merged using **Squash and Merge** to keep the history clean on the parent branch (`dev` or `main`).
* Ensure the branch is up-to-date with the base branch before merging.
