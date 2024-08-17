package util

import (
	"fmt"
	"os"
)

func ReadInputFile(name string) (*os.File, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}

	return file, nil
}
