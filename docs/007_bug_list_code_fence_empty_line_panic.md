---
date: 2026-03-22T00:00:00
status: "pending"
type: "bug"
title: "Bug: Panic on Empty Lines in List-Nested Fenced Code Blocks"
---

# Task 007: Bug: Panic on Empty Lines in List-Nested Fenced Code Blocks

## Goal

Fix a runtime panic that occurs when converting Markdown containing a fenced code block nested inside a list item when the fenced code contains empty lines.

## Reproduction

Build and run:

```bash
go build -o bin/md2html .
cat samples/article-002.md | ./bin/md2html > /dev/null
```

Observed result:

- CLI panics with:

```text
panic: runtime error: slice bounds out of range [4:0]
main.processCodeBlock(...)
converter.go:264
```

Input file: `samples/article-002.md`

## Suspected Root Cause

Fenced code block with empty lines nested inside a list item. Example bellows.

In `processCodeBlock` (`converter.go`), `depth` is derived from the indentation of the opening fence line. During code line processing, each line is sliced with `lines[idx][depth:]`.

For empty lines or lines shorter than `depth`, this causes an out-of-range slice and crashes the CLI.

## Example Markdown snippet that triggers the panic:

1. **Create the `game` object.**
   - Create a a new function that returns a `game` object, and test it.
    ```js
    test("createGame returns expected shape", () => {
      const game = createGame(["Luke", "Leia"]);

      expect(game.players).toHaveLength(2);
    });
    ```
2. **Move the `players` array into the `game` object.**
   - Refactor the code

## Required Behaviour

- Conversion must never panic for valid Markdown input containing fenced code blocks in list items.
- Empty lines inside fenced code blocks must be preserved safely.
- Language fences such as ```` ```js ```` inside list items must continue to work.
- Existing code block rendering format (`<div class="code"><pre><code>...</code></pre></div>`) must remain unchanged.

## Acceptance Criteria

- `cat samples/article-002.md | ./bin/md2html > /dev/null` exits successfully (no panic).
- Converting `samples/article-002.md` no longer panics.
- Add a focused regression test in `converter_test.go` for a list-nested fenced code block with at least one empty line in the fenced body.
- `go test ./...` passes.
- `go build -o bin/md2html .` succeeds.

## Implementation Notes

- Prefer fixing the behavior in `converter.go` rather than weakening tests.
- Keep the fix minimal and localized to code-block processing logic.
- Ensure the indentation trim operation is bounds-safe for short or empty lines.
