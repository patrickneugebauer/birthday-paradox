# Claude Code Instructions

## Project
See [docs/STRUCTURE.md](docs/STRUCTURE.md) for full layout.
See [docs/PHILOSOPHY.md](docs/PHILOSOPHY.md) for guiding principles on data collection and resource efficiency.
`cmd/app/` — Go CLI orchestrating Docker builds and benchmark runs.
Build: `go build ./cmd/app/`

## Plans
When you plan, write plans to markdown files in root instead of output in terminal and before implementing. Never describe a plan only in chat. Write the file, then stop — don't prompt or ask to proceed.

Update [STRUCTURE.md](STRUCTURE.md) whenever plans or code changes affect project layout or architecture.

When important philosophical decisions emerge during planning or implementation, prompt to add them to [PHILOSOPHY.md](PHILOSOPHY.md).

## Go
Write idiomatic Go.

## Responses
Terse. Lead with the result. Keep this file minimal.

## Models
Choose lowest cost model for each task unless specified.

## Dockerfiles
Comment any time you work around unexpected behavior or version constraints (e.g., version that doesn't exist, library incompatibilities, base image limitations, unusual command syntax). Document gotchas so future changes aren't made in ignorance of why the current approach exists.

**Approach:** Make it work first (full build environment, no cleanup), then optimize for size incrementally.

**Before creating a Dockerfile:** Always check the official installation instructions:
1. **GitHub README first** — check `projects/*/README.md` or repo root README.md for build/install instructions
2. **GitHub website** — some projects have installation sections in their README
3. **Official website** — find via GitHub API: run `./app stars`, look up language in `artifacts/star-raw-payloads.jsonl`, check the `homepage` field from the GitHub API response
4. Only then try to build from source or use package managers

## Security
Never download repos or arbitrary files directly to the filesystem. For testing tarballs or releases, extract and inspect only within Docker containers.
