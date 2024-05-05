package util

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func CreateOutputDirectory(folderName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	dirPath := filepath.Join(cwd, "output", fmt.Sprintf("%v", folderName))
	err = os.MkdirAll(dirPath, os.ModePerm) // Use MkdirAll to create parent directories if they don't exist
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	return dirPath, nil
}

func CreateFile(dirPath string, fileName string, extension string) (*os.File, error) {
	if fileName == "" {
		return nil, fmt.Errorf("empty file name")
	}

	if extension == "" {
		return nil, fmt.Errorf("empty file extension")
	}

	filePath := filepath.Join(dirPath, fmt.Sprintf("%v.%v", fileName, extension))

	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func ReadCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file %v: %w", filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to parse file as CSV for %v: %w", filePath, err)
	}

	return data, nil
}
