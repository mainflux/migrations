package migrations

import (
	mf13thingspostgres "github.com/mainflux/mainflux/things/postgres"
	mf13userspostgres "github.com/mainflux/mainflux/users/postgres"
)

type ThingsConfig13 struct {
	DBConfig           mf13thingspostgres.Config
	ThingsCSVPath      string
	ChannelsCSVPath    string
	ConnectionsCSVPath string
}

type UsersConfig13 struct {
	DBConfig     mf13userspostgres.Config
	UsersCSVPath string
}

type Config struct {
	LogLevel       string
	FromVersion    string
	ToVersion      string
	Operation      string
	ThingsConfig13 ThingsConfig13
	UsersConfig13  UsersConfig13
	UsersURL       string
	ThingsURL      string
	UserIdentity   string
	UserSecret     string
}
