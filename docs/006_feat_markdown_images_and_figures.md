---
date: 2026-03-21T00:00:00
status: "pending"
type: "feature"
title: "Markdown Images and Figure Rendering"
---

# Task 006: Markdown Images and Figure Rendering

## Goal

Add image support for Markdown image syntax and preserve raw HTML `<img>` tags in the generated HTML body.

## Current Evidence

- The current converter does not document support for Markdown image syntax such as `![Tux, the Linux mascot](/assets/images/tux.png)`.
- Raw HTML `<img>` tags embedded in Markdown should remain usable for template-driven content and article bodies.
- The sample template at `samples/blog-template.html` already uses `<figure>` and `<img>` for the cover image, which makes image behavior a practical part of the project's output model.

## Required Behavior

- Render a Markdown image such as `![Tux, the Linux mascot](/assets/images/tux.png)` as:

```html
<img src="/assets/images/tux.png" alt="Tux, the Linux mascot">
```

- Do not wrap a Markdown image in `<p>` when the source line is only the image.
- Preserve raw HTML `<img>` tags found in the Markdown input instead of escaping or rewriting them.
- Add special handling for Markdown image alt text that begins with `figure:`.
- When the alt text starts with `figure:`, strip the prefix, trim surrounding whitespace, and render the image as:

```html
<figure>
  <img src="/imgs/tux.png" alt="Linux mascot known as tux">
  <figcaption>Linux mascot known as tux</figcaption>
</figure>
```

- The `figure:` behavior applies only to Markdown image syntax, not to raw HTML `<img>` tags.
- For `figure:` images, both the rendered `alt` attribute and the `<figcaption>` content must use the trimmed text after the prefix.

## Implementation Notes

- Extend the converter's line processing so Markdown image syntax is recognized before paragraph fallback.
- Keep the implementation conservative and focused on standard inline Markdown image syntax: `![alt](src)`.
- Treat raw HTML `<img>` tags as trusted passthrough content at the converter layer for this task.
- Ensure existing escaping behavior for normal text still applies outside of the explicit raw HTML `<img>` passthrough case.
- Preserve current behavior for templates and metadata handling; this task is limited to body conversion.
- If inline image parsing and block-level line rendering interact, prefer behavior that keeps a line containing only a Markdown image from being enclosed in paragraph tags.

## Acceptance Criteria

- A Markdown image `![Tux, the Linux mascot](/assets/images/tux.png)` renders as `<img src="/assets/images/tux.png" alt="Tux, the Linux mascot">`.
- A Markdown image whose alt text starts with `figure:` renders as a `<figure>` with an `<img>` and `<figcaption>`.
- The `figure:` prefix is removed from both the `alt` attribute and the caption text, and surrounding whitespace is trimmed.
- Raw HTML `<img>` tags in Markdown pass through unchanged into the rendered HTML body.
- A line containing only a standard Markdown image is not wrapped in `<p>` tags.
- Existing paragraph, code block, list, template, and front matter behavior does not regress.

## Validation

- Add focused converter tests for standard Markdown image rendering.
- Add focused converter tests for `figure:` Markdown image rendering.
- Add a regression test for raw HTML `<img>` passthrough.
- Run `go test ./...`.
- Build the CLI with `go build -o bin/md2html .`.
- Validate the rendered output against a sample Markdown input containing a plain Markdown image, a `figure:` Markdown image, and a raw HTML `<img>` tag.