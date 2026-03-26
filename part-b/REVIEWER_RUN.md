# Reviewer Quick Run

## Prerequisite
- Go 1.22+ installed

## Verify tests
```bash
go test ./...
```

## Run with included sample CSV files
```bash
go run ./cmd/episodes-cleaner --input ./input/episodes_input_smoke.csv --outdir ./output/smoke
go run ./cmd/episodes-cleaner --input ./input/episodes_input_main.csv --outdir ./output/main
go run ./cmd/episodes-cleaner --input ./input/episodes_input_edge_cases.csv --outdir ./output/edge_cases
```

## Expected generated files
- `output/*/episodes_clean.csv`
- `output/*/report.md`
