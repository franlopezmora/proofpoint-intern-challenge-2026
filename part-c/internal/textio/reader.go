package textio

import (
	"fmt"
	"os"
)

func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read input file: %w", err)
	}
	return string(content), nil
}
