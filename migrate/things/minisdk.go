package things

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	mf14sdk "github.com/mainflux/mainflux/pkg/sdk/go/0140"
	mf13things "github.com/mainflux/mainflux/things"
	mf13postgres "github.com/mainflux/mainflux/things/postgres"
	"github.com/mainflux/migrations/migrate/users"
)

const (
	offset = uint64(0)
	limit  = uint64(100)
)

var (
	retrievThingsOps    = "retrieving things"
	retrievChannelsOps  = "retrieving channels"
	retrievConnOps      = "retrieving connections"
	writeThingsOps      = "writing things to csv file"
	writeChannelsOps    = "writing channels to csv file"
	writeConnectionsOps = "writing connections to csv file"
	writeOp             = "write to file"
	dirOp               = "create directory"
	fileOp              = "create file"
	closeOp             = "close file"
	retrieveErrString   = "error %v occured at offset: %d and total: %d during %s"
	fileErrString       = "failed to %s with error %v during %s"
	readErrString       = "error %v occured during %s"
)

// MFRetrieveThings retrieves existing things from the database
func MFRetrieveThings(ctx context.Context, db mf13postgres.Database) (mf13things.Page, error) {
	thingsPage, err := RetrieveAllThings(ctx, db, mf13things.PageMetadata{Offset: offset, Limit: limit})
	if err != nil {
		return mf13things.Page{}, fmt.Errorf(retrieveErrString, err, offset, limit, retrievThingsOps)
	}
	o := uint64(100)
	for o < thingsPage.Total {
		ths, err := RetrieveAllThings(ctx, db, mf13things.PageMetadata{Offset: o, Limit: limit})
		if err != nil {
			return mf13things.Page{}, fmt.Errorf(retrieveErrString, err, o, limit, retrievThingsOps)
		}
		thingsPage.Things = append(thingsPage.Things, ths.Things...)
		o = o + 100
	}

	return thingsPage, nil
}

// MFRetrieveChannels retrieves existing channels from the database
func MFRetrieveChannels(ctx context.Context, db mf13postgres.Database) (mf13things.ChannelsPage, error) {
	channelsPage, err := RetrieveAllChannels(ctx, db, mf13things.PageMetadata{Offset: offset, Limit: limit})
	if err != nil {
		return mf13things.ChannelsPage{}, fmt.Errorf(retrieveErrString, err, offset, limit, retrievChannelsOps)
	}
	o := uint64(100)
	for o < channelsPage.Total {
		chs, err := RetrieveAllChannels(ctx, db, mf13things.PageMetadata{Offset: o, Limit: limit})
		if err != nil {
			return mf13things.ChannelsPage{}, fmt.Errorf(retrieveErrString, err, o, limit, retrievChannelsOps)
		}
		channelsPage.Channels = append(channelsPage.Channels, chs.Channels...)
		o = o + 100
	}

	return channelsPage, nil
}

// MFRetrieveConnections retrieves existing things to channels connection from the database
func MFRetrieveConnections(ctx context.Context, db mf13postgres.Database) (ConnectionsPage, error) {
	connectionsPage, err := RetrieveAllConnections(ctx, db, mf13things.PageMetadata{Offset: offset, Limit: limit})
	if err != nil {
		return ConnectionsPage{}, fmt.Errorf(retrieveErrString, err, offset, limit, retrievConnOps)
	}
	o := uint64(100)
	for o < connectionsPage.Total {
		conns, err := RetrieveAllConnections(ctx, db, mf13things.PageMetadata{Offset: o, Limit: limit})
		if err != nil {
			return ConnectionsPage{}, fmt.Errorf(retrieveErrString, err, o, limit, retrievConnOps)
		}
		connectionsPage.Connections = append(connectionsPage.Connections, conns.Connections...)
		o = o + 100
	}

	return connectionsPage, nil
}

// ThingsToCSV saves things to the provided csv file
// The format of the things csv file is ID,Key,Name,Owner,Metadata
func ThingsToCSV(filePath string, things []mf13things.Thing) error {
	if err := createrDir(filePath); err != nil {
		return fmt.Errorf(fileErrString, dirOp, err, writeThingsOps)
	}
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf(fileErrString, fileOp, err, writeThingsOps)
	}

	w := csv.NewWriter(f)

	records := [][]string{{"ID", "Key", "Name", "Owner", "Metadata"}}
	for _, thing := range things {
		record := []string{thing.ID, thing.Key, thing.Name, thing.Owner, fmt.Sprintf("%v", thing.Metadata)}
		records = append(records, record)
	}

	if err := w.WriteAll(records); err != nil {
		return fmt.Errorf(fileErrString, writeOp, err, writeThingsOps)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf(fileErrString, closeOp, err, writeThingsOps)
	}
	return nil
}

