package things13

import (
	"context"
	"encoding/csv"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	mf13things "github.com/mainflux/mainflux/things"
	mf13postgres "github.com/mainflux/mainflux/things/postgres"
	util "github.com/mainflux/migrations/internal"
)

var (
	defaultOffset         = uint64(0)
	retrievThingsOps      = "retrieving things"
	retrievChannelsOps    = "retrieving channels"
	retrievConnOps        = "retrieving connections"
	writeThingsOps        = "writing things to csv file"
	writeChannelsOps      = "writing channels to csv file"
	writeConnectionsOps   = "writing connections to csv file"
	retrieveErrString     = "error %v occured at offset: %d and total: %d during %s"
	totalThingsQuery      = "SELECT COUNT(*) FROM things;"
	totalChannelsQuery    = "SELECT COUNT(*) FROM channels;"
	totalConnectionsQuery = "SELECT COUNT(*) FROM connections;"
)

// RetrieveAndWriteThings retrieves existing things from the database and saves them to the provided csv file
func RetrieveAndWriteThings(ctx context.Context, db mf13postgres.Database, filePath string) error {
	out := make(chan []mf13things.Thing)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer close(out)
		return RetrieveThings(ctx, db, out)
	})
	g.Go(func() error {
		return ThingsToCSV(ctx, filePath, out)
	})

	return g.Wait()
}

// RetrieveThings retrieves existing things from the database
func RetrieveThings(ctx context.Context, db mf13postgres.Database, allThings chan<- []mf13things.Thing) error {
	totolThings, err := total(ctx, db, totalThingsQuery, map[string]interface{}{})
	if err != nil {
		return err
	}
	offset := defaultOffset
	limit := util.UpdateLimit(totolThings)

	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	for {
		wg.Add(1)
		go func(offset uint64, errCh chan<- error) {
			defer wg.Done()

			thingsPage, err := dbRetrieveThings(ctx, db, mf13things.PageMetadata{Offset: offset, Limit: limit})
			if err != nil {
				errCh <- fmt.Errorf(retrieveErrString, err, offset, limit, retrievThingsOps)
			}

			select {
			case <-ctx.Done():
				return
			case allThings <- thingsPage.Things:
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

		if offset+limit >= totolThings {
			break
		}
		offset = offset + limit
	}

	wg.Wait()
	return nil
}

// ThingsToCSV saves things to the provided csv file
// The format of the things csv file is ID,Key,Name,Owner,Metadata
func ThingsToCSV(ctx context.Context, filePath string, allThings <-chan []mf13things.Thing) error {
	f, err := util.CreateFile(filePath, writeThingsOps)
	if err != nil {
		return err
	}

	defer func() {
		if ferr := f.Close(); ferr != nil && err == nil {
			err = ferr
		}
	}()

	w := csv.NewWriter(f)

	records := [][]string{{"ID", "Key", "Name", "Owner", "Metadata"}}
	for things := range allThings {
		for _, thing := range things {
			metadata, err := util.MetadataToString(thing.Metadata)
			if err != nil {
				return err
			}
			record := []string{thing.ID, thing.Key, thing.Name, thing.Owner, metadata}
			records = append(records, record)
		}
	}

	return util.WriteData(ctx, w, f, records, writeThingsOps)
}

// RetrieveAndWriteChannels retrieves existing channels from the database and saves them to the provided csv file
func RetrieveAndWriteChannels(ctx context.Context, db mf13postgres.Database, filePath string) error {
	out := make(chan []mf13things.Channel)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer close(out)
		return RetrieveChannels(ctx, db, out)
	})
	g.Go(func() error {
		return ChannelsToCSV(ctx, filePath, out)
	})

	return g.Wait()
}

