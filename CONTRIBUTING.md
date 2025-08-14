# Contributing to docker-lint

Thank you for your interest in contributing! This guide explains how to work on **docker-lint** in a way that 
maintains code quality, feature parity with Hadolint, and compatibility with automated build agents like ChatGPT 
Codex.

---

## Development Workflow

1. **Fork & Clone**

    * Fork the repo and clone it locally.
    * Set up the Go toolchain (Go 1.24+).

2. **Understand the Task Plan**

    * The build plan is defined in `.codex/plan.yml`.
    * Tasks must be completed **in order**, respecting dependencies.
    * Each task is small and verifiable — do not bundle multiple unrelated changes.

3. **Guardrails**

    * All Go files must have the standard two-line header:

      ```go
      // file: <filepath/name>
      // (c) 2025 Asymmetric Effort, LLC. <scaldwell@asymmetric-effort.com>
      ```
    * Every function, method, and type must have a proper docstring.
    * Never delete functionality when refactoring.
    * No network access in default code paths.
    * Do not modify `LICENSE` or `docs/parity.md` without approval.

4. **Coding Standards**

    * Use idiomatic Go formatting (`gofmt`).
    * Pass `golangci-lint` checks.
    * Keep commits small and atomic.
    * Write modular code; one function per `.cpp`/`.go` file where applicable.

5. **Testing Requirements**

    * Maintain ≥80% **integration** coverage.
    * Run:

      ```bash
      go test ./... -covermode=atomic -coverpkg=./... -coverprofile=coverage.out
      go tool cover -func=coverage.out
      ```
    * Integration tests live in `test/integration/`.
    * Use `testdata/` for golden files and fixtures.

6. **Plugin Development**

    * Follow the `internal/plugins` API contract.
    * Each plugin must have:

        * A unique `ID()`
        * `Languages()` array
        * Deterministic `Analyze()`
    * Respect size/time limits and map findings back to Dockerfile positions.

7. **Configuration**

    * Keep `internal/config` YAML schema backwards compatible.
    * Document any changes in `docs/configuration.md`.

8. **Running Locally**

    * Use the provided `Taskfile.yml` to run common targets:

      ```bash
      task all
      ```
    * Or run individual tasks per `.codex/plan.yml`.

9. **Pull Requests**

    * Reference the task ID(s) from `.codex/plan.yml` in your PR description.
    * Include before/after test output when modifying rules or formatters.
    * Ensure CI passes (matrix build + coverage gate).

---

## Working with Codex/AI Agents

* Agents follow `.codex/plan.yml`.
* Each task is reviewed for:

    * File header compliance
    * Docstrings
    * Passing tests for the affected modules
* Avoid making changes outside the defined task scope.

---

## License

By contributing, you agree that your contributions will be licensed under the repository’s license (MIT).
