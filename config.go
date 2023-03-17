package migrations

import thingsPostgres "github.com/mainflux/mainflux/things/postgres"

type ThingsConfig struct {
	DBConfig           thingsPostgres.Config
	ThingsCSVPath      string
	ChannelsCSVPath    string
	ConnectionsCSVPath string
}

type Config struct {
	FromVersion  string
	ToVersion    string
	Operation    string
	ThingsConfig ThingsConfig
	UsersURL     string
	ThingsURL    string
	UserIdentity string
	UserSecret   string
}
