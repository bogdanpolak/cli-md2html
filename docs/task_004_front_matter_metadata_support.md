---
date: 2026-03-15T16:00:00
status: "pending"
title: "Front Matter Metadata Support"
---
# Task 004: Front Matter Metadata Support

## Goal

Parse leading YAML front matter into template data so HTML templates can render document metadata in addition to the converted body.

## Current Evidence

- `ConvertMarkdownToHTML` currently executes templates with only `Title` and `Content` populated.
- `GenerateHtmlBody` strips leading YAML front matter, so metadata is removed from the rendered body but never exposed to the template.
- `samples/blog-template.html` references these placeholders: `Title`, `Description`, `Date`, `Author`, `Language`, `CoverImage`, `CoverImageCaption`, `Content`, and `PageFooter`.
- `samples/001-post-pl.md` defines these front matter keys: `postId`, `language`, `date`, `author`, `title`, `description`, `coverImage`, `coverImageCaption`, and `intro`.
- The sample template expects `PageFooter`, but the sample front matter does not define it. That value must therefore resolve to an empty string instead of causing template execution to fail.

## Required Behavior

- Parse a leading YAML front matter block before template execution.
- Keep the current behavior where front matter is not emitted into the HTML body.
- Expose front matter values to templates using exported template field names that match the template placeholders.
- Support the sample template placeholders at minimum: `Title`, `Description`, `Date`, `Author`, `Language`, `CoverImage`, `CoverImageCaption`, `Content`, and `PageFooter`.
- Map sample YAML keys to template fields case-sensitively at the template layer, including at least:
  - `title` -> `Title`
  - `description` -> `Description`
  - `date` -> `Date`
  - `author` -> `Author`
  - `language` -> `Language`
  - `coverImage` -> `CoverImage`
  - `coverImageCaption` -> `CoverImageCaption`
- Leave unsupported or currently unused keys such as `postId` and `intro` harmless unless they are explicitly surfaced by the implementation.
- Any template placeholder that is not defined in the front matter must default to the empty string. Missing metadata must not produce a template execution error.
- Preserve CLI title override semantics. If `-title` is provided, it should continue to take precedence over front matter `title` for the `Title` template field.

## Implementation Notes

- Prefer a single front matter parsing pass that returns both the stripped markdown body and a metadata structure for template execution.
- Keep the template data model focused and explicit. Avoid introducing behavior unrelated to front matter and template population.
- Use a template data representation that tolerates absent fields cleanly, such as a map-backed structure or a struct that preinitializes known fields to empty strings.
- Ensure `Content` remains the converted HTML body and is always available to the template.
- Add tests around both direct `ConvertMarkdownToHTML` usage and CLI behavior if template rendering changes are visible through the command.

## Acceptance Criteria

- A template containing `{{.Author}}`, `{{.Date}}`, `{{.Language}}`, `{{.Description}}`, `{{.CoverImage}}`, and `{{.CoverImageCaption}}` receives values from YAML front matter in `samples/001-post-pl.md`.
- `{{.PageFooter}}` renders as an empty string when the front matter omits that key.
- Template execution no longer fails when optional metadata fields are absent.
- `{{.Content}}` still contains the converted HTML body with the front matter removed.
- `{{.Title}}` resolves from CLI `-title` when provided, otherwise from front matter `title`, otherwise the existing default title behavior.
- Existing markdown conversion behavior remains unchanged for documents without front matter.

## Validation

- Run `go test ./...`.
- Build the CLI with `go build -o bin/md2html .`.
- Generate HTML with `samples/001-post-pl.md` and `samples/blog-template.html`, then verify the rendered output contains populated metadata fields and an empty footer placeholder without template errors.