# Versioning Guidelines

This project follows [Semantic Versioning 2.0.0](https://semver.org/) for managing version numbers.

## VERSION File

The current version of the project is stored in the [VERSION](../../VERSION) file at the project root. This should be the single source of truth for the application's version.

## Versioning Rules

* **Semantic Versioning**: Version numbers are in the format `MAJOR.MINOR.PATCH`.
  * **MAJOR**: Incremented for incompatible API changes.
  * **MINOR**: Incremented for adding functionality in a backwards compatible manner.
  * **PATCH**: Incremented for backwards compatible bug fixes.
* **Increment on Commit**: Every commit that introduces a change should ideally be accompanied by an update to the `VERSION` file, following the rules above.
* **Documentation**: Before making a commit, consult this file and the [Git Commit Guidelines](commits.md) to ensure consistency.
