package things

import (
	"context"
	"database/sql/driver"
	"encoding/json"

	mferrors "github.com/mainflux/mainflux/pkg/errors"
	mfthings "github.com/mainflux/mainflux/things"
	thingsPostgres "github.com/mainflux/mainflux/things/postgres"
)

type dbThing struct {
	ID       string `db:"id"`
	Owner    string `db:"owner"`
	Name     string `db:"name"`
	Key      string `db:"key"`
	Metadata []byte `db:"metadata"`
}

// dbMetadata type for handling metadata properly in database/sql.
type dbMetadata map[string]interface{}

type dbChannel struct {
	ID       string     `db:"id"`
	Owner    string     `db:"owner"`
	Name     string     `db:"name"`
	Metadata dbMetadata `db:"metadata"`
}

type Connection struct {
	ChannelID    string
	ChannelOwner string
	ThingID      string
	ThingOwner   string
}

type ConnectionsPage struct {
	mfthings.PageMetadata
	Connections []Connection
}

type dbConnection struct {
	ChannelID    string `db:"channel_id"`
	ChannelOwner string `db:"channel_owner"`
	ThingID      string `db:"thing_id"`
	ThingOwner   string `db:"thing_owner"`
}

func RetrieveAllThings(ctx context.Context, db thingsPostgres.Database, pm mfthings.PageMetadata) (mfthings.Page, error) {
	q := `SELECT id, owner, name, key, metadata FROM things LIMIT :limit OFFSET :offset;`

	params := map[string]interface{}{
		"limit":  pm.Limit,
		"offset": pm.Offset,
	}

	rows, err := db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return mfthings.Page{}, mferrors.Wrap(mferrors.ErrViewEntity, err)
	}
	defer rows.Close()

	var items []mfthings.Thing
	for rows.Next() {
		dbth := dbThing{}
		if err := rows.StructScan(&dbth); err != nil {
			return mfthings.Page{}, mferrors.Wrap(mferrors.ErrViewEntity, err)
		}

		th, err := toThing(dbth)
		if err != nil {
			return mfthings.Page{}, mferrors.Wrap(mferrors.ErrViewEntity, err)
		}

		items = append(items, th)
	}

	cq := `SELECT COUNT(*) FROM things;`

	total, err := total(ctx, db, cq, params)
	if err != nil {
		return mfthings.Page{}, mferrors.Wrap(mferrors.ErrViewEntity, err)
	}

	page := mfthings.Page{
		Things: items,
		PageMetadata: mfthings.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}

	return page, nil
}

func RetrieveAllChannels(ctx context.Context, db thingsPostgres.Database, pm mfthings.PageMetadata) (mfthings.ChannelsPage, error) {
	q := `SELECT id, owner, name, metadata FROM channels LIMIT :limit OFFSET :offset;`

	params := map[string]interface{}{
		"limit":  pm.Limit,
		"offset": pm.Offset,
	}

	rows, err := db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return mfthings.ChannelsPage{}, mferrors.Wrap(mferrors.ErrViewEntity, err)
	}
	defer rows.Close()

	var items []mfthings.Channel
	for rows.Next() {
		dbch := dbChannel{}
		if err := rows.StructScan(&dbch); err != nil {
			return mfthings.ChannelsPage{}, mferrors.Wrap(mferrors.ErrViewEntity, err)
		}

		ch := toChannel(dbch)
		items = append(items, ch)
	}

	cq := `SELECT COUNT(*) FROM channels;`

	total, err := total(ctx, db, cq, params)
	if err != nil {
		return mfthings.ChannelsPage{}, mferrors.Wrap(mferrors.ErrViewEntity, err)
	}

	page := mfthings.ChannelsPage{
		Channels: items,
		PageMetadata: mfthings.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}

	return page, nil
}

func RetrieveAllConnections(ctx context.Context, db thingsPostgres.Database, pm mfthings.PageMetadata) (ConnectionsPage, error) {
	q := `SELECT channel_id, channel_owner, thing_id, thing_owner FROM connections LIMIT :limit OFFSET :offset;`

	params := map[string]interface{}{
		"limit":  pm.Limit,
		"offset": pm.Offset,
	}

	rows, err := db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return ConnectionsPage{}, mferrors.Wrap(mferrors.ErrViewEntity, err)
	}
	defer rows.Close()

	var items []Connection
	for rows.Next() {
		dbconn := dbConnection{}
		if err := rows.StructScan(&dbconn); err != nil {
			return ConnectionsPage{}, mferrors.Wrap(mferrors.ErrViewEntity, err)
		}

		conn := Connection(dbconn)
		items = append(items, conn)
	}

	cq := `SELECT COUNT(*) FROM connections;`

	total, err := total(ctx, db, cq, params)
	if err != nil {
		return ConnectionsPage{}, mferrors.Wrap(mferrors.ErrViewEntity, err)
	}

	page := ConnectionsPage{
		Connections: items,
		PageMetadata: mfthings.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}

	return page, nil
}

func total(ctx context.Context, db thingsPostgres.Database, query string, params interface{}) (uint64, error) {
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

func toThing(dbth dbThing) (mfthings.Thing, error) {
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(dbth.Metadata), &metadata); err != nil {
		return mfthings.Thing{}, mferrors.Wrap(mferrors.ErrMalformedEntity, err)
	}

	return mfthings.Thing{
		ID:       dbth.ID,
		Owner:    dbth.Owner,
		Name:     dbth.Name,
		Key:      dbth.Key,
		Metadata: metadata,
	}, nil
}

// Scan implements the database/sql scanner interface.
// When interface is nil `m` is set to nil.
// If error occurs on casting data then m points to empty metadata.
func (m *dbMetadata) Scan(value interface{}) error {
	if value == nil {
		// m = nil
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		// m = &dbMetadata{}
		return mferrors.ErrScanMetadata
	}

	if err := json.Unmarshal(b, m); err != nil {
		return err
	}

	return nil
}

// Value implements database/sql valuer interface.
func (m dbMetadata) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, err
}

func toChannel(ch dbChannel) mfthings.Channel {
	return mfthings.Channel{
		ID:       ch.ID,
		Owner:    ch.Owner,
		Name:     ch.Name,
		Metadata: ch.Metadata,
	}
}
