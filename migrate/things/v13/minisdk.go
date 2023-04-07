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
	retrieveErrString     = "error %v occurred at offset: %d and total: %d during %s"
	totalThingsQuery      = "SELECT COUNT(*) FROM things;"
	totalChannelsQuery    = "SELECT COUNT(*) FROM channels;"
	totalConnectionsQuery = "SELECT COUNT(*) FROM connections;"
)

// RetrieveAndWriteThings retrieves existing things from the database and saves them to the provided csv file.
func RetrieveAndWriteThings(ctx context.Context, database mf13postgres.Database, filePath string) error {
	out := make(chan []mf13things.Thing)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		defer close(out)

		return retrieveThings(ctx, database, out)
	})
	eg.Go(func() error {
		return thingsToCSV(ctx, filePath, out)
	})

	return eg.Wait()
}

// retrieveThings retrieves existing things from the database.
func retrieveThings(ctx context.Context, database mf13postgres.Database, allThings chan<- []mf13things.Thing) error {
	totolThings, err := total(ctx, database, totalThingsQuery, map[string]interface{}{})
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

			thingsPage, err := dbRetrieveThings(ctx, database, mf13things.PageMetadata{Offset: offset, Limit: limit})
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
		offset += limit
	}

	wg.Wait()

	return nil
}

// thingsToCSV saves things to the provided csv file
// The format of the things csv file is ID,Key,Name,Owner,Metadata.
func thingsToCSV(ctx context.Context, filePath string, allThings <-chan []mf13things.Thing) error {
	file, err := util.CreateFile(filePath, writeThingsOps)
	if err != nil {
		return err
	}

	defer func() {
		if ferr := file.Close(); ferr != nil && err == nil {
			err = ferr
		}
	}()

	writer := csv.NewWriter(file)

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

	return util.WriteData(ctx, writer, file, records, writeThingsOps)
}

// RetrieveAndWriteChannels retrieves existing channels from the database and saves them to the provided csv file.
func RetrieveAndWriteChannels(ctx context.Context, database mf13postgres.Database, filePath string) error {
	out := make(chan []mf13things.Channel)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		defer close(out)

		return retrieveChannels(ctx, database, out)
	})
	eg.Go(func() error {
		return channelsToCSV(ctx, filePath, out)
	})

	return eg.Wait()
}

// retrieveChannels retrieves existing channels from the database.
func retrieveChannels(ctx context.Context, database mf13postgres.Database, allChannels chan<- []mf13things.Channel) error {
	totolChannels, err := total(ctx, database, totalChannelsQuery, map[string]interface{}{})
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

			channelsPage, err := dbRetrieveChannels(ctx, database, mf13things.PageMetadata{Offset: offset, Limit: limit})
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
		offset += limit
	}

	wg.Wait()

	return nil
}

// channelsToCSV saves channels to the provided csv file
// The format of the channels csv file is ID,Name,Owner,Metadata.
func channelsToCSV(ctx context.Context, filePath string, allChannels <-chan []mf13things.Channel) error {
	file, err := util.CreateFile(filePath, writeChannelsOps)
	if err != nil {
		return err
	}

	defer func() {
		if ferr := file.Close(); ferr != nil && err == nil {
			err = ferr
		}
	}()

	writer := csv.NewWriter(file)

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

	return util.WriteData(ctx, writer, file, records, writeChannelsOps)
}

// RetrieveAndWriteConnections retrieves existing things to channels connection from the database and saves them to the provided csv file.
func RetrieveAndWriteConnections(ctx context.Context, database mf13postgres.Database, filePath string) error {
	out := make(chan []Connection)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		defer close(out)

		return retrieveConnections(ctx, database, out)
	})
	eg.Go(func() error {
		return connectionsToCSV(ctx, filePath, out)
	})

	return eg.Wait()
}

// retrieveConnections retrieves existing things to channels connection from the database.
func retrieveConnections(ctx context.Context, database mf13postgres.Database, allConnections chan<- []Connection) error {
	totolConnections, err := total(ctx, database, totalConnectionsQuery, map[string]interface{}{})
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

			connectionsPage, err := dbRetrieveConnections(ctx, database, mf13things.PageMetadata{Offset: offset, Limit: limit})
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
		offset += limit
	}

	wg.Wait()

	return nil
}

// connectionsToCSV saves connections to the provided csv file
// The format of the connections csv file is ChannelID,ChannelOwner,ThingID,ThingOwner.
func connectionsToCSV(ctx context.Context, filePath string, inconn <-chan []Connection) error {
	file, err := util.CreateFile(filePath, writeConnectionsOps)
	if err != nil {
		return err
	}

	defer func() {
		if ferr := file.Close(); ferr != nil && err == nil {
			err = ferr
		}
	}()

	writer := csv.NewWriter(file)

	records := [][]string{{"ChannelID", "ChannelOwner", "ThingID", "ThingOwner"}}
	for connections := range inconn {
		for _, conn := range connections {
			record := []string{conn.ChannelID, conn.ChannelOwner, conn.ThingID, conn.ThingOwner}
			records = append(records, record)
		}
	}

	return util.WriteData(ctx, writer, file, records, writeConnectionsOps)
}
