package csvio

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

const ExpectedColumns = 5

func ReadRows(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open input CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = false

	var rows [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read CSV row: %w", err)
		}
		rows = append(rows, normalizeRow(record))
	}

	if len(rows) > 0 && looksLikeHeader(rows[0]) {
		rows = rows[1:]
	}

	return rows, nil
}

func normalizeRow(record []string) []string {
	normalized := make([]string, ExpectedColumns)

	switch {
	case len(record) == ExpectedColumns:
		copy(normalized, record)
	case len(record) < ExpectedColumns:
		copy(normalized, record)
	case len(record) > ExpectedColumns:
		normalized[0] = record[0]
		normalized[1] = record[1]
		normalized[2] = record[2]
		normalized[3] = strings.Join(record[3:len(record)-1], ",")
		normalized[4] = record[len(record)-1]
	}

	return normalized
}

func looksLikeHeader(record []string) bool {
	if len(record) < ExpectedColumns {
		return false
	}
	values := make([]string, 0, ExpectedColumns)
	for i := 0; i < ExpectedColumns; i++ {
		values = append(values, strings.ToLower(strings.TrimSpace(record[i])))
	}

	expected := []string{"series name", "season number", "episode number", "episode title", "air date"}
	for i := range expected {
		if values[i] != expected[i] {
			return false
		}
	}
	return true
}
