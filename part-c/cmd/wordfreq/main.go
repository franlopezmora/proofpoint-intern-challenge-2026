package main

import (
	"flag"
	"fmt"
	"os"

	"proofpoint-flm-part-c/internal/analyzer"
	"proofpoint-flm-part-c/internal/textio"
)

func main() {
	var inputPath string

	flag.StringVar(&inputPath, "input", "", "Path to input text file")
	flag.Parse()

	if inputPath == "" {
		exitWithError(fmt.Errorf("missing required argument: --input"))
	}

	text, err := textio.ReadFile(inputPath)
	if err != nil {
		exitWithError(err)
	}

	counts := analyzer.CountWords(text)
	top := analyzer.TopN(counts, 10)

	fmt.Println("Top 10 words")
	for i, item := range top {
		fmt.Printf("%d. %s: %d\n", i+1, item.Word, item.Count)
	}
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
