package migrate

import (
	"fmt"
	"log"

	docker "github.com/fsouza/go-dockerclient"
	mainflux13 "github.com/mainflux/mainflux"
	mf13thingspostgres "github.com/mainflux/mainflux/things/postgres"
	mf13userspostgres "github.com/mainflux/mainflux/users/postgres"
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
	thingsContainerName := mainflux13.Env(envThingsDBContainerName, defThingsDBContainerName)
	thingsContainerNetwork := mainflux13.Env(envThingsDBContainerNetwork, defThingsDBContainerNetwork)

	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("failed to create docker client with error %v", err)
	}
	imgs, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		log.Fatalf("failed to list containers with error %v", err)
	}
	thingsDBHost := mainflux13.Env(envThingsDBHost, defThingsDBHost)
	for _, img := range imgs {
		if img.Names[0] == fmt.Sprintf("/%s", thingsContainerName) {
			thingsDBHost = img.Networks.Networks[thingsContainerNetwork].IPAddress
		}
	}

	usersContainerName := mainflux13.Env(envUsersDBContainerName, defUsersDBContainerName)
	usersContainerNetwork := mainflux13.Env(envUsersDBContainerNetwork, defUsersDBContainerNetwork)

	client, err = docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("failed to create docker client with error %v", err)
	}
	imgs, err = client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		log.Fatalf("failed to list containers with error %v", err)
	}
	usersDBHost := mainflux13.Env(envUsersDBHost, defUsersDBHost)
	for _, img := range imgs {
		if img.Names[0] == fmt.Sprintf("/%s", usersContainerName) {
			usersDBHost = img.Networks.Networks[usersContainerNetwork].IPAddress
		}
	}

	tdbConfig := mf13thingspostgres.Config{
		Host:        thingsDBHost,
		Port:        mainflux13.Env(envThingsDBPort, defThingsDBPort),
		User:        mainflux13.Env(envThingsDBUser, defThingsDBUser),
		Pass:        mainflux13.Env(envThingsDBPass, defThingsDBPass),
		Name:        mainflux13.Env(envThingsDB, defThingsDB),
		SSLMode:     mainflux13.Env(envThingsDBSSLMode, defThingsDBSSLMode),
		SSLCert:     mainflux13.Env(envThingsDBSSLCert, defThingsDBSSLCert),
		SSLKey:      mainflux13.Env(envThingsDBSSLKey, defThingsDBSSLKey),
		SSLRootCert: mainflux13.Env(envThingsDBSSLRootCert, defThingsDBSSLRootCert),
	}

	udbConfig := mf13userspostgres.Config{
		Host:        usersDBHost,
		Port:        mainflux13.Env(envUsersDBPort, defUsersDBPort),
		User:        mainflux13.Env(envUsersDBUser, defUsersDBUser),
		Pass:        mainflux13.Env(envUsersDBPass, defUsersDBPass),
		Name:        mainflux13.Env(envUsersDB, defUsersDB),
		SSLMode:     mainflux13.Env(envUsersDBSSLMode, defUsersDBSSLMode),
		SSLCert:     mainflux13.Env(envUsersDBSSLCert, defUsersDBSSLCert),
		SSLKey:      mainflux13.Env(envUsersDBSSLKey, defUsersDBSSLKey),
		SSLRootCert: mainflux13.Env(envUsersDBSSLRootCert, defUsersDBSSLRootCert),
	}

	thConfig := migrations.ThingsConfig13{
		DBConfig:           tdbConfig,
		ThingsCSVPath:      mainflux13.Env(envThingsCSVPath, defThingsCSVPath),
		ChannelsCSVPath:    mainflux13.Env(envChannelsCSVPath, defChannelsCSVPath),
		ConnectionsCSVPath: mainflux13.Env(envConnectionCSVPath, defConnectionCSVPath),
	}

	uConfig := migrations.UsersConfig13{
		DBConfig:     udbConfig,
		UsersCSVPath: mainflux13.Env(envUsersCSVPath, defUsersCSVPath),
	}

	return migrations.Config{
		LogLevel:       mainflux13.Env(envLogLevel, defLogLevel),
		ThingsConfig13: thConfig,
		UsersConfig13:  uConfig,
		UsersURL:       mainflux13.Env(envUsersURL, defUsersURL),
		ThingsURL:      mainflux13.Env(envThingsURL, defThingsURL),
		UserIdentity:   mainflux13.Env(envUserIdentity, defUserIdentity),
		UserSecret:     mainflux13.Env(envUserSecret, defUserSecret),
	}
}
