package things13

import (
	"context"
	"encoding/csv"
	"fmt"

	"golang.org/x/sync/errgroup"

	mf13things "github.com/mainflux/mainflux/things"
	mf13postgres "github.com/mainflux/mainflux/things/postgres"
	util "github.com/mainflux/migrations/internal"
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

// RetrieveAndWriteThings retrieves existing things from the database and saves them to the provided csv file
func RetrieveAndWriteThings(ctx context.Context, db mf13postgres.Database, filePath string) error {
	out := make(chan []mf13things.Thing)

	g := new(errgroup.Group)

	g.Go(func() error {
		return RetrieveThings(ctx, db, out)
	})
	g.Go(func() error {
		return ThingsToCSV(filePath, out)
	})

	return g.Wait()
}

// RetrieveThings retrieves existing things from the database
func RetrieveThings(ctx context.Context, db mf13postgres.Database, outth chan<- []mf13things.Thing) error {
	o := offset
	l := limit
	for {
		thingsPage, err := dbRetrieveThings(ctx, db, mf13things.PageMetadata{Offset: o, Limit: l})
		if err != nil {
			return fmt.Errorf(retrieveErrString, err, o, l, retrievThingsOps)
		}
		outth <- thingsPage.Things
		if o+l >= thingsPage.Total {
			break
		}

		o = o + l
		l = util.UpdateLimit(thingsPage.Total)
	}
	close(outth)
	return nil

}

// ThingsToCSV saves things to the provided csv file
// The format of the things csv file is ID,Key,Name,Owner,Metadata
func ThingsToCSV(filePath string, inth <-chan []mf13things.Thing) error {
	f, err := util.CreateFile(filePath, writeThingsOps)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)

	records := [][]string{{"ID", "Key", "Name", "Owner", "Metadata"}}
	for things := range inth {
		for _, thing := range things {
			metadata, err := util.MetadataToString(thing.Metadata)
			if err != nil {
				return err
			}
			record := []string{thing.ID, thing.Key, thing.Name, thing.Owner, metadata}
			records = append(records, record)
		}
	}

	return util.WriteData(w, f, records, writeThingsOps)
}

// RetrieveAndWriteChannels retrieves existing channels from the database and saves them to the provided csv file
func RetrieveAndWriteChannels(ctx context.Context, db mf13postgres.Database, filePath string) error {
	out := make(chan []mf13things.Channel)

	g := new(errgroup.Group)

	g.Go(func() error {
		return RetrieveChannels(ctx, db, out)
	})
	g.Go(func() error {
		return ChannelsToCSV(filePath, out)
	})

	return g.Wait()
}

// RetrieveChannels retrieves existing channels from the database
func RetrieveChannels(ctx context.Context, db mf13postgres.Database, outch chan<- []mf13things.Channel) error {
	o := offset
	l := limit
	for {
		channelsPage, err := dbRetrieveChannels(ctx, db, mf13things.PageMetadata{Offset: o, Limit: l})
		if err != nil {
			return fmt.Errorf(retrieveErrString, err, o, l, retrievChannelsOps)
		}
		outch <- channelsPage.Channels
		if o+l >= channelsPage.Total {
			break
		}

		o = o + l
		l = util.UpdateLimit(channelsPage.Total)
	}
	close(outch)
	return nil
}

// ChannelsToCSV saves channels to the provided csv file
// The format of the channels csv file is ID,Name,Owner,Metadata
func ChannelsToCSV(filePath string, inch <-chan []mf13things.Channel) error {
	f, err := util.CreateFile(filePath, writeChannelsOps)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)

	records := [][]string{{"ID", "Name", "Owner", "Metadata"}}
	for channels := range inch {
		for _, channel := range channels {
			metadata, err := util.MetadataToString(channel.Metadata)
			if err != nil {
				return err
			}
			record := []string{channel.ID, channel.Name, channel.Owner, metadata}
			records = append(records, record)
		}
	}
	return util.WriteData(w, f, records, writeChannelsOps)
}

// RetrieveAndWriteConnections retrieves existing things to channels connection from the database and saves them to the provided csv file
func RetrieveAndWriteConnections(ctx context.Context, db mf13postgres.Database, filePath string) error {
	out := make(chan []Connection)

	g := new(errgroup.Group)

	g.Go(func() error {
		return RetrieveConnections(ctx, db, out)
	})
	g.Go(func() error {
		return ConnectionsToCSV(filePath, out)
	})

	return g.Wait()
}

// RetrieveConnections retrieves existing things to channels connection from the database
func RetrieveConnections(ctx context.Context, db mf13postgres.Database, outconn chan<- []Connection) error {
	o := offset
	l := limit
	for {
		connectionsPage, err := dbRetrieveConnections(ctx, db, mf13things.PageMetadata{Offset: o, Limit: l})
		if err != nil {
			return fmt.Errorf(retrieveErrString, err, offset, limit, retrievConnOps)
		}
		outconn <- connectionsPage.Connections
		if o+l >= connectionsPage.Total {
			break
		}

		o = o + l
		l = util.UpdateLimit(connectionsPage.Total)
	}
	close(outconn)
	return nil
}

// ConnectionsToCSV saves connections to the provided csv file
// The format of the connections csv file is ChannelID,ChannelOwner,ThingID,ThingOwner
func ConnectionsToCSV(filePath string, inconn <-chan []Connection) error {
	f, err := util.CreateFile(filePath, writeConnectionsOps)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)

	records := [][]string{{"ChannelID", "ChannelOwner", "ThingID", "ThingOwner"}}
	for connections := range inconn {
		for _, conn := range connections {
			record := []string{conn.ChannelID, conn.ChannelOwner, conn.ThingID, conn.ThingOwner}
			records = append(records, record)
		}
	}
	return util.WriteData(w, f, records, writeConnectionsOps)
}
