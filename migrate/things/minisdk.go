package things

import (
	"context"
	"encoding/csv"
	"fmt"

	mf14sdk "github.com/mainflux/mainflux/pkg/sdk/go/0140"
	mf13things "github.com/mainflux/mainflux/things"
	mf13postgres "github.com/mainflux/mainflux/things/postgres"
	"github.com/mainflux/migrations/internal/util"
	"github.com/mainflux/migrations/migrate/users"
)

var (
	offset              = uint64(0)
	limit               = uint64(100)
	retrievThingsOps    = "retrieving things"
	retrievChannelsOps  = "retrieving channels"
	retrievConnOps      = "retrieving connections"
	writeThingsOps      = "writing things to csv file"
	writeChannelsOps    = "writing channels to csv file"
	writeConnectionsOps = "writing connections to csv file"
	retrieveErrString   = "error %v occured at offset: %d and total: %d during %s"
	readErrString       = "error %v occured during %s"
)

// MFRetrieveThings retrieves existing things from the database
func MFRetrieveThings(ctx context.Context, db mf13postgres.Database) (mf13things.Page, error) {
	thingsPage, err := RetrieveAllThings(ctx, db, mf13things.PageMetadata{Offset: offset, Limit: limit})
	if err != nil {
		return mf13things.Page{}, fmt.Errorf(retrieveErrString, err, offset, limit, retrievThingsOps)
	}
	o := limit
	limit = util.UpdateLimit(thingsPage.Total)
	for o < thingsPage.Total {
		ths, err := RetrieveAllThings(ctx, db, mf13things.PageMetadata{Offset: o, Limit: limit})
		if err != nil {
			return mf13things.Page{}, fmt.Errorf(retrieveErrString, err, o, limit, retrievThingsOps)
		}
		thingsPage.Things = append(thingsPage.Things, ths.Things...)
		o = o + limit
	}

	return thingsPage, nil
}

// MFRetrieveChannels retrieves existing channels from the database
func MFRetrieveChannels(ctx context.Context, db mf13postgres.Database) (mf13things.ChannelsPage, error) {
	channelsPage, err := RetrieveAllChannels(ctx, db, mf13things.PageMetadata{Offset: offset, Limit: limit})
	if err != nil {
		return mf13things.ChannelsPage{}, fmt.Errorf(retrieveErrString, err, offset, limit, retrievChannelsOps)
	}
	o := limit
	limit = util.UpdateLimit(channelsPage.Total)
	for o < channelsPage.Total {
		chs, err := RetrieveAllChannels(ctx, db, mf13things.PageMetadata{Offset: o, Limit: limit})
		if err != nil {
			return mf13things.ChannelsPage{}, fmt.Errorf(retrieveErrString, err, o, limit, retrievChannelsOps)
		}
		channelsPage.Channels = append(channelsPage.Channels, chs.Channels...)
		o = o + limit
	}

	return channelsPage, nil
}

// MFRetrieveConnections retrieves existing things to channels connection from the database
func MFRetrieveConnections(ctx context.Context, db mf13postgres.Database) (ConnectionsPage, error) {
	connectionsPage, err := RetrieveAllConnections(ctx, db, mf13things.PageMetadata{Offset: offset, Limit: limit})
	if err != nil {
		return ConnectionsPage{}, fmt.Errorf(retrieveErrString, err, offset, limit, retrievConnOps)
	}
	o := limit
	limit = util.UpdateLimit(connectionsPage.Total)
	for o < connectionsPage.Total {
		conns, err := RetrieveAllConnections(ctx, db, mf13things.PageMetadata{Offset: o, Limit: limit})
		if err != nil {
			return ConnectionsPage{}, fmt.Errorf(retrieveErrString, err, o, limit, retrievConnOps)
		}
		connectionsPage.Connections = append(connectionsPage.Connections, conns.Connections...)
		o = o + limit
	}

	return connectionsPage, nil
}

// ThingsToCSV saves things to the provided csv file
// The format of the things csv file is ID,Key,Name,Owner,Metadata
func ThingsToCSV(filePath string, things []mf13things.Thing) error {
	f, err := util.CreateFile(filePath, writeThingsOps)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)

	records := [][]string{{"ID", "Key", "Name", "Owner", "Metadata"}}
	for _, thing := range things {
		metadata, err := util.MetadataToString(thing.Metadata)
		if err != nil {
			return err
		}
		record := []string{thing.ID, thing.Key, thing.Name, thing.Owner, metadata}
		records = append(records, record)
	}

	return util.WriteData(w, f, records, writeThingsOps)
}

// ChannelsToCSV saves channels to the provided csv file
// The format of the channels csv file is ID,Name,Owner,Metadata
func ChannelsToCSV(filePath string, channels []mf13things.Channel) error {
	f, err := util.CreateFile(filePath, writeChannelsOps)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)

	records := [][]string{{"ID", "Name", "Owner", "Metadata"}}
	for _, channel := range channels {
		metadata, err := util.MetadataToString(channel.Metadata)
		if err != nil {
			return err
		}
		record := []string{channel.ID, channel.Name, channel.Owner, metadata}
		records = append(records, record)
	}

	return util.WriteData(w, f, records, writeChannelsOps)
}

// ConnectionsToCSV saves connections to the provided csv file
// The format of the connections csv file is ChannelID,ChannelOwner,ThingID,ThingOwner
func ConnectionsToCSV(filePath string, connections []Connection) error {
	f, err := util.CreateFile(filePath, writeConnectionsOps)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)

	records := [][]string{{"ChannelID", "ChannelOwner", "ThingID", "ThingOwner"}}
	for _, conn := range connections {
		record := []string{conn.ChannelID, conn.ChannelOwner, conn.ThingID, conn.ThingOwner}
		records = append(records, record)
	}

	return util.WriteData(w, f, records, writeConnectionsOps)
}

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
