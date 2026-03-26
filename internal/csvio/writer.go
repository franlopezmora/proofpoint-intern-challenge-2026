package csvio

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"proofpoint-flm/internal/domain"
)

func WriteCleanCSV(outDir string, episodes []domain.Episode) (string, error) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("create output directory: %w", err)
	}

	path := filepath.Join(outDir, "episodes_clean.csv")
	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create clean CSV: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"SeriesName", "SeasonNumber", "EpisodeNumber", "EpisodeTitle", "AirDate"}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("write CSV header: %w", err)
	}

	for _, ep := range episodes {
		row := []string{
			ep.SeriesName,
			strconv.Itoa(ep.SeasonNumber),
			strconv.Itoa(ep.EpisodeNumber),
			ep.EpisodeTitle,
			ep.AirDate,
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("write CSV row: %w", err)
		}
	}

	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("flush CSV writer: %w", err)
	}

	return path, nil
}
