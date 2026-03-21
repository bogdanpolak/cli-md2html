---
name: tdd-task-workflow
description: 'TDD workflow for docs/task_*.md. Use for sequential task impl with RED/GREEN/REFACTOR commits, focused tests, sample-based CLI checks, terse updates, minimal unrelated edits.'
argument-hint: 'Which task files should be implemented?'
---

# TDD Task Workflow

Use for one or more `docs/task_*.md` files.

## Use When

- task file exists in `docs/`
- user wants TDD phases
- user wants small commits

## Inputs

- task file paths
- scope / commit constraints if any

## Defaults

- tests: `go test ./...`
- sample: `samples/001-post-pl.md`
- impl: `converter.go`
- cli: `main.go`

## Flow

1. Read task. Extract reqs, acceptance, validation, out-of-scope.
2. If unclear, ask before edits.
3. Do tasks one by one.
4. Keep plan short.
5. RED:
   - add/change tests first
   - run smallest test scope
   - verify fail reason
   - commit `RED: {task title}`
6. GREEN:
   - implement minimum code
   - avoid non-required refactor
   - rerun relevant tests
   - if rendering/CLI touched, validate with `samples/001-post-pl.md`
   - commit `GREEN: {task title}`
7. REFACTOR:
   - do at least one cleanup
   - no behavior change
   - rerun affected tests
   - commit `REFACTOR: {short refactor title}`
8. Report blockers, key decisions, test results, final result.

## Rules

- IMPORTANT: always create a commit for each phase, even not asked. This is the core of the workflow.
- task unclear: clarify first, use tool `askQuestions`
- task vs tests conflict: task doc wins unless clearly wrong
- suite large: start narrow, widen after GREEN
- unrelated issue: ignore unless blocking
- use subagents if they save context
- if using subagent, print one line: `[SubagentName] <current activity>`
- be terse
- no repeated plans
- no long code excerpts unless needed
- report blockers, key decisions, final results

## Done When

- acceptance met
- RED/GREEN/REFACTOR all done
- commit names match format
- relevant tests pass
- CLI checked when applicable
- no unrelated cleanup mixed in
