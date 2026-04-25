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
