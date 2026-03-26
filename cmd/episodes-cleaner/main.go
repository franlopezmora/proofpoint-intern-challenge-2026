package main

import (
	"flag"
	"fmt"
	"os"

	"proofpoint-flm/internal/service"
)

func main() {
	var inputPath string
	var outDir string

	flag.StringVar(&inputPath, "input", "", "Path to input CSV file")
	flag.StringVar(&outDir, "outdir", "", "Output directory for episodes_clean.csv and report.md")
	flag.Parse()

	if err := service.ValidateArgs(inputPath, outDir); err != nil {
		exitWithError(err)
	}

	cleaner := service.NewCleaner()
	result, err := cleaner.Run(inputPath, outDir)
	if err != nil {
		exitWithError(err)
	}

	fmt.Printf("Generated %s\n", result.CleanCSVPath)
	fmt.Printf("Generated %s\n", result.ReportPath)
	fmt.Printf("Input: %d | Output: %d | Discarded: %d | Corrected: %d | Duplicates: %d\n",
		result.Metrics.TotalInputRecords,
		result.Metrics.TotalOutputRecords,
		result.Metrics.DiscardedEntries,
		result.Metrics.CorrectedEntries,
		result.Metrics.DuplicatesDetected,
	)
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
