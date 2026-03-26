# Part C - Word Frequency Analyzer (Go CLI)

This CLI reads a text file and prints the top 10 most frequent words.

## Run

```bash
go run ./cmd/wordfreq --input ./input/sample.txt
```

## Tokenization and normalization policy

- Text is converted to lowercase.
- A word is any contiguous Unicode alphanumeric sequence (`unicode.IsLetter` or `unicode.IsDigit`).
- Punctuation and special characters are treated as separators.
- Apostrophes (`'`), hyphens (`-`), and underscores (`_`) are separators.

This policy keeps behavior deterministic and simple while matching the challenge requirements.

## Output format

```text
Top 10 words
1. the: 23
2. and: 17
...
```

If fewer than 10 unique words exist, all available words are shown.

## Project structure

- `cmd/wordfreq/main.go`: CLI entrypoint and argument/error handling.
- `internal/textio/reader.go`: text file reading.
- `internal/analyzer/analyzer.go`: tokenization, counting, ranking.
- `internal/analyzer/analyzer_test.go`: unit tests for critical logic.

## Tests

```bash
go test ./...
```
