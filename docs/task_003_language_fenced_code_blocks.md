# Task 003: Support Language-Qualified Fenced Code Blocks

## Goal

Handle fenced code blocks that start with a language hint such as ` ```pascal ` so they render as normal code blocks instead of paragraph text.

## Current Evidence

- The sample file `samples/001-post-pl.md` contains multiple fences opened with ` ```pascal `.
- The converter only recognizes a code fence when the trimmed line is exactly ` ``` `.
- Generated output in `samples/001-post-pl.html` shows opening fence lines rendered as paragraphs like `<p>```pascal</p>`.
- The next plain closing fence is incorrectly treated as the end of a code block, so normal prose between fenced examples is swallowed into `<section class="code">`.

## Reproduction

1. Build the CLI: `go build -o bin/md2html .`
2. Generate HTML: `./bin/md2html -input samples/001-post-pl.md -output samples/001-post-pl.html`
3. Inspect the first code example in the generated HTML.

Observed behavior:

- The opening ` ```pascal ` line is emitted as a paragraph.
- Code lines are emitted as separate paragraphs instead of a `<pre><code>` block.
- The following plain ` ``` ` line starts a code block late and captures non-code prose until the next matching fence.

Expected behavior:

- Opening fences with an optional language suffix should start a code block.
- Closing fences should still end the current code block.
- The language suffix may be ignored for now if syntax highlighting is out of scope, but the block structure must remain correct.

## Implementation Notes

- Extend fence detection in `converter.go` so it accepts lines that begin with triple backticks and optionally continue with a language token.
- Preserve current behavior for plain ` ``` ` fences.
- Ensure list-aware code block handling also uses the same fence detection rules.
- Add tests for language-qualified fences both at top level and inside lists.

## Acceptance Criteria

- A markdown block opened with ` ```pascal ` and closed with ` ``` ` renders as one `<section class="code">` containing `<pre><code>...</code></pre>`.
- The prose immediately following such a block is no longer captured inside the code section.
- Existing plain fenced code block tests continue to pass.
- New regression tests cover at least one language-qualified top-level block and one language-qualified block inside a list.

## Validation

- Run `go test ./...`.
- Regenerate `samples/001-post-pl.html` and verify the first fenced examples render as code sections instead of paragraphs.