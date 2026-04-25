# JSONL → SQLite3 Migration Plan

## Philosophy

This project prioritizes **data collection** and **low-resource efficiency**:

- **Memory first**: Stream everything. Never load entire files, result sets, or error lists into memory. Process one row at a time, write one row at a time.
- **Minimal dependencies**: Only stdlib + one sqlite3 driver. No ORMs, no frameworks, no extras.
- **Fail loudly**: Errors surface immediately with clear context (which row, which operation). No silent failures, no graceful degradation.
- **Temp then final (batch-consistent)**: Commands write results as 'temp' status during execution. On successful completion, all old 'final' rows are deleted and all temp rows promoted to 'final'. Result: each table contains either all 'final' (previous successful run) or all 'temp' (current run), never a mix. On failure, temp rows remain for inspection/retry.
- **Data preservation**: Keep all data. Orphaned records (from deleted languages/tags) persist and relationships auto-restore via natural keys when data is re-added. No CASCADE DELETE.
- **Idempotent reruns**: Commands can safely run multiple times. INSERT OR REPLACE handles duplicate rows. Previous temp results are overwritten; previous final results are deleted then replaced.
- **Natural keys only**: Use language, tag, (language, status), or (tag, status) as primary keys. No synthetic IDs. Relationships are reversible.
- **One command at a time**: Single database connection per command invocation. No pooling, no concurrency. Sufficient because only one command runs at a time.
- **Atomic transactions**: Each command wraps all writes in a single transaction. On error, accumulated writes are committed (preserving partial progress) and command exits.

## Problem
The readme task currently loads all JSONL files into memory as large slices and maps, then marshals them all to JSON just to read them. This doesn't scale — we need efficient querying and joining across multiple data sources.

## Solution
Add SQLite3 as the primary data store. All commands read from and write to DB only (no JSONL/SH generation). On first run, populate `languages` table from reference CSVs. Queries driven by `Tag` and `Language` as natural primary keys (no synthetic IDs). Keep existing JSONL/SH files untouched (backward compat, but unused after migration).

## Database Schema

### Core Tables

```sql
CREATE TABLE languages (
  language TEXT NOT NULL,
  status TEXT NOT NULL,  -- 'temp' or 'final'
  year INT,
  wiki_url TEXT,
  github_url TEXT,
  PRIMARY KEY (language, status)
);

CREATE TABLE solutions (
  tag TEXT NOT NULL,
  status TEXT NOT NULL,  -- 'temp' or 'final'
  language TEXT NOT NULL,
  dockerfile TEXT NOT NULL,
  runtime TEXT,
  PRIMARY KEY (tag, status)
);

CREATE TABLE builds (
  tag TEXT NOT NULL,
  status TEXT NOT NULL,  -- 'temp' or 'final'
  PRIMARY KEY (tag, status)
);

CREATE TABLE sizes (
  tag TEXT NOT NULL,
  status TEXT NOT NULL,  -- 'temp' or 'final'
  size_mb REAL NOT NULL,
  PRIMARY KEY (tag, status)
);

CREATE TABLE runs (
  tag TEXT NOT NULL,
  status TEXT NOT NULL,  -- 'temp' or 'final'
  iterations INT NOT NULL,
  sample_size INT NOT NULL,
  percent REAL NOT NULL,
  seconds REAL NOT NULL,
  ips INT NOT NULL,
  PRIMARY KEY (tag, status)
);

CREATE TABLE stars (
  language TEXT NOT NULL,
  status TEXT NOT NULL,  -- 'temp' or 'final'
  stars INT NOT NULL,
  PRIMARY KEY (language, status)
);

CREATE TABLE github_api_responses (
  language TEXT NOT NULL,
  status TEXT NOT NULL,  -- 'temp' or 'final'
  response JSON NOT NULL,
  fetched_at TIMESTAMP NOT NULL,
  PRIMARY KEY (language, status)
);

CREATE TABLE readme (
  tag TEXT PRIMARY KEY,
  language TEXT NOT NULL,
  runtime TEXT,
  year INT,
  wiki_url TEXT,
  wiki_display TEXT,
  github_url TEXT,
  github_display TEXT,
  stars INT,
  size_mb REAL,
  ips INT
  -- No foreign keys: denormalized, persists independently
);
```

