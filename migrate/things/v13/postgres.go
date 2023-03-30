package things13

import (
	"context"
	"encoding/json"

	mf13errors "github.com/mainflux/mainflux/pkg/errors"
	mf13things "github.com/mainflux/mainflux/things"
	mf13postgres "github.com/mainflux/mainflux/things/postgres"
)

type dbThing struct {
	ID       string `db:"id"`
	Owner    string `db:"owner"`
	Name     string `db:"name"`
	Key      string `db:"key"`
	Metadata []byte `db:"metadata"`
}

type dbChannel struct {
	ID       string `db:"id"`
	Owner    string `db:"owner"`
	Name     string `db:"name"`
	Metadata []byte `db:"metadata"`
}

type Connection struct {
	ChannelID    string `db:"channel_id"`
	ChannelOwner string `db:"channel_owner"`
	ThingID      string `db:"thing_id"`
	ThingOwner   string `db:"thing_owner"`
}

type ConnectionsPage struct {
	mf13things.PageMetadata
	Connections []Connection
}

// dbRetrieveThings retrieves things from the database with the given page navigation parameters
func dbRetrieveThings(ctx context.Context, db mf13postgres.Database, pm mf13things.PageMetadata) (mf13things.Page, error) {
	q := `SELECT id, owner, name, key, metadata FROM things LIMIT :limit OFFSET :offset;`

	params := map[string]interface{}{
		"limit":  pm.Limit,
		"offset": pm.Offset,
	}

	rows, err := db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return mf13things.Page{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}
	defer rows.Close()

	var items []mf13things.Thing
	for rows.Next() {
		dbth := dbThing{}
		if err := rows.StructScan(&dbth); err != nil {
			return mf13things.Page{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
		}

		th, err := toThing(dbth)
		if err != nil {
			return mf13things.Page{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
		}

		items = append(items, th)
	}

	cq := `SELECT COUNT(*) FROM things;`

	total, err := total(ctx, db, cq, params)
	if err != nil {
		return mf13things.Page{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}

	page := mf13things.Page{
		Things: items,
		PageMetadata: mf13things.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}

	return page, nil
}

// dbRetrieveChannels retrieves things from the database with the given page navigation parameters
func dbRetrieveChannels(ctx context.Context, db mf13postgres.Database, pm mf13things.PageMetadata) (mf13things.ChannelsPage, error) {
	q := `SELECT id, owner, name, metadata FROM channels LIMIT :limit OFFSET :offset;`

	params := map[string]interface{}{
		"limit":  pm.Limit,
		"offset": pm.Offset,
	}

	rows, err := db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return mf13things.ChannelsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}
	defer rows.Close()

	var items []mf13things.Channel
	for rows.Next() {
		dbch := dbChannel{}
		if err := rows.StructScan(&dbch); err != nil {
			return mf13things.ChannelsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
		}

		ch, err := toChannel(dbch)
		if err != nil {
			return mf13things.ChannelsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
		}

		items = append(items, ch)
	}

	cq := `SELECT COUNT(*) FROM channels;`

	total, err := total(ctx, db, cq, params)
	if err != nil {
		return mf13things.ChannelsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}

	page := mf13things.ChannelsPage{
		Channels: items,
		PageMetadata: mf13things.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}

	return page, nil
}

// dbRetrieveConnections retrieves things from the database with the given page navigation parameters
func dbRetrieveConnections(ctx context.Context, db mf13postgres.Database, pm mf13things.PageMetadata) (ConnectionsPage, error) {
	q := `SELECT channel_id, channel_owner, thing_id, thing_owner FROM connections LIMIT :limit OFFSET :offset;`

	params := map[string]interface{}{
		"limit":  pm.Limit,
		"offset": pm.Offset,
	}

	rows, err := db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return ConnectionsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}
	defer rows.Close()

	var items []Connection
	for rows.Next() {
		dbconn := Connection{}
		if err := rows.StructScan(&dbconn); err != nil {
			return ConnectionsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
		}

		conn := Connection(dbconn)
		items = append(items, conn)
	}

	cq := `SELECT COUNT(*) FROM connections;`

	total, err := total(ctx, db, cq, params)
	if err != nil {
		return ConnectionsPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}

	page := ConnectionsPage{
		Connections: items,
		PageMetadata: mf13things.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}

	return page, nil
}

func total(ctx context.Context, db mf13postgres.Database, query string, params interface{}) (uint64, error) {
	rows, err := db.NamedQueryContext(ctx, query, params)
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

func toThing(dbth dbThing) (mf13things.Thing, error) {
	var metadata map[string]interface{}
	if dbth.Metadata != nil {
		if err := json.Unmarshal([]byte(dbth.Metadata), &metadata); err != nil {
			return mf13things.Thing{}, mf13errors.Wrap(mf13errors.ErrMalformedEntity, err)
		}
	}
	return mf13things.Thing{
		ID:       dbth.ID,
		Owner:    dbth.Owner,
		Name:     dbth.Name,
		Key:      dbth.Key,
		Metadata: metadata,
	}, nil
}

func toChannel(dch dbChannel) (mf13things.Channel, error) {
	var metadata map[string]interface{}
	if dch.Metadata != nil {
		if err := json.Unmarshal([]byte(dch.Metadata), &metadata); err != nil {
			return mf13things.Channel{}, mf13errors.Wrap(mf13errors.ErrMalformedEntity, err)
		}
	}

	return mf13things.Channel{
		ID:       dch.ID,
		Owner:    dch.Owner,
		Name:     dch.Name,
		Metadata: metadata,
	}, nil
}
