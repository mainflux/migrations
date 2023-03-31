package things14

import (
	"fmt"
	"log"
	"sync"

	mf14sdk "github.com/mainflux/mainflux/pkg/sdk/go/0140"
	util "github.com/mainflux/migrations/internal"
	users "github.com/mainflux/migrations/migrate/users/v14"
	"golang.org/x/sync/errgroup"
)

var (
	limit = 100
)

// ReadAndCreateThings reads things from the provided csv file and creates them
func ReadAndCreateThings(sdk mf14sdk.SDK, usersPath, filePath, token string) error {
	g := new(errgroup.Group)
	thchan := make(chan []string, limit)

	g.Go(func() error {
		return util.ReadInBatch(filePath, "creating things", thchan)
	})
	g.Go(func() error {
		return CreateThings(sdk, usersPath, token, thchan)
	})

	return g.Wait()
}

// CreateThings creates things from the provided csv file
// The format of the things csv file is ID,Key,Name,Owner,Metadata
func CreateThings(sdk mf14sdk.SDK, usersPath, token string, inth <-chan []string) error {
	ths := []mf14sdk.Thing{}
	var wg sync.WaitGroup

	for record := range inth {
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
			wg.Add(1)
			go func(things []mf14sdk.Thing) {
				if _, err := sdk.CreateThings(things, token); err != nil {
					log.Fatalf("Failed to create things with error %v", err)
				}
				defer wg.Done()
			}(ths)
			ths = []mf14sdk.Thing{}
		}
	}

	// Create remaining things
	if len(ths) > 0 {
		if _, err := sdk.CreateThings(ths, token); err != nil {
			return fmt.Errorf("failed to create things with error %v", err)
		}
	}

	wg.Wait()
	return nil
}

func ReadAndCreateChannels(sdk mf14sdk.SDK, usersPath, filePath, token string) error {
	g := new(errgroup.Group)
	chchan := make(chan []string, limit)

	g.Go(func() error {
		return util.ReadInBatch(filePath, "creating channels", chchan)
	})
	g.Go(func() error {
		return CreateChannels(sdk, usersPath, token, chchan)
	})

	return g.Wait()
}

// CreateChannels creates channels from the provided csv file
// The format of the channels csv file is ID,Name,Owner,Metadata
func CreateChannels(sdk mf14sdk.SDK, usersPath, token string, inch <-chan []string) error {
	chs := []mf14sdk.Channel{}
	var wg sync.WaitGroup

	for record := range inch {
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
			wg.Add(1)
			go func(channels []mf14sdk.Channel) {
				if _, err := sdk.CreateChannels(channels, token); err != nil {
					log.Fatalf("Failed to create things with error %v", err)
				}
				defer wg.Done()
			}(chs)
			chs = []mf14sdk.Channel{}
		}
	}

	// Create remaining channels
	if len(chs) > 0 {
		if _, err := sdk.CreateChannels(chs, token); err != nil {
			return fmt.Errorf("failed to create channels with error %v", err)
		}
	}

	wg.Wait()
	return nil
}

func ReadAndCreateConnections(sdk mf14sdk.SDK, filePath, token string) error {
	g := new(errgroup.Group)
	connchan := make(chan []string, limit)

	g.Go(func() error {
		return util.ReadInBatch(filePath, "creating connections", connchan)
	})
	g.Go(func() error {
		return CreateConnections(sdk, token, connchan)
	})

	return g.Wait()
}

// CreateConnections creates policies for things to read and write to the
// specified channels. The format of the connections csv file is
// ChannelID,ChannelOwner,ThingID,ThingOwner
func CreateConnections(sdk mf14sdk.SDK, token string, inconn <-chan []string) error {
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
						log.Fatalf("failed to connect things with ids %v to channel with id %s with error %v", conn.ThingIDs, conn.ChannelIDs, err)
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
				log.Fatalf("failed to connect things with ids %v to channel with id %s with error %v", conn.ThingIDs, conn.ChannelIDs, err)
			}
		}
	}

	wg.Wait()

	return nil
}

// Helper function to check if a thing contains a given element
func contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
