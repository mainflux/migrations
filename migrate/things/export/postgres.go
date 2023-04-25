package export

import (
	"context"
	"encoding/json"

	mf13errors "github.com/mainflux/mainflux/pkg/errors"
	"github.com/mainflux/mainflux/things/postgres" // for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0
)

// dbRetrieveThings retrieves things from the database with the given page navigation parameters.
func dbRetrieveThings(ctx context.Context, query string, database postgres.Database, pm pageMetadata) (thingsPage, error) {
	params := map[string]interface{}{
		"limit":  pm.Limit,
		"offset": pm.Offset,
	}

	rows, err := database.NamedQueryContext(ctx, query, params)
	if err != nil {
		return thingsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}
	defer rows.Close()

	var items []thing
	for rows.Next() {
		dbth := dbThing{}
		if err := rows.StructScan(&dbth); err != nil {
			return thingsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
		}

		th, err := toThing(dbth)
		if err != nil {
			return thingsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
		}

		items = append(items, th)
	}

	cq := `SELECT COUNT(*) FROM things;`

	total, err := total(ctx, database, cq, params)
	if err != nil {
		return thingsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}

	page := thingsPage{
		Things: items,
		pageMetadata: pageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}

	return page, nil
}

// dbRetrieveChannels retrieves things from the database with the given page navigation parameters.
func dbRetrieveChannels(ctx context.Context, query string, database postgres.Database, pm pageMetadata) (channelsPage, error) {
	params := map[string]interface{}{
		"limit":  pm.Limit,
		"offset": pm.Offset,
	}

	rows, err := database.NamedQueryContext(ctx, query, params)
	if err != nil {
		return channelsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}
	defer rows.Close()

	var items []channel
	for rows.Next() {
		dbch := dbChannel{}
		if err := rows.StructScan(&dbch); err != nil {
			return channelsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
		}

		ch, err := toChannel(dbch)
		if err != nil {
			return channelsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
		}

		items = append(items, ch)
	}

	cq := `SELECT COUNT(*) FROM channels;`

	total, err := total(ctx, database, cq, params)
	if err != nil {
		return channelsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}

	page := channelsPage{
		Channels: items,
		pageMetadata: pageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}

	return page, nil
}

// dbRetrieveConnections retrieves things from the database with the given page navigation parameters.
func dbRetrieveConnections(ctx context.Context, query string, database postgres.Database, pm pageMetadata) (connectionsPage, error) {
	params := map[string]interface{}{
		"limit":  pm.Limit,
		"offset": pm.Offset,
	}

	rows, err := database.NamedQueryContext(ctx, query, params)
	if err != nil {
		return connectionsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}
	defer rows.Close()

	var items []connection
	for rows.Next() {
		dbconn := connection{}
		if err := rows.StructScan(&dbconn); err != nil {
			return connectionsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
		}

		conn := dbconn
		items = append(items, conn)
	}

	cq := `SELECT COUNT(*) FROM connections;`

	total, err := total(ctx, database, cq, params)
	if err != nil {
		return connectionsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}

	page := connectionsPage{
		Connections: items,
		pageMetadata: pageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}

	return page, nil
}

func total(ctx context.Context, database postgres.Database, query string, params interface{}) (uint64, error) {
	rows, err := database.NamedQueryContext(ctx, query, params)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	total := uint64(0)
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, err
		}
	}

	return total, nil
}

// toThing converts a dbThing to a thing.
// this is for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0.
func toThing(dbth dbThing) (thing, error) {
	var metadata map[string]interface{}
	if dbth.Metadata != nil {
		if err := json.Unmarshal(dbth.Metadata, &metadata); err != nil {
			return thing{}, mf13errors.Wrap(mf13errors.ErrMalformedEntity, err)
		}
	}

	return thing{
		ID:       dbth.ID,
		Owner:    dbth.Owner,
		Name:     dbth.Name,
		Key:      dbth.Key,
		Metadata: metadata,
	}, nil
}

// toChannel converts a dbChannel to a channel.
// this is for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0.
func toChannel(dch dbChannel) (channel, error) {
	var metadata map[string]interface{}
	if dch.Metadata != nil {
		if err := json.Unmarshal(dch.Metadata, &metadata); err != nil {
			return channel{}, mf13errors.Wrap(mf13errors.ErrMalformedEntity, err)
		}
	}

	return channel{
		ID:       dch.ID,
		Owner:    dch.Owner,
		Name:     dch.Name,
		Metadata: metadata,
	}, nil
}
