# Project Philosophy

This project prioritizes **data collection** and **low-resource efficiency**.

## Guiding Principles

**Memory first**: Stream everything. Never load entire files, result sets, or error lists into memory. Process one row at a time, write one row at a time.

**Minimal dependencies**: Only stdlib + one sqlite3 driver. No ORMs, no frameworks, no extras.

**Fail loudly**: Errors surface immediately with clear context (which row, which operation). No silent failures, no graceful degradation.

**Temp then final (batch-consistent)**: Commands write results as 'temp' status during execution. On successful completion, all old 'final' rows are deleted and all temp rows promoted to 'final'. Result: each table contains either all 'final' (previous successful run) or all 'temp' (current run), never a mix. On failure, temp rows remain for inspection/retry.

**Data preservation**: Keep all data. Orphaned records (from deleted languages/tags) persist and relationships auto-restore via natural keys when data is re-added. No CASCADE DELETE.

**Idempotent reruns**: Commands can safely run multiple times. INSERT OR REPLACE handles duplicate rows. Previous temp results are overwritten; previous final results are deleted then replaced.

**Natural keys only**: Use language, tag, (language, status), or (tag, status) as primary keys. No synthetic IDs. Relationships are reversible.

**One command at a time**: Single database connection per command invocation. No pooling, no concurrency. Sufficient because only one command runs at a time.

**Atomic transactions**: Each command wraps all writes in a single transaction. On error, accumulated writes are committed (preserving partial progress) and command exits.

## Data Integrity

No CASCADE DELETE constraints are defined; orphaned records are intentionally preserved. Natural keys (language, tag, status) make relationships reversible:
- If a language is deleted from `languages` but its `stars` records remain with language='Python', re-adding Python later will automatically re-establish the relationship.
- Synthetic IDs would make this recovery impossible; natural keys preserve historical continuity.
- Temp and final results are independent, allowing safe reruns and updates without data loss.

## Memory Efficiency

Tasks stream SQL results row-by-row without loading all records into memory. Writes are batched in transactions and committed once per task completion. On error mid-loop: commit accumulated writes, log error, exit (no error list kept in memory).

## Benchmark Dimensions

The project benchmarks along **4 independent dimensions**. Each Dockerfile represents a unique point in this 4D space:

1. **Language**: The programming language (bash, C, Python, Go, etc.)
2. **Runtime**: The language implementation — compiler, interpreter, or VM (e.g. gcc, clang, cpython, node, dotnet, jvm)
3. **DataStructure**: The algorithm variant used in the solution (e.g. array, list, vector, set, associative-array, hash-table)
4. **ExecutionMethod**: Optional; only populated when the same language/runtime/data_structure can be executed in multiple ways (e.g. build, publish, run, ocamlopt, mix-build)

Dockerfile filenames encode these dimensions as `Dockerfile[.{runtime}][.{data_structure}][.{execution_method}]`. Null dimensions are omitted from the filename and tag. Docker tags are constructed by combining all four dimensions: `bday/{language}:{runtime}[.{data_structure}][.{execution_method}]` with all tag components lowercased. Periods separate all dimensions within the tag.

Examples:
- `bday/python:cpython.set`
- `bday/bash:bash.associative-array`
- `bday/csharp:dotnet.array.build`

When adding new solutions, name Dockerfiles directly according to the convention without needing to update reference data.

## Dockerfile Best Practices

**Single-stage builds**: All Dockerfiles are single-stage. No multi-stage builds with builder images. If the build fails, the error should surface clearly in the build log. Within a single stage, it's fine to delete build-time dependencies after they're no longer needed (e.g. `apk del gcc cmake` after compiling).

**Minimal official base images**: Choose the slimmest sensible base image:
- Start with `alpine` or `slim` variants where available (e.g. `python:3.11-slim`, `node:20-alpine`).
- Use `alpine` only if dependencies are available there; if the language or significant dependencies require many additional packages, use a `slim` variant (e.g. `debian:bookworm-slim`) instead. Don't fight alpine; pick the right base.
- Always use **official images** from Docker Hub (golang, python, node, rust, etc.). Never pull from community, unofficial, or third-party image repositories.

**Predictable, maintainable**: Each Dockerfile should be self-contained and easy to understand at a glance. No clever layering tricks. Use `apk add --no-cache` (alpine) or `apt-get install -y --no-install-recommends` (debian) to keep images lean.

**COPY before source files**: All setup steps (package installation, compilation, downloads) must come before the `COPY` of source files. This preserves Docker layer caching: developers can modify the source file and rebuild without re-running expensive setup steps. Always use `COPY`, not `ADD`, unless extracting archives is required.

**WORKDIR placement and directory**: `WORKDIR` must be set immediately after `FROM`, before any `RUN` or `COPY` commands. This establishes context for all subsequent commands and preserves layer caching when source files change.

Use `/app` as the standard working directory. Exception: if the base image is specifically tied to a user/path (e.g. `ocaml/opam:alpine` uses `/home/opam`), respect that path only if it's required by the image's entrypoint or tooling.
