package service

import (
	"testing"

	"proofpoint-flm/internal/normalize"
)

func TestSanitizeRowDiscardWhenSeriesMissing(t *testing.T) {
	row := []string{"  ", "1", "1", "Pilot", "2020-01-01"}
	outcome := sanitizeRow(row, 0)
	if !outcome.Discarded {
		t.Fatalf("expected row to be discarded when series is missing")
	}
}

func TestSanitizeRowDiscardWhenEpisodeTitleAndAirDateAllMissing(t *testing.T) {
	row := []string{"Lost", "1", "", "", ""}
	outcome := sanitizeRow(row, 0)
	if !outcome.Discarded {
		t.Fatalf("expected row to be discarded when episode number, title and air date are missing")
	}
}

func TestSanitizeRowAppliesDefaults(t *testing.T) {
	row := []string{"  the office ", "one", "-2", "pilot", "not a date"}
	outcome := sanitizeRow(row, 4)
	if outcome.Discarded {
		t.Fatalf("did not expect discard")
	}
	ep := outcome.Episode
	if ep.SeriesName != "The Office" {
		t.Fatalf("expected normalized series name, got %q", ep.SeriesName)
	}
	if ep.SeasonNumber != 0 || ep.EpisodeNumber != 0 {
		t.Fatalf("expected defaults for season and episode numbers")
	}
	if ep.EpisodeTitle != "Pilot" {
		t.Fatalf("expected normalized title")
	}
	if ep.AirDate != normalize.UnknownAirDate {
		t.Fatalf("expected default air date")
	}
	if !outcome.Corrected {
		t.Fatalf("expected corrected flag")
	}
}
