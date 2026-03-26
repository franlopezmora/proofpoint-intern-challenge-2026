package normalize

import "testing"

func TestNormalizeSeriesName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "trim and collapse", input: "  game   OF thrones ", expected: "Game Of Thrones"},
		{name: "empty", input: "   ", expected: ""},
		{name: "control chars", input: "lost\x00\x1f", expected: "Lost"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := NormalizeSeriesName(tc.input)
			if got != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestNormalizeEpisodeTitle(t *testing.T) {
	got, _ := NormalizeEpisodeTitle(" ")
	if got != UntitledEpisode {
		t.Fatalf("expected default title %q, got %q", UntitledEpisode, got)
	}

	got, _ = NormalizeEpisodeTitle("  tHe  Long  night ")
	if got != "The Long Night" {
		t.Fatalf("expected normalized title, got %q", got)
	}
}

func TestParseNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{name: "valid", input: "12", expected: 12},
		{name: "negative", input: "-1", expected: 0},
		{name: "float", input: "3.5", expected: 0},
		{name: "word", input: "one", expected: 0},
		{name: "double sign", input: "--2", expected: 0},
		{name: "empty", input: "", expected: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotSeason, _ := ParseSeasonNumber(tc.input)
			gotEpisode, _ := ParseEpisodeNumber(tc.input)
			if gotSeason != tc.expected {
				t.Fatalf("season expected %d, got %d", tc.expected, gotSeason)
			}
			if gotEpisode != tc.expected {
				t.Fatalf("episode expected %d, got %d", tc.expected, gotEpisode)
			}
		})
	}
}

func TestParseAirDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "iso", input: "2022-11-03", expected: "2022-11-03"},
		{name: "slash", input: "2022/11/03", expected: "2022-11-03"},
		{name: "invalid text", input: "not a date", expected: UnknownAirDate},
		{name: "invalid date", input: "2022-40-99", expected: UnknownAirDate},
		{name: "zero date", input: "0000-00-00", expected: UnknownAirDate},
		{name: "empty", input: "", expected: UnknownAirDate},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := ParseAirDate(tc.input)
			if got != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}
