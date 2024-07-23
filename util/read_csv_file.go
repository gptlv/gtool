package util

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ReadCsvFile(filePath string) ([][]string, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file %v: %w", filePath, err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to parse file as CSV for %v: %w", filePath, err)
	}

	return data, nil

}
