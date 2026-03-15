# AGENTS

## Project Summary

`cli-md2html` is a small Go CLI that converts Markdown into HTML with optional templating and BDD coverage for CLI behavior.

## Repository Layout

- `main.go`: CLI flags, file I/O, preview mode.
- `converter.go`: Markdown to HTML conversion logic.
- `converter_test.go`: unit and integration tests for conversion behavior.
- `bdd_test.go`, `features/`, `steps/`: Godog scenarios for CLI and template behavior.
- `samples/`: sample Markdown inputs.
- `ekon-template.html`: example HTML template.

## Working Rules

- Clarify before starting work.
- Use subagents, reduce main-context usage.
- Keep changes focused and minimal.
- Preserve the existing formatting and naming style.
- Prefer fixing behavior in `converter.go` rather than adjusting tests unless the tests are clearly wrong.
- Do not mix feature work with unrelated cleanup.
- After CLI-affecting changes, validate against samples in `samples/`.

## Tasks

- tasks are in `docs/`, include all implementation details.
- each task should be self-contained, feature or bug fix.
- when task is created, it needs YAML front matter with date, status, and title.
- status = "pending" | "in_progress" | "completed".
- title should be concise.
- date should be in ISO format.
- once task is completed, update status.

## Validation

- Build the project with `go build -o bin/md2html .`.
- Run unit and BDD tests with `go test ./...`.
- If changing CLI behavior, review `features/*.feature` and `steps/step_definitions.go` together.
- For markdown parsing changes, check both `converter_test.go` and README feature claims.
