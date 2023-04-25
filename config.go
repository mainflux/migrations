package migrations

import (
	thingspostgres "github.com/mainflux/mainflux/things/postgres" // for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0
	userspostgres "github.com/mainflux/mainflux/users/postgres"   // for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0
)

type ThingsConfig struct {
	DBConfig           thingspostgres.Config
	ThingsCSVPath      string
	ChannelsCSVPath    string
	ConnectionsCSVPath string
}

type UsersConfig struct {
	DBConfig     userspostgres.Config
	UsersCSVPath string
}

type Config struct {
	LogLevel                string
	FromVersion             string
	ToVersion               string
	Operation               string
	ThingsConfig            ThingsConfig
	UsersConfig             UsersConfig
	UsersURL                string
	ThingsURL               string
	UserIdentity            string
	UserSecret              string
	UsersRetrievalSQL       string
	ThingsRetrievalSQL      string
	ChannelsRetrievalSQL    string
	ConnectionsRetrievalSQL string
}
