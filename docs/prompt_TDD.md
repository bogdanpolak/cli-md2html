Tasks:
{provide a list of tasks to implement, for example: task_001_block_quote_support.md}

Your job is to implement selected tasks, execute sequentially, for each task:
- clarify if needed
- plan implementation
- write tests and run for the requirements (RED phase)
- git commit, message format `RED: {title of the task}` - for example `RED: Block quote support`
- implement requirements using TDD approach
- run tests
- validate cli tool using #file:001-post-pl.md (it has both cases)
- in that phase reduce code refactorings to required only changes, next phase `REFACTOR` is dedicated for that
- git commit, message format `GREEN: {title of the task}` - for example `GREEN: Block quote support`
- Review all refactorings and apply at least one - in TDD style (it can be main logic or tests)
- git commit, message format `REFACTOR: {title of the refactoring}` - for example `REFACTOR: Simplify block quote tests`

# Working Rules

- ask any clarifying question before start
- use subagents to save context window
- when a subagent is used, show one short status line: `[SubagentName] <current activity>`
- maximize token efficiency. Prefer terse bullets, no repeated plans, no long code excerpts unless required
- report only blockers, key decisions, and final results
