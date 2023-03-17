package things

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go/0140"
	things "github.com/mainflux/mainflux/things"
	thingsPostgres "github.com/mainflux/mainflux/things/postgres"
	"github.com/mainflux/migrations/migrate/users"
)

const (
	offset = uint64(0)
	limit  = uint64(100)
)

func MFRetrieveThings(ctx context.Context, db thingsPostgres.Database) (things.Page, error) {
	thingsPage, err := RetrieveAllThings(ctx, db, things.PageMetadata{Offset: offset, Limit: limit})
	if err != nil {
		return things.Page{}, err
	}
	o := uint64(100)
	for o < thingsPage.Total {
		ths, err := RetrieveAllThings(ctx, db, things.PageMetadata{Offset: o, Limit: limit})
		if err != nil {
			return things.Page{}, err
		}
		thingsPage.Things = append(thingsPage.Things, ths.Things...)
		o = o + 100
	}

	return thingsPage, nil
}

func MFRetrieveChannels(ctx context.Context, db thingsPostgres.Database) (things.ChannelsPage, error) {
	channelsPage, err := RetrieveAllChannels(ctx, db, things.PageMetadata{Offset: offset, Limit: limit})
	if err != nil {
		return things.ChannelsPage{}, err
	}
	o := uint64(100)
	for o < channelsPage.Total {
		chs, err := RetrieveAllChannels(ctx, db, things.PageMetadata{Offset: o, Limit: limit})
		if err != nil {
			return things.ChannelsPage{}, err
		}
		channelsPage.Channels = append(channelsPage.Channels, chs.Channels...)
		o = o + 100
	}

	return channelsPage, nil
}

func MFRetrieveConnections(ctx context.Context, db thingsPostgres.Database) (ConnectionsPage, error) {
	connectionsPage, err := RetrieveAllConnections(ctx, db, things.PageMetadata{Offset: offset, Limit: limit})
	if err != nil {
		return ConnectionsPage{}, err
	}
	o := uint64(100)
	for o < connectionsPage.Total {
		conns, err := RetrieveAllConnections(ctx, db, things.PageMetadata{Offset: o, Limit: limit})
		if err != nil {
			return ConnectionsPage{}, err
		}
		connectionsPage.Connections = append(connectionsPage.Connections, conns.Connections...)
		o = o + 100
	}

	return connectionsPage, nil
}

func ThingsToCSV(filePath string, things []things.Thing) error {
	if err := createrDir(filePath); err != nil {
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)

	records := [][]string{{"ID", "Key", "Name", "Owner", "Metadata"}}
	for _, thing := range things {
		record := []string{thing.ID, thing.Key, thing.Name, thing.Owner, fmt.Sprintf("%v", thing.Metadata)}
		records = append(records, record)
	}

	if err := w.WriteAll(records); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func ChannelsToCSV(filePath string, channels []things.Channel) error {
	if err := createrDir(filePath); err != nil {
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)

	records := [][]string{{"ID", "Name", "Owner", "Metadata"}}
	for _, channel := range channels {
		record := []string{channel.ID, channel.Name, channel.Owner, fmt.Sprintf("%v", channel.Metadata)}
		records = append(records, record)
	}
	if err := w.WriteAll(records); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func CreateThings(sdk mfsdk.SDK, filePath string, token string) error {
	records, err := readData(filePath)
	if err != nil {
		return err
	}
	ths := []mfsdk.Thing{}
	for _, record := range records {
		thing := mfsdk.Thing{
			ID:    record[0],
			Name:  record[2],
			Owner: users.GetUserID(record[3]),
			Credentials: mfsdk.Credentials{
				Secret: record[1],
			},
			Status: mfsdk.EnabledStatus,
		}
		ths = append(ths, thing)
	}
	if _, err := sdk.CreateThings(ths, token); err != nil {
		return err
	}
	return nil
}

func CreateChannels(sdk mfsdk.SDK, filePath string, token string) error {
	records, err := readData(filePath)
	if err != nil {
		return err
	}
	chs := []mfsdk.Channel{}
	for _, record := range records {
		channel := mfsdk.Channel{
			ID:      record[0],
			Name:    record[1],
			OwnerID: users.GetUserID(record[2]),
			Status:  mfsdk.EnabledStatus,
		}
		chs = append(chs, channel)
	}
	if _, err := sdk.CreateChannels(chs, token); err != nil {
		return err
	}
	return nil
}

func CreateConnections(sdk mfsdk.SDK, filePath string, token string) error {
	records, err := readData(filePath)
	if err != nil {
		return err
	}

	for _, record := range records {
		// ChannelID,ChannelOwner,ThingID,ThingOwner
		if err := sdk.ConnectThing(record[2], record[0], token); err != nil {
			return err
		}
	}

	return nil
}

func ConnectionsToCSV(filePath string, connections []Connection) error {
	if err := createrDir(filePath); err != nil {
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)

	records := [][]string{{"ChannelID", "ChannelOwner", "ThingID", "ThingOwner"}}
	for _, conn := range connections {
		record := []string{conn.ChannelID, conn.ChannelOwner, conn.ThingID, conn.ThingOwner}
		records = append(records, record)
	}
	if err := w.WriteAll(records); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		return err
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
