package util

import (
	"encoding/csv"
	"fmt"
	"os"
)

func WriteCsvFile(data [][]string, filepath string) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	err = writer.WriteAll(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file %s: %w", filepath, err)
	}

	return nil
}
