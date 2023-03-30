package things14

import (
	"fmt"

	mf14sdk "github.com/mainflux/mainflux/pkg/sdk/go/0140"
	"github.com/mainflux/migrations/internal/util"
	users "github.com/mainflux/migrations/migrate/users/v14"
)

var (
	limit         = uint64(100)
	readErrString = "error %v occured during %s"
)

// CreateThings creates things from the provided csv file
// The format of the things csv file is ID,Key,Name,Owner,Metadata
func CreateThings(sdk mf14sdk.SDK, usersPath, filePath, token string) error {
	records, err := util.ReadData(filePath)
	if err != nil {
		return fmt.Errorf(readErrString, err, "creating things")
	}
	ths := []mf14sdk.Thing{}
	for _, record := range records {
		metadata, err := util.MetadataFromString(record[4])
		if err != nil {
			return err
		}
		thing := mf14sdk.Thing{
			ID:    record[0],
			Name:  record[2],
			Owner: users.GetUserID(usersPath, record[3]),
			Credentials: mf14sdk.Credentials{
				Secret: record[1],
			},
			Metadata: metadata,
			Status:   mf14sdk.EnabledStatus,
		}
		ths = append(ths, thing)
		if len(ths) <= int(limit) {
			if _, err := sdk.CreateThings(ths, token); err != nil {
				return fmt.Errorf("failed to create things with error %v", err)
			}
			ths = []mf14sdk.Thing{}
		}
	}

	return nil
}

// CreateChannels creates channels from the provided csv file
// The format of the channels csv file is ID,Name,Owner,Metadata
func CreateChannels(sdk mf14sdk.SDK, usersPath, filePath, token string) error {
	records, err := util.ReadData(filePath)
	if err != nil {
		return fmt.Errorf(readErrString, err, "creating channels")
	}
	chs := []mf14sdk.Channel{}
	for _, record := range records {
		metadata, err := util.MetadataFromString(record[3])
		if err != nil {
			return err
		}
		channel := mf14sdk.Channel{
			ID:       record[0],
			Name:     record[1],
			OwnerID:  users.GetUserID(usersPath, record[2]),
			Metadata: metadata,
			Status:   mf14sdk.EnabledStatus,
		}
		chs = append(chs, channel)
		if len(chs) <= int(limit) {
			if _, err := sdk.CreateChannels(chs, token); err != nil {
				return fmt.Errorf("failed to create channel with error %v", err)
			}
			chs = []mf14sdk.Channel{}
		}
	}

	return nil
}

// CreateConnections creates policies for things to read and write to the
// specified channels. The format of the connections csv file is
// ChannelID,ChannelOwner,ThingID,ThingOwner
func CreateConnections(sdk mf14sdk.SDK, filePath string, token string) error {
	records, err := util.ReadData(filePath)
	if err != nil {
		return fmt.Errorf(readErrString, err, "creating connections")
	}

	for _, record := range records {
		if err := sdk.ConnectThing(record[2], record[0], token); err != nil {
			return fmt.Errorf("failed to connect thing with id %s to channel with id %s with error %v", record[2], record[0], err)
		}
	}

	return nil
}
