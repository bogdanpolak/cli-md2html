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

- Keep changes focused and minimal.
- Preserve the existing formatting and naming style.
- Prefer fixing behavior in `converter.go` rather than adjusting tests unless the tests are clearly wrong.
- Do not mix feature work with unrelated cleanup.

## Validation

- Run unit and BDD tests with `go test ./...`.
- If changing CLI behavior, review `features/*.feature` and `steps/step_definitions.go` together.
- For markdown parsing changes, check both `converter_test.go` and README feature claims.

## Current Known Gaps

- Block quotes are asserted in tests but not implemented.
- YAML frontmatter stripping is asserted in tests but not implemented.
