package migrations

import (
	mf13thingspostgres "github.com/mainflux/mainflux/things/postgres"
	mf13userspostgres "github.com/mainflux/mainflux/users/postgres"
)

type ThingsConfig struct {
	DBConfig           mf13thingspostgres.Config
	ThingsCSVPath      string
	ChannelsCSVPath    string
	ConnectionsCSVPath string
}

type UsersConfig struct {
	DBConfig     mf13userspostgres.Config
	UsersCSVPath string
}

type Config struct {
	LogLevel     string
	FromVersion  string
	ToVersion    string
	Operation    string
	ThingsConfig ThingsConfig
	UsersConfig  UsersConfig
	UsersURL     string
	ThingsURL    string
	UserIdentity string
	UserSecret   string
}
