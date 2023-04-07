package things14

import (
	"context"
	"fmt"
	"log"
	"sync"

	mf14sdk "github.com/mainflux/mainflux/pkg/sdk/go/0140"
	util "github.com/mainflux/migrations/internal"
	users "github.com/mainflux/migrations/migrate/users/v14"
	"golang.org/x/sync/errgroup"
)

var limit = 100

// ReadAndCreateThings reads things from the provided csv file and creates them.
func ReadAndCreateThings(ctx context.Context, sdk mf14sdk.SDK, usersPath, filePath, token string) error {
	thchan := make(chan []string, limit)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return util.ReadInBatch(ctx, filePath, "creating things", thchan)
	})
	eg.Go(func() error {
		return createThings(ctx, sdk, usersPath, token, thchan)
	})

	return eg.Wait()
}

// createThings creates things from the provided csv file
// The format of the things csv file is ID,Key,Name,Owner,Metadata.
func createThings(ctx context.Context, sdk mf14sdk.SDK, usersPath, token string, inth <-chan []string) error {
	ths := []mf14sdk.Thing{}
	errCh := make(chan error)
	var wg sync.WaitGroup

	for record := range inth {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

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
		if len(ths) >= limit {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			wg.Add(1)
			go func(things []mf14sdk.Thing, errCh chan<- error) {
				defer wg.Done()

				if _, err := sdk.CreateThings(things, token); err != nil {
					errCh <- fmt.Errorf("Failed to create things with error %w", err)

					return
				}

				errCh <- nil
			}(ths, errCh)
			ths = []mf14sdk.Thing{}
		}
	}

	// Create remaining things
	if len(ths) > 0 {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if _, err := sdk.CreateThings(ths, token); err != nil {
			return fmt.Errorf("failed to create things with error %w", err)
		}
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadAndCreateChannels(ctx context.Context, sdk mf14sdk.SDK, usersPath, filePath, token string) error {
	chchan := make(chan []string, limit)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return util.ReadInBatch(ctx, filePath, "creating channels", chchan)
	})
	eg.Go(func() error {
		return createChannels(ctx, sdk, usersPath, token, chchan)
	})

	return eg.Wait()
}

// createChannels creates channels from the provided csv file
// The format of the channels csv file is ID,Name,Owner,Metadata.
func createChannels(ctx context.Context, sdk mf14sdk.SDK, usersPath, token string, inch <-chan []string) error {
	chs := []mf14sdk.Channel{}
	errCh := make(chan error)
	var wg sync.WaitGroup

	for record := range inch {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

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
		if len(chs) >= limit {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			wg.Add(1)
			go func(channels []mf14sdk.Channel, errCh chan<- error) {
				defer wg.Done()

				if _, err := sdk.CreateChannels(channels, token); err != nil {
					errCh <- fmt.Errorf("Failed to create things with error %w", err)

					return
				}

				errCh <- nil
			}(chs, errCh)
			chs = []mf14sdk.Channel{}
		}
	}

	// Create remaining channels
	if len(chs) > 0 {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if _, err := sdk.CreateChannels(chs, token); err != nil {
			return fmt.Errorf("failed to create channels with error %w", err)
		}
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadAndCreateConnections(ctx context.Context, sdk mf14sdk.SDK, filePath, token string) error {
	connchan := make(chan []string, limit)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return util.ReadInBatch(ctx, filePath, "creating connections", connchan)
	})
	eg.Go(func() error {
		return createConnections(sdk, token, connchan)
	})

	return eg.Wait()
}

// createConnections creates policies for things to read and write to the
// specified channels. The format of the connections csv file is
// ChannelID,ChannelOwner,ThingID,ThingOwner.
func createConnections(sdk mf14sdk.SDK, token string, inconn <-chan []string) error {
	thingIDsByChannelID := make(map[string][]string)
	for record := range inconn {
		channelID := record[0]
		thingID := record[2]
		if !contains(thingIDsByChannelID[channelID], thingID) {
			thingIDsByChannelID[channelID] = append(thingIDsByChannelID[channelID], thingID)
		}
	}

	var wg sync.WaitGroup

	conns := []mf14sdk.ConnectionIDs{}
	for channelID, thingIDs := range thingIDsByChannelID {
		conn := mf14sdk.ConnectionIDs{
			ChannelIDs: []string{channelID},
			ThingIDs:   thingIDs,
		}
		conns = append(conns, conn)
		if len(conns) >= limit {
			wg.Add(1)
			go func(conns []mf14sdk.ConnectionIDs) {
				for _, conn := range conns {
					if err := sdk.Connect(conn, token); err != nil {
						log.Fatalf("failed to connect things %v to channels %s with error %v", conn.ThingIDs, conn.ChannelIDs, err)
					}
				}
				defer wg.Done()
			}(conns)
			conns = []mf14sdk.ConnectionIDs{}
		}
	}

	if len(conns) > 0 {
		for _, conn := range conns {
			if err := sdk.Connect(conn, token); err != nil {
				log.Fatalf("failed to connect things %v to channels %s with error %v", conn.ThingIDs, conn.ChannelIDs, err)
			}
		}
	}

	wg.Wait()

	return nil
}

// Helper function to check if a thing contains a given element.
func contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}

	return false
}
