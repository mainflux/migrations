package things13

import (
	"context"
	"encoding/csv"
	"fmt"

	mf13things "github.com/mainflux/mainflux/things"
	mf13postgres "github.com/mainflux/mainflux/things/postgres"
	"github.com/mainflux/migrations/internal/util"
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
)

// RetrieveThings retrieves existing things from the database
func RetrieveThings(ctx context.Context, db mf13postgres.Database) (mf13things.Page, error) {
	thingsPage, err := dbRetrieveThings(ctx, db, mf13things.PageMetadata{Offset: offset, Limit: limit})
	if err != nil {
		return mf13things.Page{}, fmt.Errorf(retrieveErrString, err, offset, limit, retrievThingsOps)
	}
	o := limit
	limit = util.UpdateLimit(thingsPage.Total)
	for o < thingsPage.Total {
		ths, err := dbRetrieveThings(ctx, db, mf13things.PageMetadata{Offset: o, Limit: limit})
		if err != nil {
			return mf13things.Page{}, fmt.Errorf(retrieveErrString, err, o, limit, retrievThingsOps)
		}
		thingsPage.Things = append(thingsPage.Things, ths.Things...)
		o = o + limit
	}

	return thingsPage, nil
}

// RetrieveChannels retrieves existing channels from the database
func RetrieveChannels(ctx context.Context, db mf13postgres.Database) (mf13things.ChannelsPage, error) {
	channelsPage, err := dbRetrieveChannels(ctx, db, mf13things.PageMetadata{Offset: offset, Limit: limit})
	if err != nil {
		return mf13things.ChannelsPage{}, fmt.Errorf(retrieveErrString, err, offset, limit, retrievChannelsOps)
	}
	o := limit
	limit = util.UpdateLimit(channelsPage.Total)
	for o < channelsPage.Total {
		chs, err := dbRetrieveChannels(ctx, db, mf13things.PageMetadata{Offset: o, Limit: limit})
		if err != nil {
			return mf13things.ChannelsPage{}, fmt.Errorf(retrieveErrString, err, o, limit, retrievChannelsOps)
		}
		channelsPage.Channels = append(channelsPage.Channels, chs.Channels...)
		o = o + limit
	}

	return channelsPage, nil
}

// RetrieveConnections retrieves existing things to channels connection from the database
func RetrieveConnections(ctx context.Context, db mf13postgres.Database) (ConnectionsPage, error) {
	connectionsPage, err := dbRetrieveConnections(ctx, db, mf13things.PageMetadata{Offset: offset, Limit: limit})
	if err != nil {
		return ConnectionsPage{}, fmt.Errorf(retrieveErrString, err, offset, limit, retrievConnOps)
	}
	o := limit
	limit = util.UpdateLimit(connectionsPage.Total)
	for o < connectionsPage.Total {
		conns, err := dbRetrieveConnections(ctx, db, mf13things.PageMetadata{Offset: o, Limit: limit})
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
