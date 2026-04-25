# Claude Code Instructions

## Project
See [STRUCTURE.md](STRUCTURE.md) for full layout.
See [PHILOSOPHY.md](PHILOSOPHY.md) for guiding principles on data collection and resource efficiency.
`cmd/app/` — Go CLI orchestrating Docker builds and benchmark runs.
Build: `go build ./cmd/app/`

## Plans
Always write plans to `plans/YYYYMMDD-short-description.md` before implementing. Never describe a plan only in chat. Write the file, then stop — don't prompt or ask to proceed.

Update [STRUCTURE.md](STRUCTURE.md) whenever plans or code changes affect project layout or architecture.

When important philosophical decisions emerge during planning or implementation, prompt to add them to [PHILOSOPHY.md](PHILOSOPHY.md).

## Go
Write idiomatic Go.

## Responses
Terse. Lead with the result. Keep this file minimal.

## Models
Choose lowest cost model for each task unless specified.
