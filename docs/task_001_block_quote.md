# Task 001: Implement Block Quote Support

## Goal

Add Markdown block quote support so lines starting with `>` are rendered as HTML block quotes.

## Current Evidence

- `TestBlockQuoteConversion` in `converter_test.go` is failing.
- The converter currently treats `>` lines as plain paragraph text and escapes the marker.

## Required Behavior

- Render a single-line quote such as `> text` as `<blockquote>text</blockquote>`.
- Process inline formatting inside block quotes using the existing inline pipeline.
- Treat a quote as a callout only when the quoted text begins with a phrase that ends at the first `:`.
- Ignore callout detection if `.`, `,`, `;`, or `-` appears before that first `:`.
- Wrap the detected callout label, including the `:`, as `<strong>...</strong>` inside the block quote.
- Keep output formatting consistent with the rest of `GenerateHtmlBody`, including trailing newlines.

## Implementation Notes

- Update the main line-processing flow in `GenerateHtmlBody` or `processSingleLine` so quote lines are recognized before paragraph fallback.
- Reuse `processInlineElements` for quote content.
- Detect callouts only on single-line quotes for now. Multi-line quote support is out of scope for this task.
- The callout label is always everything from the start of the quote up to and including the first `:` when no earlier separator invalidates it.

## Acceptance Criteria

- `TestBlockQuoteConversion/01_Basic_Block_Quote` passes.
- `TestBlockQuoteConversion/02_Block_Quote_with_Callout` passes with a properly closed `<strong>` tag.
- Quotes containing `.`, `,`, `;`, or `-` before the first `:` remain regular block quotes and do not trigger callout formatting.
- Quote lines are no longer emitted as escaped paragraph text.
- Inline code, links, bold, and italic formatting still work inside block quotes.
- No regressions in existing header, list, code block, or paragraph tests.

## Validation

- Run `go test ./...`.
- Add or refine tests for multi-line quotes only if support is implemented.
