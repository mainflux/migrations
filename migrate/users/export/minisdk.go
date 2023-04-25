package exportusers

import (
	"context"
	"encoding/csv"
	"fmt"
	"sync"

	"github.com/mainflux/mainflux/users/postgres" // for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0
	util "github.com/mainflux/migrations/internal"
	"golang.org/x/sync/errgroup"
)

var (
	defaultOffset     = uint64(0)
	retrievUsersOps   = "retrieving users"
	writeUsersOps     = "writing users to csv file"
	retrieveErrString = "error %v occurred at offset: %d and total: %d during %s"
	totalUsersQuery   = "SELECT COUNT(*) FROM users;"
)

// RetrieveAndWriteUsers retrieves existing users from the database and saves them to the provided csv file.
func RetrieveAndWriteUsers(ctx context.Context, query string, database postgres.Database, filePath string) error {
	out := make(chan []user)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		defer close(out)

		return retrieveUsers(ctx, query, database, out)
	})
	eg.Go(func() error {
		return usersToCSV(ctx, filePath, out)
	})

	return eg.Wait()
}

// retrieveUsers retrieves existing users from the database.
func retrieveUsers(ctx context.Context, query string, database postgres.Database, allUsers chan<- []user) error {
	totalUsers, err := total(ctx, database, totalUsersQuery, map[string]interface{}{})
	if err != nil {
		return err
	}
	offset := defaultOffset
	limit := util.UpdateLimit(totalUsers)

	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	for {
		wg.Add(1)
		go func(offset uint64, errCh chan<- error) {
			defer wg.Done()

			usersPage, err := dbRetrieveUsers(ctx, query, database, pageMetadata{Offset: offset, Limit: limit})
			if err != nil {
				errCh <- fmt.Errorf(retrieveErrString, err, offset, limit, retrievUsersOps)
			}
			select {
			case <-ctx.Done():
				return
			case allUsers <- usersPage.Users:
			}
			errCh <- nil
		}(offset, errCh)

		// Wait for the goroutine to finish or return an error.
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			if err != nil {
				return err
			}
		}

		if offset+limit >= totalUsers {
			break
		}
		offset += limit
	}

	wg.Wait()

	return nil
}

// usersToCSV saves users to the provided csv file
// The format of the users csv file is ID,Email,Password,Metadata.
func usersToCSV(ctx context.Context, filePath string, allUsers <-chan []user) error {
	file, err := util.CreateFile(filePath, writeUsersOps)
	if err != nil {
		return err
	}

	defer func() {
		if ferr := file.Close(); ferr != nil && err == nil {
			err = ferr
		}
	}()

	w := csv.NewWriter(file)

	records := [][]string{{"ID", "Email", "Password", "Metadata"}}
	for users := range allUsers {
		for _, user := range users {
			metadata, err := util.MetadataToString(user.Metadata)
			if err != nil {
				return err
			}
			record := []string{user.ID, user.Email, user.Password, metadata}
			records = append(records, record)
		}
	}

	return util.WriteData(ctx, w, file, records, writeUsersOps)
}
