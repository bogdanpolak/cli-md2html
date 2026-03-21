---
date: 2026-03-21T00:00:00
status: "complete"
type: "bug"
title: "Bug: Render Code Blocks With div Wrapper"
---

# Task 005: Bug: Render Code Blocks With div Wrapper

## Goal

Standardize fenced code block HTML so both top-level and list-nested code blocks render with `<div class="code">` instead of `<section class="code">`.

- Code blocks are expected to render with a dedicated wrapper class `code`, but the wrapper element is being changed from `section` to `div`.
- The change touches both top-level fenced code blocks and fenced code blocks nested inside ordered or unordered lists.

## Expected Behavior

- A fenced code block opened and closed with triple backticks must render inside `<div class="code">`.
- The same wrapper must be used when the fenced code block appears inside a list item.
- Existing escaping behavior inside `<pre><code>` must remain unchanged.
- Tests that verify HTML output should consistently assert against `<div class="code">`.

## Acceptance Criteria

- Top-level fenced code blocks render as `<div class="code">` containing `<pre><code>...</code></pre>`.
- List-nested fenced code blocks also render as `<div class="code">` with the same inner structure.
- No existing converter tests for code blocks fail because of mixed `section` versus `div` wrappers.
- The complex document test asserts the new wrapper consistently.
- `go test ./...` passes with the updated converter and test expectations.
- Command-line BDD scenarios do not require a separately prepared `bin/md2html` binary before the test run.
