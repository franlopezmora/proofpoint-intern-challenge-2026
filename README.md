# Episodes Cleaner CLI (Go)

This project provides a Go CLI that cleans TV episode records from a corrupted/incomplete CSV and generates:

- `episodes_clean.csv`
- `report.md`

## Run

```bash
go run ./cmd/episodes-cleaner --input ./input/episodes.csv --outdir ./output
```

If `--input` or `--outdir` is missing, the CLI exits with a clear error.

## Outputs

1. `episodes_clean.csv` with exact columns:
   `SeriesName,SeasonNumber,EpisodeNumber,EpisodeTitle,AirDate`
2. `report.md` with metrics and a deduplication summary.

## Project structure

- `cmd/episodes-cleaner/main.go`: CLI entrypoint and flag/error handling.
- `internal/service/cleaner.go`: orchestration of the end-to-end cleaning flow.
- `internal/csvio`: robust CSV read/write (including uneven row lengths).
- `internal/normalize`: text normalization and numeric/date parsing.
- `internal/dedupe`: candidate-key generation, grouping, and best-record selection.
- `internal/report`: report content generation and file output.
- `internal/domain`: domain models.
- `internal/*/*_test.go`: focused unit tests for critical logic.

## Key implementation decisions

- Standard library only.
- Rows with fewer than 5 columns are padded with empty values.
- Rows with more than 5 columns are recovered using:
  - column 0: series name
  - column 1: season number
  - column 2: episode number
  - last column: air date
  - columns in-between: episode title (joined with commas)
- A standard header row (`Series Name, Season Number, Episode Number, Episode Title, Air Date`) is detected and skipped if present.
- Output is deterministic and sorted by:
  `SeriesName asc, SeasonNumber asc, EpisodeNumber asc, EpisodeTitle asc`.

## Corrected entries policy

A record is counted as corrected if at least one field is normalized or defaulted, including:

- trim/collapse whitespace
- text case normalization
- invalid/missing season or episode number -> `0`
- missing title -> `Untitled Episode`
- invalid/missing air date -> `Unknown` (or normalized to ISO `YYYY-MM-DD`)

## Deduplication strategy

Comparison keys use normalized text (`trim + collapse spaces + lowercase`).
Records are grouped as duplicates when any rule matches:

1. `(series, season, episode)`
2. `(series, 0, episode, title)` when `season == 0`
3. `(series, season, 0, title)` when `episode == 0`

For each duplicate group, the kept record is selected by:

1. valid `AirDate` over `Unknown`
2. known `EpisodeTitle` over `Untitled Episode`
3. known `SeasonNumber` and `EpisodeNumber` (both `> 0`)
4. first appearance in input order as final tie-breaker

## Tests

```bash
go test ./...
```

## Reviewer quick run

See `REVIEWER_RUN.md` for copy-paste commands using the included sample CSV files.