// ChannelsToCSV saves channels to the provided csv file
// The format of the channels csv file is ID,Name,Owner,Metadata
func ChannelsToCSV(filePath string, channels []mf13things.Channel) error {
	if err := createrDir(filePath); err != nil {
		return fmt.Errorf(fileErrString, dirOp, err, writeChannelsOps)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf(fileErrString, fileOp, err, writeChannelsOps)
	}

	w := csv.NewWriter(f)

	records := [][]string{{"ID", "Name", "Owner", "Metadata"}}
	for _, channel := range channels {
		record := []string{channel.ID, channel.Name, channel.Owner, fmt.Sprintf("%v", channel.Metadata)}
		records = append(records, record)
	}
	if err := w.WriteAll(records); err != nil {
		return fmt.Errorf(fileErrString, writeOp, err, writeChannelsOps)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf(fileErrString, closeOp, err, writeChannelsOps)
	}
	return nil
}

// ConnectionsToCSV saves connections to the provided csv file
// The format of the connections csv file is ChannelID,ChannelOwner,ThingID,ThingOwner
func ConnectionsToCSV(filePath string, connections []Connection) error {
	if err := createrDir(filePath); err != nil {
		return fmt.Errorf(fileErrString, dirOp, err, writeConnectionsOps)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf(fileErrString, fileOp, err, writeConnectionsOps)
	}

	w := csv.NewWriter(f)

	records := [][]string{{"ChannelID", "ChannelOwner", "ThingID", "ThingOwner"}}
	for _, conn := range connections {
		record := []string{conn.ChannelID, conn.ChannelOwner, conn.ThingID, conn.ThingOwner}
		records = append(records, record)
	}
	if err := w.WriteAll(records); err != nil {
		return fmt.Errorf(fileErrString, writeOp, err, writeConnectionsOps)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf(fileErrString, closeOp, err, writeConnectionsOps)
	}
	return nil
}

// CreateThings creates things from the provided csv file
// The format of the things csv file is ID,Key,Name,Owner,Metadata
func CreateThings(sdk mf14sdk.SDK, filePath string, token string) error {
	records, err := readData(filePath)
	if err != nil {
		return fmt.Errorf(readErrString, err, "creating things")
	}
	ths := []mf14sdk.Thing{}
	for _, record := range records {
		thing := mf14sdk.Thing{
			ID:    record[0],
			Name:  record[2],
			Owner: users.GetUserID(record[3]),
			Credentials: mf14sdk.Credentials{
				Secret: record[1],
			},
			Status: mf14sdk.EnabledStatus,
		}
		ths = append(ths, thing)
	}
	if _, err := sdk.CreateThings(ths, token); err != nil {
		return fmt.Errorf("failed to create things with error %v", err)
	}
	return nil
}

// CreateChannels creates channels from the provided csv file
// The format of the channels csv file is ID,Name,Owner,Metadata
func CreateChannels(sdk mf14sdk.SDK, filePath string, token string) error {
	records, err := readData(filePath)
	if err != nil {
		return fmt.Errorf(readErrString, err, "creating channels")
	}
	chs := []mf14sdk.Channel{}
	for _, record := range records {
		channel := mf14sdk.Channel{
			ID:      record[0],
			Name:    record[1],
			OwnerID: users.GetUserID(record[2]),
			Status:  mf14sdk.EnabledStatus,
		}
		chs = append(chs, channel)
	}
	if _, err := sdk.CreateChannels(chs, token); err != nil {
		return fmt.Errorf("failed to create channel with error %v", err)
	}
	return nil
}

// CreateConnections creates policies for things to read and write to the
// specified channels. The format of the connections csv file is
// ChannelID,ChannelOwner,ThingID,ThingOwner
func CreateConnections(sdk mf14sdk.SDK, filePath string, token string) error {
	records, err := readData(filePath)
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

func createrDir(path string) error {
	var dir string = filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	return nil
}

func readData(fileName string) ([][]string, error) {
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
