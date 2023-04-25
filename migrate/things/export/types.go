package export

type metadata map[string]interface{}

type pageMetadata struct {
	Total  uint64
	Offset uint64 `json:"offset,omitempty"`
	Limit  uint64 `json:"limit,omitempty"`
}

// thing is a thing entity.
// this is for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0.
type thing struct {
	ID       string
	Owner    string
	Name     string
	Key      string
	Metadata metadata
}

// thingsPage is a page of things.
// this is for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0.
type thingsPage struct {
	pageMetadata
	Things []thing
}

// dbThing is a thing entity in the database.
// this is for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0.
type dbThing struct {
	ID       string `db:"id"`
	Owner    string `db:"owner"`
	Name     string `db:"name"`
	Key      string `db:"key"`
	Metadata []byte `db:"metadata"`
}

// channel is a channel entity.
// this is for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0.
type channel struct {
	ID       string
	Owner    string
	Name     string
	Metadata map[string]interface{}
}

// channelsPage is a page of channels.
// this is for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0.
type channelsPage struct {
	pageMetadata
	Channels []channel
}

// dbChannel is a channel entity in the database.
// this is for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0.
type dbChannel struct {
	ID       string `db:"id"`
	Owner    string `db:"owner"`
	Name     string `db:"name"`
	Metadata []byte `db:"metadata"`
}

// connection is a connection entity.
// this is for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0.
type connection struct {
	ChannelID    string `db:"channel_id"`
	ChannelOwner string `db:"channel_owner"`
	ThingID      string `db:"thing_id"`
	ThingOwner   string `db:"thing_owner"`
}

// connectionsPage is a page of connections.
// this is for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0.
type connectionsPage struct {
	pageMetadata
	Connections []connection
}
