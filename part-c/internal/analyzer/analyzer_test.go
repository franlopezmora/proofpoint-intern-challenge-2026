package analyzer

import (
	"reflect"
	"testing"
)

func TestTokenizeIgnoresPunctuationAndSpecialChars(t *testing.T) {
	input := `Hello, world! (Go-lang) isn't_bad? go123...`
	got := Tokenize(input)
	want := []string{"hello", "world", "go", "lang", "isn", "t", "bad", "go123"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected tokens\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestCountWordsCaseInsensitive(t *testing.T) {
	counts := CountWords("Data data DATA dAtA")
	if counts["data"] != 4 {
		t.Fatalf("expected data=4, got %d", counts["data"])
	}
	if len(counts) != 1 {
		t.Fatalf("expected 1 unique word, got %d", len(counts))
	}
}

func TestTopNSortsByFrequencyThenAlphabetically(t *testing.T) {
	counts := map[string]int{
		"banana": 3,
		"apple":  3,
		"carrot": 2,
		"date":   2,
	}

	got := TopN(counts, 10)
	want := []WordCount{
		{Word: "apple", Count: 3},
		{Word: "banana", Count: 3},
		{Word: "carrot", Count: 2},
		{Word: "date", Count: 2},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected ranking\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestTopNReturnsOnlyAvailableWords(t *testing.T) {
	counts := map[string]int{"a": 2, "b": 1}
	got := TopN(counts, 10)
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
}
