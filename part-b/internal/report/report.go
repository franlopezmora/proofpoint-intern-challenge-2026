package report

import (
	"fmt"
	"os"
	"path/filepath"
)

type Metrics struct {
	TotalInputRecords  int
	TotalOutputRecords int
	DiscardedEntries   int
	CorrectedEntries   int
	DuplicatesDetected int
	DeduplicationNotes string
}

func Write(outDir string, metrics Metrics) (string, error) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("create output directory: %w", err)
	}

	path := filepath.Join(outDir, "report.md")
	content := buildReport(metrics)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("write report file: %w", err)
	}

	return path, nil
}

func buildReport(metrics Metrics) string {
	return fmt.Sprintf(`# Cleaning Report

- Total number of input records: %d
- Total number of output records: %d
- Number of discarded entries: %d
- Number of corrected entries: %d
- Number of duplicates detected: %d

## Deduplication strategy
%s

## Corrected entries policy
An entry is counted as corrected when at least one field was normalized or defaulted (trim/collapse spaces, case normalization, number fallback to 0, title fallback to "Untitled Episode", or air date fallback to "Unknown"/normalized ISO date).
`,
		metrics.TotalInputRecords,
		metrics.TotalOutputRecords,
		metrics.DiscardedEntries,
		metrics.CorrectedEntries,
		metrics.DuplicatesDetected,
		metrics.DeduplicationNotes,
	)
}
