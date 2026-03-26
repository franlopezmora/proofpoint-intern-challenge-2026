# Technical Challenge - Intern Program 2026

This repository is organized by challenge section:

- `part-b/`: Streaming episodes CSV cleaning + deduplication CLI
- `part-c/`: Word frequency analysis CLI

## Part B

```bash
cd part-b
go test ./...
go run ./cmd/episodes-cleaner --input ./input/episodes_input_main.csv --outdir ./output/main
```

## Part C

```bash
cd part-c
go test ./...
go run ./cmd/wordfreq --input ./input/sample.txt
```
