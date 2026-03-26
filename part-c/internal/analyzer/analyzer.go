package analyzer

import (
	"sort"
	"strings"
	"unicode"
)

type WordCount struct {
	Word  string
	Count int
}

func CountWords(text string) map[string]int {
	counts := make(map[string]int)
	for _, token := range Tokenize(text) {
		counts[token]++
	}
	return counts
}

func TopN(counts map[string]int, n int) []WordCount {
	if n <= 0 {
		return nil
	}

	items := make([]WordCount, 0, len(counts))
	for word, count := range counts {
		items = append(items, WordCount{Word: word, Count: count})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Count != items[j].Count {
			return items[i].Count > items[j].Count
		}
		return items[i].Word < items[j].Word
	})

	if len(items) < n {
		return items
	}
	return items[:n]
}

func Tokenize(text string) []string {
	lower := strings.ToLower(text)
	var tokens []string
	var current []rune

	flush := func() {
		if len(current) == 0 {
			return
		}
		tokens = append(tokens, string(current))
		current = current[:0]
	}

	for _, r := range lower {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current = append(current, r)
			continue
		}
		flush()
	}
	flush()

	return tokens
}
