package users13

import (
	"context"
	"encoding/csv"
	"fmt"

	mf13users "github.com/mainflux/mainflux/users"
	mf13postgres "github.com/mainflux/mainflux/users/postgres"
	util "github.com/mainflux/migrations/internal"
	"golang.org/x/sync/errgroup"
)

var (
	offset            = uint64(0)
	limit             = uint64(100)
	retrievUsersOps   = "retrieving users"
	writeUsersOps     = "writing users to csv file"
	retrieveErrString = "error %v occured at offset: %d and total: %d during %s"
)

func RetrieveAndWriteUsers(ctx context.Context, db mf13postgres.Database, filePath string) error {
	out := make(chan []mf13users.User)

	g := new(errgroup.Group)

	g.Go(func() error {
		return RetrieveUsers(ctx, db, out)
	})
	g.Go(func() error {
		return UsersToCSV(filePath, out)
	})

	return g.Wait()
}

// RetrieveUsers retrieves existing users from the database
func RetrieveUsers(ctx context.Context, db mf13postgres.Database, outusers chan<- []mf13users.User) error {
	o := offset
	l := limit
	for {
		usersPage, err := dbRetrieveUsers(ctx, db, mf13users.PageMetadata{Offset: o, Limit: l})
		if err != nil {
			return fmt.Errorf(retrieveErrString, err, o, l, retrievUsersOps)
		}
		outusers <- usersPage.Users
		if o+l >= usersPage.Total {
			break
		}

		o = o + l
		l = util.UpdateLimit(usersPage.Total)
	}
	close(outusers)
	return nil
}

// UsersToCSV saves users to the provided csv file
// The format of the users csv file is ID,Email,Password,Metadata
func UsersToCSV(filePath string, inusers <-chan []mf13users.User) error {
	f, err := util.CreateFile(filePath, writeUsersOps)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)

	records := [][]string{{"ID", "Email", "Password", "Metadata"}}
	for users := range inusers {
		for _, user := range users {
			metadata, err := util.MetadataToString(user.Metadata)
			if err != nil {
				return err
			}
			record := []string{user.ID, user.Email, user.Password, metadata}
			records = append(records, record)
		}
	}
	return util.WriteData(w, f, records, writeUsersOps)
}
