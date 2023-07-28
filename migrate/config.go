package migrate

import (
	"fmt"
	"log"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/mainflux/mainflux"
	thingspostgres "github.com/mainflux/mainflux/things/postgres" // for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0
	userspostgres "github.com/mainflux/mainflux/users/postgres"   // for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0
	"github.com/mainflux/migrations"
)

const (
	defLogLevel                 = "info"
	defThingsDBHost             = "localhost"
	defThingsDBPort             = "5432"
	defThingsDBUser             = "mainflux"
	defThingsDBPass             = "mainflux"
	defThingsDB                 = "things"
	defThingsDBSSLMode          = "disable"
	defThingsDBSSLCert          = ""
	defThingsDBSSLKey           = ""
	defThingsDBSSLRootCert      = ""
	defThingsDBContainerName    = "mainflux-things-db"
	defThingsDBContainerNetwork = "docker_mainflux-base-net"
	defThingsCSVPath            = "csv/things.csv"
	defChannelsCSVPath          = "csv/channels.csv"
	defConnectionCSVPath        = "csv/connections.csv"
	defUsersDBHost              = "localhost"
	defUsersDBPort              = "5432"
	defUsersDBUser              = "mainflux"
	defUsersDBPass              = "mainflux"
	defUsersDB                  = "users"
	defUsersDBSSLMode           = "disable"
	defUsersDBSSLCert           = ""
	defUsersDBSSLKey            = ""
	defUsersDBSSLRootCert       = ""
	defUsersDBContainerName     = "mainflux-users-db"
	defUsersDBContainerNetwork  = "docker_mainflux-base-net"
	defUsersCSVPath             = "csv/users.csv"
	defUsersURL                 = "http://localhost"
	defThingsURL                = "http://localhost"
	defUserIdentity             = "admin@example.com"
	defUserSecret               = "12345678"

	envLogLevel                 = "MF_LOG_LEVEL"
	envThingsDBHost             = "MF_THINGS_DB_HOST"
	envThingsDBPort             = "MF_THINGS_DB_PORT"
	envThingsDBUser             = "MF_THINGS_DB_USER"
	envThingsDBPass             = "MF_THINGS_DB_PASS"
	envThingsDB                 = "MF_THINGS_DB"
	envThingsDBSSLMode          = "MF_THINGS_DB_SSL_MODE"
	envThingsDBSSLCert          = "MF_THINGS_DB_SSL_CERT"
	envThingsDBSSLKey           = "MF_THINGS_DB_SSL_KEY"
	envThingsDBSSLRootCert      = "MF_THINGS_DB_SSL_ROOT_CERT"
	envThingsDBContainerName    = "MF_THINGS_CONTAINER_NAME"
	envThingsDBContainerNetwork = "MF_THINGS_CONTAINER_NETWORK"
	envThingsCSVPath            = "MF_THINGS_CSV_PATH"
	envChannelsCSVPath          = "MF_CHANNELS_CSV_PATH"
	envConnectionCSVPath        = "MF_CONNECTIONS_CSV_PATH"
	envUsersDBHost              = "MF_USERS_DB_HOST"
	envUsersDBPort              = "MF_USERS_DB_PORT"
	envUsersDBUser              = "MF_USERS_DB_USER"
	envUsersDBPass              = "MF_USERS_DB_PASS"
	envUsersDB                  = "MF_USERS_DB"
	envUsersDBSSLMode           = "MF_USERS_DB_SSL_MODE"
	envUsersDBSSLCert           = "MF_USERS_DB_SSL_CERT"
	envUsersDBSSLKey            = "MF_USERS_DB_SSL_KEY"
	envUsersDBSSLRootCert       = "MF_USERS_DB_SSL_ROOT_CERT"
	envUsersDBContainerName     = "MF_USERS_CONTAINER_NAME"
	envUsersDBContainerNetwork  = "MF_USERS_CONTAINER_NETWORK"
	envUsersCSVPath             = "MF_USERS_CSV_PATH"
	envUsersURL                 = "MF_USERS_URL"
	envThingsURL                = "MF_THINGS_URL"
	envUserIdentity             = "MF_USER_IDENTITY"
	envUserSecret               = "MF_USER_SECRET"
)

