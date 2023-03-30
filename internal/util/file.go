package util

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

const (
	retrievUsersOps   = "retrieving users"
	writeUsersOps     = "writing users to csv file"
	writeOp           = "write to file"
	dirOp             = "create directory"
	fileOp            = "create file"
	closeOp           = "close file"
	retrieveErrString = "error %v occured at offset: %d and total: %d during %s"
	fileErrString     = "failed to %s with error %v during %s"
)

// CreaterDir creates a directory if it doesn't already exist
func CreaterDir(path string) error {
	var dir string = filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	return nil
}

// CreateFile creates a file if it doesn't already exist and opens it
func CreateFile(filePath, operation string) (*os.File, error) {
	if err := CreaterDir(filePath); err != nil {
		return nil, fmt.Errorf(fileErrString, dirOp, err, operation)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf(fileErrString, fileOp, err, operation)
	}
	return f, nil
}

// ReadAllData reads data from from the provided csv file
func ReadAllData(fileName string) ([][]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return [][]string{}, err
	}

	reader := csv.NewReader(f)

	// skip first line
	if _, err := reader.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	if err := f.Close(); err != nil {
		return [][]string{}, err
	}
	return records, nil
}

func WriteData(writer *csv.Writer, file *os.File, records [][]string, operation string) error {
	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf(fileErrString, writeOp, err, operation)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf(fileErrString, closeOp, err, operation)
	}
	return nil
}
