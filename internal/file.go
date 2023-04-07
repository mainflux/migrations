package util

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mainflux/mainflux/pkg/errors"
)

const (
	openOp        = "open file"
	writeOp       = "write to file"
	closeOp       = "close file"
	dirOp         = "create directory"
	fileOp        = "create file"
	fileErrString = "failed to %s %s during %s with error : %w"
	filePerm      = 0700
	batchSize     = 100
)

// createrDirectory creates a directory if it doesn't already exist.
func createrDirectory(path string) error {
	var dir = filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, filePerm); err != nil {
			return err
		}
	}

	return nil
}

// CreateFile creates a file if it doesn't already exist and opens it.
func CreateFile(filePath, operation string) (*os.File, error) {
	if err := createrDirectory(filePath); err != nil {
		return nil, fmt.Errorf(fileErrString, dirOp, filePath, operation, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf(fileErrString, fileOp, filePath, operation, err)
	}

	return file, nil
}

// ReadInBatch reads data from from the provided csv file in batches.
func ReadInBatch(ctx context.Context, filePath, operation string, outth chan<- []string) error {
	defer close(outth)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf(fileErrString, openOp, filePath, operation, err)
	}

	defer func() {
		if ferr := file.Close(); ferr != nil && err == nil {
			err = fmt.Errorf(fileErrString, closeOp, file.Name(), operation, ferr)
		}
	}()

	reader := csv.NewReader(file)

	// skip first line
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read csv header from file %s during %s: %w", filePath, operation, err)
	}

	// use a buffered channel to reduce overhead of sending records
	recordCh := make(chan []string, batchSize)
	errCh := make(chan error, 1)

	// use a goroutine to read from the file and send records to the channel
	go func(errCh chan<- error) {
		defer close(recordCh)
		for {
			record, err := reader.Read()
			if errors.Contains(err, io.EOF) {
				errCh <- nil
				break
			}
			if err != nil {
				errCh <- fmt.Errorf("failed to read csv data from file %s during %s: %w", filePath, operation, err)

				break
			}
			recordCh <- record
		}
	}(errCh)

	// read records from the channel and send to output thread
	for record := range recordCh {
		outth <- record
	}

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadAllData reads data from from the provided csv file.
func ReadAllData(fileName, operation string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return [][]string{}, err
	}

	defer func() {
		if ferr := file.Close(); ferr != nil && err == nil {
			err = fmt.Errorf(fileErrString, closeOp, file.Name(), operation, ferr)
		}
	}()

	reader := csv.NewReader(file)

	// skip first line
	if _, err := reader.Read(); err != nil {
		return [][]string{}, fmt.Errorf("failed to read csv header from file %s: %w", fileName, err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return [][]string{}, fmt.Errorf("failed to read csv data from file %s: %w", fileName, err)
	}

	return records, nil
}

// WriteData writes data to the provided csv file.
func WriteData(ctx context.Context, writer *csv.Writer, file *os.File, records [][]string, operation string) error {
	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf(fileErrString, writeOp, file.Name(), operation, err)
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if writer.Error() != nil {
		return fmt.Errorf("writer error on file %s during %s with error: %w", file.Name(), operation, writer.Error())
	}

	return nil
}
