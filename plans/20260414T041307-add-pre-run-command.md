# Plan: Project Markdown Document + Pre-Run Command

## Context

The birthday-paradox project benchmarks the birthday paradox algorithm across 50+ language implementations using Docker containers. The CLI (`cmd/app/main.go`) supports three commands: `pre-build`, `build`, and `run`.

Currently `pre-build` generates both `build.sh` and `run.sh` (with a hardcoded 1000 iterations). The `run` command then reads `run.sh`, but overrides the iteration count per language by reading previous IPS stats from `results.json` at runtime. The goal is to:
1. Add a `pre-run` command that owns `run.sh` generation (with scaled iteration counts baked in).
2. Remove `run.sh` generation from `pre-build`.
3. Remove IPS-scaling logic from `run` — iteration counts are now baked into `run.sh` by `pre-run`. The `run` command still executes each line individually, captures stdout, parses results, and maintains the three-file write pattern (`temp-results.json` → `results.json`, with previous `results.json` archived to `previous-results.json` by `pre-run`).

---

## Changes

### 1. Create `internal/tasks/pre-run.go`

New `PreRun()` function that:
- Calls the existing `discoverBuilds(solutionsDir)` to get all `BuildConfig` entries (reuses existing code in `pre-build.go`).
- Loads IPS stats from `results.json` (the current results, not `previous-results.json`) using `loadPreviousStats()` (moved from `run.go` — or kept in `run.go` and made package-level, since both files are in the same `tasks` package).
- For each config, selects `iters = prevStats[tag]` if available, else `defaultIterations` (1000).
- Writes `run.sh` with one `docker run --rm <tag> <iters>` line per image using `saveScript()`.

### 2. Modify `internal/tasks/pre-build.go`

- Remove the `runLines` slice and the `saveScript(runScript, runLines, false)` call.
- Keep `build.sh` and `dockerfiles.yaml` generation unchanged.
- Update the log line to say "Generated N build targets." instead of "N build/run targets."

### 3. Modify `internal/tasks/run.go`

- Remove the IPS-scaling logic: the `loadPreviousStats()` call, the `iters` override, and `fields[len(fields)-1] = strconv.Itoa(iters)`.
- Keep the archiving step (`os.Rename(resultsFile, previousResultsFile)`) in `run` — it belongs here because `run` is what writes new results. This preserves the three-file pattern (`previous-results.json`, `temp-results.json`, `results.json`) across any number of consecutive `run` invocations without requiring `pre-run` to be called in between.
- Execute each line from `run.sh` directly via `exec.Command(fields[0], fields[1:]...)` without modification.
- Keep `parseOutput()` and the temp-file write pattern unchanged.
- `loadPreviousStats()` moves to `pre-run.go` since `run.go` no longer needs it. Both files are in the same `tasks` package so it remains accessible.

### 4. Register `pre-run` in `cmd/app/main.go`

Add a `case "pre-run": finalErr = tasks.PreRun()` branch to the switch statement.

---

## Critical Files

| File | Change |
|------|--------|
| `internal/tasks/pre-run.go` | **Create** |
| `internal/tasks/pre-build.go` | Remove run.sh generation (lines 39–49) |
| `internal/tasks/run.go` | Remove scaling logic and archive step |
| `cmd/app/main.go` | Add `pre-run` case |

---

## Verification

1. Run `go build ./cmd/app/` — must compile with no errors.
2. Run `./app pre-build` — should produce `build.sh` and `dockerfiles.yaml`, but **not** `run.sh`.
3. Run `./app pre-run` — should produce `run.sh` with iterations scaled from `results.json` (or 1000 if no results yet). Confirm `run.sh` exists and contains correct docker run lines.
4. Run `./app run` — should execute `run.sh` as-is, saving results to `results.json`.
5. Run `./app pre-run` again — `run.sh` should now reflect IPS values from the just-saved `results.json`.
