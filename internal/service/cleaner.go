package service

import (
	"fmt"

	"proofpoint-flm/internal/csvio"
	"proofpoint-flm/internal/dedupe"
	"proofpoint-flm/internal/domain"
	"proofpoint-flm/internal/normalize"
	"proofpoint-flm/internal/report"
)

type Cleaner struct{}

type RunResult struct {
	CleanCSVPath string
	ReportPath   string
	Metrics      report.Metrics
}

func NewCleaner() Cleaner {
	return Cleaner{}
}

func (c Cleaner) Run(inputPath, outDir string) (RunResult, error) {
	rows, err := csvio.ReadRows(inputPath)
	if err != nil {
		return RunResult{}, err
	}

	cleaned := make([]domain.Episode, 0, len(rows))
	corrected := 0
	discarded := 0

	for index, row := range rows {
		outcome := sanitizeRow(row, index)
		if outcome.Corrected {
			corrected++
		}
		if outcome.Discarded {
			discarded++
			continue
		}
		cleaned = append(cleaned, outcome.Episode)
	}

	deduped := dedupe.Apply(cleaned)

	cleanPath, err := csvio.WriteCleanCSV(outDir, deduped.Episodes)
	if err != nil {
		return RunResult{}, err
	}

	metrics := report.Metrics{
		TotalInputRecords:  len(rows),
		TotalOutputRecords: len(deduped.Episodes),
		DiscardedEntries:   discarded,
		CorrectedEntries:   corrected,
		DuplicatesDetected: deduped.DuplicatesDetected,
		DeduplicationNotes: deduped.StrategyDescription,
	}

	reportPath, err := report.Write(outDir, metrics)
	if err != nil {
		return RunResult{}, err
	}

	return RunResult{
		CleanCSVPath: cleanPath,
		ReportPath:   reportPath,
		Metrics:      metrics,
	}, nil
}

func sanitizeRow(row []string, inputOrder int) domain.RecordOutcome {
	seriesRaw := safeAt(row, 0)
	seasonRaw := safeAt(row, 1)
	episodeRaw := safeAt(row, 2)
	titleRaw := safeAt(row, 3)
	airDateRaw := safeAt(row, 4)

	seriesName, seriesChanged := normalize.NormalizeSeriesName(seriesRaw)
	if seriesName == "" {
		return domain.RecordOutcome{Discarded: true, Corrected: seriesChanged || seriesRaw == ""}
	}

	seasonNumber, seasonChanged := normalize.ParseSeasonNumber(seasonRaw)
	episodeNumber, episodeChanged := normalize.ParseEpisodeNumber(episodeRaw)
	episodeTitle, titleChanged := normalize.NormalizeEpisodeTitle(titleRaw)
	airDate, airDateChanged := normalize.ParseAirDate(airDateRaw)

	isMissingInfo := episodeNumber == 0 && !normalize.IsTitleKnown(episodeTitle) && !normalize.IsAirDateKnown(airDate)
	if isMissingInfo {
		return domain.RecordOutcome{
			Discarded: true,
			Corrected: seriesChanged || seasonChanged || episodeChanged || titleChanged || airDateChanged,
		}
	}

	return domain.RecordOutcome{
		Episode: domain.Episode{
			SeriesName:    seriesName,
			SeasonNumber:  seasonNumber,
			EpisodeNumber: episodeNumber,
			EpisodeTitle:  episodeTitle,
			AirDate:       airDate,
			InputOrder:    inputOrder,
		},
		Corrected: seriesChanged || seasonChanged || episodeChanged || titleChanged || airDateChanged,
	}
}

func safeAt(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return row[idx]
}

func ValidateArgs(input, outDir string) error {
	if input == "" {
		return fmt.Errorf("missing required argument: --input")
	}
	if outDir == "" {
		return fmt.Errorf("missing required argument: --outdir")
	}
	return nil
}
