package util

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	batchSize         = uint8(100)
	retrievUsersOps   = "retrieving users"
	writeUsersOps     = "writing users to csv file"
	openOp            = "open file"
	writeOp           = "write to file"
	closeOp           = "close file"
	dirOp             = "create directory"
	fileOp            = "create file"
	retrieveErrString = "error %v occured at offset: %d and total: %d during %s"
	fileErrString     = "failed to %s %s during %s with error : %v "
)

// createrDirectory creates a directory if it doesn't already exist
func createrDirectory(path string) error {
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
	if err := createrDirectory(filePath); err != nil {
		return nil, fmt.Errorf(fileErrString, dirOp, filePath, operation, err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf(fileErrString, fileOp, filePath, operation, err)
	}
	return f, nil
}

// ReadInBatch reads data from from the provided csv file in batches
func ReadInBatch(filePath, operation string, outth chan<- []string) error {
	defer close(outth)

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf(fileErrString, openOp, filePath, operation, err)
	}

	defer func() {
		if ferr := f.Close(); ferr != nil && err == nil {
			err = fmt.Errorf(fileErrString, closeOp, f.Name(), operation, ferr)
		}
	}()

	reader := csv.NewReader(f)

	// skip first line
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read csv header from file %s during %s: %v", filePath, operation, err)
	}

	// use a buffered channel to reduce overhead of sending records
	recordCh := make(chan []string, 100)
	errCh := make(chan error)

	// use a goroutine to read from the file and send records to the channel
	go func(errCh chan<- error) {
		defer close(recordCh)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errCh <- fmt.Errorf("failed to read csv data from file %s during %s: %v", filePath, operation, err)
				break
			}
			recordCh <- record
		}
	}(errCh)

	// read records from the channel and send to output thread
	for record := range recordCh {
		outth <- record
	}

	if err := <-errCh; err != nil {
		return err
	}

	return nil
}

// ReadAllData reads data from from the provided csv file
func ReadAllData(fileName, operation string) ([][]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return [][]string{}, err
	}

	defer func() {
		if ferr := f.Close(); ferr != nil && err == nil {
			err = fmt.Errorf(fileErrString, closeOp, f.Name(), operation, ferr)
		}
	}()

	reader := csv.NewReader(f)

	// skip first line
	if _, err := reader.Read(); err != nil {
		return [][]string{}, fmt.Errorf("failed to read csv header from file %s: %v", fileName, err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return [][]string{}, fmt.Errorf("failed to read csv data from file %s: %v", fileName, err)
	}

	return records, nil
}

// WriteData writes data to the provided csv file
func WriteData(ctx context.Context, writer *csv.Writer, file *os.File, records [][]string, operation string) error {
	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf(fileErrString, writeOp, file.Name(), operation, err)
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if writer.Error() != nil {
		return fmt.Errorf("writer error on file %s during %s with error: %v", file.Name(), operation, writer.Error())
	}

	return nil
}
