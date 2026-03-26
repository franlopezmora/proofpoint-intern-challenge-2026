package dedupe

import (
	"testing"

	"proofpoint-flm/internal/domain"
	"proofpoint-flm/internal/normalize"
)

func TestCandidateKeys(t *testing.T) {
	ep := domain.Episode{SeriesName: "The Office", SeasonNumber: 0, EpisodeNumber: 3, EpisodeTitle: "Diversity Day"}
	keys := candidateKeys(ep)
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "A|the office|0|3" {
		t.Fatalf("unexpected A key: %s", keys[0])
	}
	if keys[1] != "B|the office|0|3|diversity day" {
		t.Fatalf("unexpected B key: %s", keys[1])
	}
}

func TestBetterPriority(t *testing.T) {
	base := domain.Episode{SeriesName: "X", SeasonNumber: 0, EpisodeNumber: 2, EpisodeTitle: normalize.UntitledEpisode, AirDate: normalize.UnknownAirDate, InputOrder: 0}
	withDate := base
	withDate.AirDate = "2024-01-01"
	if !better(withDate, base) {
		t.Fatalf("record with known air date should be better")
	}

	withTitle := base
	withTitle.EpisodeTitle = "Pilot"
	if !better(withTitle, base) {
		t.Fatalf("record with known title should be better")
	}

	withNumbers := base
	withNumbers.SeasonNumber = 1
	withNumbers.EpisodeNumber = 1
	if !better(withNumbers, base) {
		t.Fatalf("record with known numbers should be better")
	}
}

func TestApplyDeduplicatesAndKeepsBest(t *testing.T) {
	episodes := []domain.Episode{
		{SeriesName: "The Office", SeasonNumber: 1, EpisodeNumber: 1, EpisodeTitle: normalize.UntitledEpisode, AirDate: normalize.UnknownAirDate, InputOrder: 0},
		{SeriesName: "The Office", SeasonNumber: 1, EpisodeNumber: 1, EpisodeTitle: "Pilot", AirDate: "2005-03-24", InputOrder: 1},
		{SeriesName: "Dark", SeasonNumber: 1, EpisodeNumber: 1, EpisodeTitle: "Secrets", AirDate: "2017-12-01", InputOrder: 2},
	}

	result := Apply(episodes)
	if result.DuplicatesDetected != 1 {
		t.Fatalf("expected 1 duplicate, got %d", result.DuplicatesDetected)
	}
	if len(result.Episodes) != 2 {
		t.Fatalf("expected 2 output episodes, got %d", len(result.Episodes))
	}

	for _, ep := range result.Episodes {
		if ep.SeriesName == "The Office" {
			if ep.EpisodeTitle != "Pilot" || ep.AirDate != "2005-03-24" {
				t.Fatalf("best record not selected for The Office")
			}
		}
	}
}