// RetrieveChannels retrieves existing channels from the database
func RetrieveChannels(ctx context.Context, db mf13postgres.Database, allChannels chan<- []mf13things.Channel) error {
	totolChannels, err := total(ctx, db, totalChannelsQuery, map[string]interface{}{})
	if err != nil {
		return err
	}
	offset := defaultOffset
	limit := util.UpdateLimit(totolChannels)

	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	for {
		wg.Add(1)
		go func(offset uint64, errCh chan<- error) {
			defer wg.Done()

			channelsPage, err := dbRetrieveChannels(ctx, db, mf13things.PageMetadata{Offset: offset, Limit: limit})
			if err != nil {
				errCh <- fmt.Errorf(retrieveErrString, err, offset, limit, retrievChannelsOps)
			}

			select {
			case <-ctx.Done():
				return
			case allChannels <- channelsPage.Channels:
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

		if offset+limit >= totolChannels {
			break
		}
		offset = offset + limit
	}

	wg.Wait()
	return nil
}

// ChannelsToCSV saves channels to the provided csv file
// The format of the channels csv file is ID,Name,Owner,Metadata
func ChannelsToCSV(ctx context.Context, filePath string, allChannels <-chan []mf13things.Channel) error {
	f, err := util.CreateFile(filePath, writeChannelsOps)
	if err != nil {
		return err
	}

	defer func() {
		if ferr := f.Close(); ferr != nil && err == nil {
			err = ferr
		}
	}()

	w := csv.NewWriter(f)

	records := [][]string{{"ID", "Name", "Owner", "Metadata"}}
	for channels := range allChannels {
		for _, channel := range channels {
			metadata, err := util.MetadataToString(channel.Metadata)
			if err != nil {
				return err
			}
			record := []string{channel.ID, channel.Name, channel.Owner, metadata}
			records = append(records, record)
		}
	}

	return util.WriteData(ctx, w, f, records, writeChannelsOps)
}

// RetrieveAndWriteConnections retrieves existing things to channels connection from the database and saves them to the provided csv file
func RetrieveAndWriteConnections(ctx context.Context, db mf13postgres.Database, filePath string) error {
	out := make(chan []Connection)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer close(out)
		return RetrieveConnections(ctx, db, out)
	})
	g.Go(func() error {
		return ConnectionsToCSV(ctx, filePath, out)
	})

	return g.Wait()
}

// RetrieveConnections retrieves existing things to channels connection from the database
func RetrieveConnections(ctx context.Context, db mf13postgres.Database, allConnections chan<- []Connection) error {
	totolConnections, err := total(ctx, db, totalConnectionsQuery, map[string]interface{}{})
	if err != nil {
		return err
	}
	offset := defaultOffset
	limit := util.UpdateLimit(totolConnections)

	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	for {
		wg.Add(1)
		go func(offset uint64, errCh chan<- error) {
			defer wg.Done()

			connectionsPage, err := dbRetrieveConnections(ctx, db, mf13things.PageMetadata{Offset: offset, Limit: limit})
			if err != nil {
				errCh <- fmt.Errorf(retrieveErrString, err, offset, limit, retrievConnOps)
			}

			select {
			case <-ctx.Done():
				return
			case allConnections <- connectionsPage.Connections:
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
		if offset+limit >= totolConnections {
			break
		}
		offset = offset + limit
	}

	wg.Wait()
	return nil
}

// ConnectionsToCSV saves connections to the provided csv file
// The format of the connections csv file is ChannelID,ChannelOwner,ThingID,ThingOwner
func ConnectionsToCSV(ctx context.Context, filePath string, inconn <-chan []Connection) error {
	f, err := util.CreateFile(filePath, writeConnectionsOps)
	if err != nil {
		return err
	}

	defer func() {
		if ferr := f.Close(); ferr != nil && err == nil {
			err = ferr
		}
	}()

	w := csv.NewWriter(f)

	records := [][]string{{"ChannelID", "ChannelOwner", "ThingID", "ThingOwner"}}
	for connections := range inconn {
		for _, conn := range connections {
			record := []string{conn.ChannelID, conn.ChannelOwner, conn.ThingID, conn.ThingOwner}
			records = append(records, record)
		}
	}

	return util.WriteData(ctx, w, f, records, writeConnectionsOps)
}
