package importthings

import (
	"context"
	"fmt"
	"log"
	"sync"

	mf13sdk "github.com/mainflux/mainflux/pkg/sdk/go/0130"
	"github.com/mainflux/migrations"
	util "github.com/mainflux/migrations/internal"
	"golang.org/x/sync/errgroup"
)

// InitSDKv13 initializes the SDK and creates a token.
func InitSDKv13(cfg migrations.Config) (mf13sdk.SDK, string, error) {
	sdkConf := mf13sdk.Config{
		UsersURL:        cfg.UsersURL,
		ThingsURL:       cfg.ThingsURL,
		MsgContentType:  mf13sdk.CTJSONSenML,
		TLSVerification: false,
	}

	sdk := mf13sdk.NewSDK(sdkConf)
	user := mf13sdk.User{
		Email:    cfg.UserIdentity,
		Password: cfg.UserSecret,
	}
	token, err := sdk.CreateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create token with error %w", err)
	}

	return sdk, token, nil
}

// ReadAndCreateThingsv13 reads things from the provided csv file and creates them.
func ReadAndCreateThingsv13(ctx context.Context, sdk mf13sdk.SDK, _, filePath, token string) error {
	thchan := make(chan []string, limit)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return util.ReadInBatch(ctx, filePath, "creating things", thchan)
	})
	eg.Go(func() error {
		return createThingsv13(ctx, sdk, token, thchan)
	})

	return eg.Wait()
}

// createThingsv13 creates things from the provided csv file
// The format of the things csv file is ID,Key,Name,Owner,Metadata.
func createThingsv13(ctx context.Context, sdk mf13sdk.SDK, token string, inth <-chan []string) error {
	ths := []mf13sdk.Thing{}
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
		thing := mf13sdk.Thing{
			ID:       record[0],
			Name:     record[2],
			Key:      record[1],
			Metadata: metadata,
		}
		ths = append(ths, thing)
		if len(ths) >= limit {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			wg.Add(1)
			go func(things []mf13sdk.Thing, errCh chan<- error) {
				defer wg.Done()

				if _, err := sdk.CreateThings(things, token); err != nil {
					errCh <- fmt.Errorf("failed to create things with error %w", err)

					return
				}

				errCh <- nil
			}(ths, errCh)
			ths = []mf13sdk.Thing{}
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

// ReadAndCreateChannelsv13 reads channels from the provided csv file and creates them.
func ReadAndCreateChannelsv13(ctx context.Context, sdk mf13sdk.SDK, _, filePath, token string) error {
	chchan := make(chan []string, limit)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return util.ReadInBatch(ctx, filePath, "creating channels", chchan)
	})
	eg.Go(func() error {
		return createChannelsv13(ctx, sdk, token, chchan)
	})

	return eg.Wait()
}

// createChannelsv13 creates channels from the provided csv file
// The format of the channels csv file is ID,Name,Owner,Metadata.
func createChannelsv13(ctx context.Context, sdk mf13sdk.SDK, token string, inch <-chan []string) error {
	chs := []mf13sdk.Channel{}
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
		channel := mf13sdk.Channel{
			ID:       record[0],
			Name:     record[1],
			Metadata: metadata,
		}
		chs = append(chs, channel)
		if len(chs) >= limit {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			wg.Add(1)
			go func(channels []mf13sdk.Channel, errCh chan<- error) {
				defer wg.Done()

				if _, err := sdk.CreateChannels(channels, token); err != nil {
					errCh <- fmt.Errorf("failed to create things with error %w", err)

					return
				}

				errCh <- nil
			}(chs, errCh)
			chs = []mf13sdk.Channel{}
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

// ReadAndCreateConnectionsv13 reads connections from the database and creates them.
func ReadAndCreateConnectionsv13(ctx context.Context, sdk mf13sdk.SDK, filePath, token string) error {
	connchan := make(chan []string, limit)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return util.ReadInBatch(ctx, filePath, "creating connections", connchan)
	})
	eg.Go(func() error {
		return createConnectionsv13(sdk, token, connchan)
	})

	return eg.Wait()
}

// createConnectionsv13 creates policies for things to read and write to the
// specified channels. The format of the connections csv file is
// ChannelID,ChannelOwner,ThingID,ThingOwner.
func createConnectionsv13(sdk mf13sdk.SDK, token string, inconn <-chan []string) error {
	thingIDsByChannelID := make(map[string][]string)
	for record := range inconn {
		channelID := record[0]
		thingID := record[2]
		if !contains(thingIDsByChannelID[channelID], thingID) {
			thingIDsByChannelID[channelID] = append(thingIDsByChannelID[channelID], thingID)
		}
	}

	var wg sync.WaitGroup

	conns := []mf13sdk.ConnectionIDs{}
	for channelID, thingIDs := range thingIDsByChannelID {
		conn := mf13sdk.ConnectionIDs{
			ChannelIDs: []string{channelID},
			ThingIDs:   thingIDs,
		}
		conns = append(conns, conn)
		if len(conns) >= limit {
			wg.Add(1)
			go func(conns []mf13sdk.ConnectionIDs) {
				for _, conn := range conns {
					if err := sdk.Connect(conn, token); err != nil {
						log.Fatalf("failed to connect things %v to channels %s with error %v", conn.ThingIDs, conn.ChannelIDs, err)
					}
				}
				defer wg.Done()
			}(conns)
			conns = []mf13sdk.ConnectionIDs{}
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
