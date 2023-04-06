package users13

import (
	"context"
	"encoding/csv"
	"fmt"
	"sync"

	mf13users "github.com/mainflux/mainflux/users"
	mf13postgres "github.com/mainflux/mainflux/users/postgres"
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

func RetrieveAndWriteUsers(ctx context.Context, database mf13postgres.Database, filePath string) error {
	out := make(chan []mf13users.User)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		defer close(out)

		return RetrieveUsers(ctx, database, out)
	})
	eg.Go(func() error {
		return UsersToCSV(ctx, filePath, out)
	})

	return eg.Wait()
}

// RetrieveUsers retrieves existing users from the database.
func RetrieveUsers(ctx context.Context, database mf13postgres.Database, allUsers chan<- []mf13users.User) error {
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

			usersPage, err := dbRetrieveUsers(ctx, database, mf13users.PageMetadata{Offset: offset, Limit: limit})
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

// UsersToCSV saves users to the provided csv file
// The format of the users csv file is ID,Email,Password,Metadata.
func UsersToCSV(ctx context.Context, filePath string, allUsers <-chan []mf13users.User) error {
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
