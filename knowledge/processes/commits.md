# Git Commit Guidelines

This project follows the [Conventional Commits](https://www.conventionalcommits.org/) specification for clear and consistent commit history.

## Commit Message Template

Save the following content to a file named `.gitmessage` in your project root:

```text
<type>(<scope>): <description>

[optional body]

[optional footer(s)]

# --- Commit Type Guide ---
# feat:     A new feature
# fix:      A bug fix
# docs:     Documentation only changes
# style:    Changes that do not affect the meaning of the code (white-space, formatting, etc)
# refactor: A code change that neither fixes a bug nor adds a feature
# perf:     A code change that improves performance
# test:     Adding missing tests or correcting existing tests
# build:    Changes that affect the build system or external dependencies
# ci:       Changes to our CI configuration files and scripts
# chore:    Other changes that don't modify src or test files
# -------------------------
# Remember:
# - Use the imperative, present tense: "change" not "changed" nor "changes"
# - Separate subject from body with a blank line
# - Limit the subject line to 50 characters
# - Wrap the body at 72 characters
# -------------------------
```

## Local Configuration

To use this template as your default for this repository:

1. Create the template file:

   ```bash
   cp knowledge/processes/commits.md .gitmessage # Or copy the text block above manually
   ```

2. Configure git to use it:

   ```bash
   git config commit.template .gitmessage
   ```

## Best Practices

* **Imperative Mood**: Use the imperative mood in the subject line (e.g., "Add feature" instead of "Added feature").
* **Conciseness**: Keep the subject line under 50 characters.
* **Detail**: Use the body for "what" and "why" explanations, wrapping at 72 characters.
* **Atomic Commits**: Keep commits focused on a single logical change.
