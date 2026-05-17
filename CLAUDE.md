# Claude Code Instructions

## Git
never commit or push unless explicitly askd to

## Project
`cmd/app/` — Go CLI orchestrating Docker builds and benchmark runs.
Build: `go build ./cmd/app/`

## Go
Write idiomatic Go. Write idiomatic code, structure code in idiomatic patterns, reach for idiomatic solutions.

## Responses
Terse. Lead with the result.

## Explicit
Do not show a dash for a zero value, only for a null or missing one.  Don't show a zero when something was null.

## Security
Never download repos or arbitrary files directly to the filesystem or install libraries. For testing tarballs or releases, extract and inspect only within Docker containers.

## Data Formats
For structured data files, prefer JSON over TSV. JSON is much safer for editing (no column misalignment issues, self-describing fields, less prone to corruption). When working with language metadata, use `all-language-data.json` as the source of truth.

# Project Structure

## Top-Level Directories

| Path | Contents |
|---|---|
| `hello-worlds/` | New solutions start here.  One folder per language — each contains Dockerfiles and hello world source code |
| `scaffolds` | 2nd step for new solutions. More refined than hello-worlds, contains most strucutre, missing middle data portion |
| `solutions/` | Final place for new solutions. |
| `cmd/app/` | Go CLI entry point (`main.go`) — dispatches commands: `map-files`, `stars`, `build`, `weigh`, `run`, `readme` |
| `internal/tasks/` | Go task implementations |
| `tables/` | Generated sorted README views |
| `artifacts/` | Generated output files (git-ignored): JSONL results, `.sh` command logs |

## Key Files

| File | creator | consumer |
|---|---|---|
| `language-info.csv` | source of information for languages | Readme |
| `artifacts/dockerfiles.jsonl` | MapFiles | Build, Readme |
| `artifacts/build-results.jsonl` | Build | Weigh, Run |
| `artifacts/weigh-results.jsonl` | Weigh | Readme |
| `artifacts/run-results.jsonl` | Run | Readme |
| `artifacts/star-results.jsonl` | Star | Readme |
| `artifacts/star-raw-payloads.jsonl` | Star | - |
| `readme.md` | Readme| - |

## Dockerfile Naming Convention

- Dockerfiles name pattern: `Dockerfile[.{runtime}][.{data_structure}][.{execution_method}]`
- tag pattern: `bday/{language}:{runtime}[.{data_structure}][.{execution_method}]` all lowercase

Examples:
- `bday/bash:bash.associative-array`
- `bday/c:clang.array`
- `bday/csharp:dotnet.array.build`
- `bday/chapel:chapel` (when data_structure = `-`)

## Guiding Principles
- avoid external dependencies when possible
- fail loudly
- write to named tempfiles during execution, move over to final file after

## Dockerfile Best Practices
- Minimize size of final image.
- Single stage builds.
- base: Use official image if it exists, if not, build from alpine/debian-slim/node, etc. Choose the slimmest sensible base image: don't fight to get alpine working if it will work on debian slim without a fight.

## Dockerfile sections
- link at top to where the installation instructions came from
- from, workdir app (or use what base recommends), run (install), copy, run (compile), entrypoint

## Source code Sections:
each starts with a comment: vars, data, calcs, format, output

## Do not delete
do not delete the database or files of saved results