### Rationale
- **languages**: Joins CSV reference data (years, wiki, github) under a single table. Has status='temp'/'final' to allow monthly updates; only 'final' rows used by other commands.
- **solutions**: Maps tag → language, dockerfile, runtime. Has status='temp'/'final' for consistency with other tables; allows tracking failed/partial map-files runs.
- **builds/sizes/runs/stars**: Composite PK `(tag/language, status)` allows temp and final results to coexist. Status is 'temp' during processing, 'final' when done. Mirrors the old `*-temp-results.jsonl` pattern.
- **github_api_responses**: Stores the full unmodified JSON response from GitHub API per language, with timestamp and status. Allows audit trail and handles schema changes gracefully. Can extract fields with `json_extract` if needed.
- **readme**: Denormalized table (no foreign keys). Stores pre-joined rows with both full URLs and extracted display text (extracted by Go code during readme-data task). Won't cascade-delete if source tables are updated; temp results in other tables don't affect finalized readme rows.
- **No artificial IDs**: Natural keys (tag, language, status) are sufficient; they're already unique.

## Data Flow

### Single Phase: Database-Driven Commands
- `main.go` opens DB connection at startup, closes at exit
- DB file created at `artifacts/bday.db` on first run (idempotent via `CREATE TABLE IF NOT EXISTS`)
- Each command reads from DB (or filesystem for map-files), writes results to DB only
- No JSONL or SH files generated (existing files left untouched)
- Connection held open for command duration, closed when command exits

### README Data & Markdown Generation
The `readme` workflow is split into two independent tasks for resilience:

**`readme-data` task: Populate readme table**
- Single atomic INSERT...SELECT query (no app memory overhead — streams within SQLite)
- Execute within transaction for atomicity
- On success: commit transaction
- On error: commit what succeeded, return error, exit

In the Go code (readme-data task):
```go
// Helper to extract display text from URLs
func extractGitHubDisplay(url string) string {
  if url == "" {
    return ""
  }
  // Extract "owner/repo" from "https://github.com/owner/repo"
  parts := strings.Split(url, "github.com/")
  if len(parts) > 1 {
    return parts[1]
  }
  return ""
}

func extractWikiDisplay(url string) string {
  if url == "" {
    return ""
  }
  // Extract "ArticleName" from "https://en.wikipedia.org/wiki/ArticleName"
  parts := strings.Split(url, "/wiki/")
  if len(parts) > 1 {
    return parts[1]
  }
  return ""
}
```

Then execute INSERT...SELECT that joins all tables:
```sql
INSERT OR REPLACE INTO readme
SELECT
  s.tag, s.language, s.runtime,
  l.year, l.wiki_url, l.github_url,
  st.stars, sz.size_mb, r.ips
FROM solutions s
JOIN languages l ON s.language = l.language AND l.status = 'final'
LEFT JOIN (SELECT language, stars FROM stars WHERE status = 'final') st 
  ON l.language = st.language
LEFT JOIN (SELECT tag, size_mb FROM sizes WHERE status = 'final') sz 
  ON s.tag = sz.tag
LEFT JOIN (SELECT tag, ips FROM runs WHERE status = 'final') r 
  ON s.tag = r.tag
ORDER BY s.tag;
```

Then in Go, read the readme table and populate `_display` columns before writing markdown:
```go
rows, _ := db.Query("SELECT tag, wiki_url, github_url FROM readme")
for rows.Next() {
  var tag, wikiURL, githubURL string
  rows.Scan(&tag, &wikiURL, &githubURL)
  wikiDisplay := extractWikiDisplay(wikiURL)
  githubDisplay := extractGitHubDisplay(githubURL)
  // INSERT or UPDATE readme with display texts
}
```
(INSERT...SELECT is atomic and memory-efficient — SQLite executes the SELECT and INSERT within a single transaction without loading all rows into app memory.)

**`write-readme` task: Format as markdown**
```go
f, _ := NewBufferedFile(readmeFile)
rows, _ := db.Query("SELECT * FROM readme ORDER BY ips DESC")
for rows.Next() {
  var row ReadmeRow
  rows.Scan(&row.Tag, &row.Language, &row.Runtime, &row.Year, &row.WikiURL, &row.WikiDisplay, &row.GitHubURL, &row.GitHubDisplay, &row.Stars, &row.SizeMB, &row.IPS)
  f.WriteString(formatMarkdownRow(row))  // write directly, don't accumulate
}
f.Close()
```
- Stream `SELECT * FROM readme ORDER BY ips DESC` one row at a time
- Format and write each row directly to markdown file via BufferedFile
- Never accumulate rows in memory; one row in memory at a time
- **Table columns in output:** `Tag | Language | Runtime | Year | Wiki | GitHub | Stars | Size (MB) | IPS`
- **Formatting rules:**
  - Runtime: Dockerfile extension (e.g., `clang`, `gcc`), blank if none
  - Stars, Size (MB), IPS: Format with thousands separators (commas) — e.g., `133,582`, `5,260`, `8,500,000`
  - GitHub: `[github_display](github_url)` or `-` if url is null
  - Wiki: `[wiki_display](wiki_url)` or `-` if url is null

**Resilience:** If `readme-data` fails mid-run, partial readme table is committed. User can inspect, fix issues, re-run `readme-data`, then call `write-readme` to format the latest batch.

