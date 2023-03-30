package users13

import (
	"context"
	"encoding/csv"
	"fmt"

	mf13users "github.com/mainflux/mainflux/users"
	mf13postgres "github.com/mainflux/mainflux/users/postgres"
	"github.com/mainflux/migrations/internal/util"
)

var (
	offset            = uint64(0)
	limit             = uint64(100)
	retrievUsersOps   = "retrieving users"
	writeUsersOps     = "writing users to csv file"
	retrieveErrString = "error %v occured at offset: %d and total: %d during %s"
)

// RetrieveUsers retrieves existing users from the database
func RetrieveUsers(ctx context.Context, db mf13postgres.Database) (mf13users.UserPage, error) {
	usersPage, err := dbRetrieveUsers(ctx, db, mf13users.PageMetadata{Offset: offset, Limit: limit})
	if err != nil {
		return mf13users.UserPage{}, fmt.Errorf(retrieveErrString, err, offset, limit, retrievUsersOps)
	}
	o := limit
	limit = util.UpdateLimit(usersPage.Total)
	for o < usersPage.Total {
		ths, err := dbRetrieveUsers(ctx, db, mf13users.PageMetadata{Offset: o, Limit: limit})
		if err != nil {
			return mf13users.UserPage{}, fmt.Errorf(retrieveErrString, err, o, limit, retrievUsersOps)
		}
		usersPage.Users = append(usersPage.Users, ths.Users...)
		o = o + limit
	}

	return usersPage, nil
}

// UsersToCSV saves users to the provided csv file
// The format of the users csv file is ID,Email,Password,Metadata
func UsersToCSV(filePath string, users []mf13users.User) error {
	f, err := util.CreateFile(filePath, writeUsersOps)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)

	records := [][]string{{"ID", "Email", "Password", "Metadata"}}
	for _, user := range users {
		metadata, err := util.MetadataToString(user.Metadata)
		if err != nil {
			return err
		}
		record := []string{user.ID, user.Email, user.Password, metadata}
		records = append(records, record)
	}

	return util.WriteData(w, f, records, writeUsersOps)
}