func LoadConfig() migrations.Config {
	thingsContainerName := mainflux.Env(envThingsDBContainerName, defThingsDBContainerName)
	thingsContainerNetwork := mainflux.Env(envThingsDBContainerNetwork, defThingsDBContainerNetwork)

	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("failed to create docker client with error %v", err)
	}
	imgs, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		log.Fatalf("failed to list containers with error %v", err)
	}
	thingsDBHost := mainflux.Env(envThingsDBHost, defThingsDBHost)
	for _, img := range imgs {
		if img.Names[0] == fmt.Sprintf("/%s", thingsContainerName) {
			thingsDBHost = img.Networks.Networks[thingsContainerNetwork].IPAddress
		}
	}

	usersContainerName := mainflux.Env(envUsersDBContainerName, defUsersDBContainerName)
	usersContainerNetwork := mainflux.Env(envUsersDBContainerNetwork, defUsersDBContainerNetwork)

	client, err = docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("failed to create docker client with error %v", err)
	}
	imgs, err = client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		log.Fatalf("failed to list containers with error %v", err)
	}
	usersDBHost := mainflux.Env(envUsersDBHost, defUsersDBHost)
	for _, img := range imgs {
		if img.Names[0] == fmt.Sprintf("/%s", usersContainerName) {
			usersDBHost = img.Networks.Networks[usersContainerNetwork].IPAddress
		}
	}

	tdbConfig := thingspostgres.Config{
		Host:        thingsDBHost,
		Port:        mainflux.Env(envThingsDBPort, defThingsDBPort),
		User:        mainflux.Env(envThingsDBUser, defThingsDBUser),
		Pass:        mainflux.Env(envThingsDBPass, defThingsDBPass),
		Name:        mainflux.Env(envThingsDB, defThingsDB),
		SSLMode:     mainflux.Env(envThingsDBSSLMode, defThingsDBSSLMode),
		SSLCert:     mainflux.Env(envThingsDBSSLCert, defThingsDBSSLCert),
		SSLKey:      mainflux.Env(envThingsDBSSLKey, defThingsDBSSLKey),
		SSLRootCert: mainflux.Env(envThingsDBSSLRootCert, defThingsDBSSLRootCert),
	}

	udbConfig := userspostgres.Config{
		Host:        usersDBHost,
		Port:        mainflux.Env(envUsersDBPort, defUsersDBPort),
		User:        mainflux.Env(envUsersDBUser, defUsersDBUser),
		Pass:        mainflux.Env(envUsersDBPass, defUsersDBPass),
		Name:        mainflux.Env(envUsersDB, defUsersDB),
		SSLMode:     mainflux.Env(envUsersDBSSLMode, defUsersDBSSLMode),
		SSLCert:     mainflux.Env(envUsersDBSSLCert, defUsersDBSSLCert),
		SSLKey:      mainflux.Env(envUsersDBSSLKey, defUsersDBSSLKey),
		SSLRootCert: mainflux.Env(envUsersDBSSLRootCert, defUsersDBSSLRootCert),
	}

	thConfig := migrations.ThingsConfig{
		DBConfig:           tdbConfig,
		ThingsCSVPath:      mainflux.Env(envThingsCSVPath, defThingsCSVPath),
		ChannelsCSVPath:    mainflux.Env(envChannelsCSVPath, defChannelsCSVPath),
		ConnectionsCSVPath: mainflux.Env(envConnectionCSVPath, defConnectionCSVPath),
	}

	uConfig := migrations.UsersConfig{
		DBConfig:     udbConfig,
		UsersCSVPath: mainflux.Env(envUsersCSVPath, defUsersCSVPath),
	}

	return migrations.Config{
		LogLevel:     mainflux.Env(envLogLevel, defLogLevel),
		ThingsConfig: thConfig,
		UsersConfig:  uConfig,
		UsersURL:     mainflux.Env(envUsersURL, defUsersURL),
		ThingsURL:    mainflux.Env(envThingsURL, defThingsURL),
		UserIdentity: mainflux.Env(envUserIdentity, defUserIdentity),
		UserSecret:   mainflux.Env(envUserSecret, defUserSecret),
	}
}