## Implementation Steps

0. **Update STRUCTURE.md** — Document all schema changes, new database file location, updated data flow, and reorganized task architecture. This should reflect the migration from JSONL files to SQLite3 as primary data store.

1. **Write sqlite3 helpers** (`internal/db/db.go` or similar)
   - `OpenDB(path string)`: Open/create database at `artifacts/bday.db`, run schema (idempotent via `CREATE TABLE IF NOT EXISTS`)
   - `EnsureLanguagesPopulated()`: Populate languages table from CSVs (idempotent, rerunnable monthly)
     - Check: all three CSV files exist (years.csv, github-links.csv, wiki-links.csv) — if missing, error immediately
     - Load years.csv: scan line-by-line, INSERT INTO languages (language, status, year) VALUES (..., 'temp', ...)
     - Load github-links.csv: scan line-by-line, INSERT OR REPLACE INTO languages (language, status, github_url)
     - Load wiki-links.csv: scan line-by-line, INSERT OR REPLACE INTO languages (language, status, wiki_url)
     - All in single transaction
     - On success: DELETE old (language, 'final') rows, UPDATE (language, 'temp') → (language, 'final')
     - On error: rollback and return error message (previous 'final' rows untouched; user can fix and retry)
   - Generic query helpers to keep code terse
   - Use standard `database/sql` package with `github.com/mattn/go-sqlite3` driver (cgo-based, wraps sqlite3 C library directly)

2. **Update main.go**
   - Open DB connection before command dispatch
   - Pass DB handle to all task functions
   - Close DB connection after command completes
   - Each task creates its own transaction via `db.Begin()` (one transaction per command)

3. **Update constants.go**
   - Add DB path: `const dbPath = filepath.Join(artifactDir, "bday.db")`
   - Add readme path: `const readmeFile = filepath.Join(".", "README.md")` or similar to write to repo root

4. **Modify all task functions — Transaction Pattern**
   All tasks follow this pattern (with full error handling):
   ```go
   func BuildWithDB(db *sql.DB) error {
     tx, err := db.Begin()
     if err != nil {
       return fmt.Errorf("begin transaction: %w", err)
     }
     defer func() {
       if r := recover(); r != nil {
         tx.Commit()  // save progress before panicking
         panic(r)
       }
     }()
     
     rows, err := tx.Query("SELECT tag FROM solutions WHERE status = 'final'")
     if err != nil {
       tx.Commit()  // save what succeeded
       return fmt.Errorf("query solutions: %w", err)
     }
     defer rows.Close()
     
     for rows.Next() {
       var tag string
       if err := rows.Scan(&tag); err != nil {
         tx.Commit()  // save what succeeded
         return fmt.Errorf("scan tag: %w", err)
       }
       result, err := DoWork(tag)
       if err != nil {
         tx.Commit()  // save what succeeded
         return fmt.Errorf("work on %s: %w", tag, err)
       }
       if err := tx.Exec("INSERT OR REPLACE INTO builds (tag, status, ...) VALUES (?, 'temp', ...)").Err(); err != nil {
         tx.Commit()  // save what succeeded
         return fmt.Errorf("insert build %s: %w", tag, err)
       }
     }
     if err := rows.Err(); err != nil {
       tx.Commit()  // save what succeeded
       return fmt.Errorf("rows iteration: %w", err)
     }
     
     // Cleanup: promote all temp to final
     if err := tx.Exec("DELETE FROM builds WHERE status = 'final'").Err(); err != nil {
       tx.Commit()
       return fmt.Errorf("delete old builds: %w", err)
     }
     if err := tx.Exec("UPDATE builds SET status = 'final' WHERE status = 'temp'").Err(); err != nil {
       tx.Commit()
       return fmt.Errorf("promote builds: %w", err)
     }
     
     if err := tx.Commit(); err != nil {
       return fmt.Errorf("commit transaction: %w", err)
     }
     return nil
   }
   ```

   Task-specific modifications (all read/write DB only; cleanup always: DELETE old 'final', UPDATE all 'temp' → 'final'):
   - `map-files`: Read solutions/ directory, INSERT to `solutions` (status='temp'). Cleanup deletes old 'final' solutions, promotes temp to final.
   - `build`: Stream DB `solutions` (WHERE status='final'), run docker build, INSERT to `builds` (status='temp'). Cleanup deletes old 'final' builds, promotes temp to final.
   - `weigh`: Stream DB `solutions` (WHERE status='final'), run docker inspect, INSERT to `sizes` (status='temp'). Cleanup deletes old 'final' sizes, promotes temp to final.
   - `run`: Stream DB `solutions` (WHERE status='final'), run docker container, INSERT to `runs` (status='temp'). Query DB for previous IPS of specific tag. Cleanup deletes old 'final' runs, promotes temp to final.
   - `stars`: Stream DB `languages` (WHERE status='final'), fetch GitHub API, INSERT to `github_api_responses` and `stars` (status='temp'). Cleanup deletes old 'final' rows from both tables, promotes temp to final.
   - `readme-data`: Atomic INSERT OR REPLACE INTO readme (no status column, naturally idempotent).
   - `write-readme`: Stream `readme` table, format as markdown, write to `readme.md` via BufferedFile. No transaction needed.

