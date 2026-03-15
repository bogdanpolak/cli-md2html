# Task 002: Implement YAML Frontmatter Skipping

## Goal

Ignore leading YAML frontmatter blocks delimited by `---` so metadata does not appear in generated HTML.

## Current Evidence

- `TestYamlFrontmatterConversion` in `converter_test.go` is failing.
- The converter currently emits frontmatter delimiter and metadata lines as normal paragraphs.

## Required Behavior

- Detect a frontmatter block only at the start of the Markdown document.
- Skip both an empty block and a populated metadata block through the closing `---` delimiter.
- Resume normal Markdown processing after the closing delimiter.

## Implementation Notes

- Strip frontmatter before the main render loop in `GenerateHtmlBody`.
- Keep the logic conservative: only treat `---` as frontmatter when it is the opening line of the document.
- Preserve current behavior for `---` appearing later in the body unless new requirements say otherwise.
- The current expected strings in `TestYamlFrontmatterConversion` do not match the provided markdown input text exactly. Confirm and correct those expectations as part of implementation before relying on them.

## Acceptance Criteria

- `TestYamlFrontmatterConversion/01_Empty_Yaml_Frontmatter` passes.
- `TestYamlFrontmatterConversion/02_Sample_Yaml_Frontmatter` passes.
- Frontmatter metadata is absent from the rendered HTML output.
- Documents without frontmatter keep current behavior.

## Validation

- Run `go test ./...`.
- Add a regression test showing that mid-document `---` is not silently stripped.
