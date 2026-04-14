# Claude Code Instructions

## Project
`solutions/` — one folder per language, each a Monte Carlo birthday problem simulation.
`cmd/app/` — Go CLI orchestrating Docker builds and benchmark runs.
Build: `go build ./cmd/app/`

## Plan-First (Strictly Enforced)
Never make code changes without a plan file. Every task is two prompts:
1. **Plan:** Enter plan mode, explore, write plan to `plans/`, call ExitPlanMode. No code changes.
2. **Execute:** Implement the plan from prompt 1.

If asked to make changes without a plan, interrupt and ask for a planning prompt first.

## Plan Files
`plans/YYYYMMDDTHHMMSS-short-description.md`
Example: `plans/20260414T143022-add-pre-run-command.md`

## Go
Write idiomatic Go.

## Responses
Terse. Lead with the result. Keep this file minimal.
