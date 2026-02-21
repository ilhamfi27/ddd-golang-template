package utils

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

func openCsvFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func ReadCsvFile(path string) ([][]string, error) {
	absPath, _ := filepath.Abs(path)
	file, err := openCsvFile(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