5. **Data flow & Cleanup**
   - No JSONL or SH file generation; all data persists in database
   - Existing JSONL/SH files left untouched (unused)
   - Commands can be rerun independently; DB is source of truth
   - **Status cleanup on success (batch-consistent):** At end of each task (map-files, build, weigh, run, stars):
     ```sql
     -- Delete all old 'final' rows from current batch
     DELETE FROM <table> WHERE status = 'final';
     -- Promote all 'temp' rows to 'final'
     UPDATE <table> SET status = 'final' WHERE status = 'temp';
     ```
   - **Result:** At any time, table contains **either** all 'final' rows (successful run) **or** all 'temp' rows (current run in progress). Never a mix.
   - **On error:** Leave 'temp' rows in transaction; user inspects, fixes issue, reruns command (INSERT OR REPLACE overwrites old temp rows)
   - **Rerun:** Previous temp rows are overwritten; previous final rows are deleted then replaced

## Philosophy
See [PHILOSOPHY.md](../PHILOSOPHY.md) for guiding principles on data preservation and memory efficiency that inform this design.

## Constraints
- **No CASCADE DELETE**: Schema uses natural keys; orphaned records are reversible via re-adding language/tag.
- **No JSONL or SH files**: All writes to DB only. Existing files left untouched for backward compatibility.
- **Status cleanup pattern**: On success, delete old 'final' rows and promote 'temp' → 'final'. On error, leave 'temp' rows for inspection/retry.
- **Error handling**: Surface DB errors with context (which row, which operation). Commit accumulated writes before exiting on error.

## Testing Notes

**Schema & Initialization**
- Test `OpenDB()` creates schema on first run (idempotent)
- Test `EnsureLanguagesPopulated()`:
  - Fails loudly if any CSV missing
  - Streams all three CSVs correctly, all in single transaction
  - First run: inserts all languages with status='temp', then promotes to 'final'
  - Rerun (monthly): inserts/replaces new/updated languages as 'temp', then cleanup promotes to 'final' and deletes old 'final'
  - Verify INSERT OR REPLACE handles duplicate languages across reruns

**Memory Efficiency**
- Verify `EnsureLanguagesPopulated()` doesn't load all CSVs into memory at once
- Verify streaming commands (build, weigh, run, stars) process one row at a time
- Verify readme-data INSERT...SELECT is atomic (no app memory overhead)
- Verify write-readme streams markdown row-by-row (no accumulation)

**Error Handling & Resilience**
- Test: run build on empty solutions table → errors loudly with "no solutions to build"
- Test: build fails mid-run (e.g., docker error on 3rd tag)
  - Verify first 2 tags committed as status='temp'
  - Verify error message includes which tag failed
  - Verify user can inspect temp rows, fix issue, rerun
  - Verify rerun with INSERT OR REPLACE overwrites temp rows for first 2 tags
  - Verify rerun succeeds and status cleanup promotes temp to final
- Test: successful run after failed run
  - Verify old status='final' rows deleted
  - Verify status='temp' rows promoted to 'final'
  - Verify no 'temp' rows remain after success

**Status & Cleanup Logic**
- Test: first run completes successfully
  - Verify all rows written with status='temp'
  - Verify cleanup runs: DELETE old 'final', UPDATE 'temp' → 'final'
  - Verify final state: only status='final' rows exist
- Test: rerun after success
  - Verify no temp rows exist (previous run was clean)
  - Verify INSERT OR REPLACE overwrites final rows with new final rows
- Test: partial failure scenario (build succeeds, weigh fails)
  - Verify builds table has status='final' rows
  - Verify sizes table has no rows (weigh failed before completing)
  - Verify user can fix and rerun weigh independently

**Data Integrity & Natural Keys**
- Run build → verify builds table populated with status='final'
- Run weigh → verify sizes table populated with status='final'
- Run run → verify runs table populated with status='final'
- Run stars → verify github_api_responses and stars both populated with status='final'
- Run readme-data → verify readme table contains joined result of all above
- Run write-readme → verify readme.md generated correctly
- Test: delete a language from languages table
  - Verify stars/github_api_responses records for that language persist (orphaned)
  - Re-add language to languages table
  - Verify relationship automatically re-establishes (natural key match)
