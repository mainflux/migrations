package migrate

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	mainflux0130 "github.com/mainflux/mainflux"
	thingsPostgres "github.com/mainflux/mainflux/things/postgres"
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
	envUsersURL                 = "MF_USERS_URL"
	envThingsURL                = "MF_THINGS_URL"
	envUserIdentity             = "MF_USER_IDENTITY"
	envUserSecret               = "MF_USER_SECRET"
)

func LoadConfig() migrations.Config {
	containerName := mainflux0130.Env(envThingsDBContainerName, defThingsDBContainerName)
	containerNetwork := mainflux0130.Env(envThingsDBContainerNetwork, defThingsDBContainerNetwork)

	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(fmt.Errorf("failed to create docker client with error %v", err))
	}
	imgs, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		panic(fmt.Errorf("failed to list containers with error %v", err))
	}
	thingsDBHost := mainflux0130.Env(envThingsDBHost, defThingsDBHost)
	for _, img := range imgs {
		if img.Names[0] == fmt.Sprintf("/%s", containerName) {
			thingsDBHost = img.Networks.Networks[containerNetwork].IPAddress
		}
	}

	dbConfig := thingsPostgres.Config{
		Host:        thingsDBHost,
		Port:        mainflux0130.Env(envThingsDBPort, defThingsDBPort),
		User:        mainflux0130.Env(envThingsDBUser, defThingsDBUser),
		Pass:        mainflux0130.Env(envThingsDBPass, defThingsDBPass),
		Name:        mainflux0130.Env(envThingsDB, defThingsDB),
		SSLMode:     mainflux0130.Env(envThingsDBSSLMode, defThingsDBSSLMode),
		SSLCert:     mainflux0130.Env(envThingsDBSSLCert, defThingsDBSSLCert),
		SSLKey:      mainflux0130.Env(envThingsDBSSLKey, defThingsDBSSLKey),
		SSLRootCert: mainflux0130.Env(envThingsDBSSLRootCert, defThingsDBSSLRootCert),
	}

	thConfig := migrations.ThingsConfig{
		DBConfig:           dbConfig,
		ThingsCSVPath:      mainflux0130.Env(envThingsCSVPath, defThingsCSVPath),
		ChannelsCSVPath:    mainflux0130.Env(envChannelsCSVPath, defChannelsCSVPath),
		ConnectionsCSVPath: mainflux0130.Env(envConnectionCSVPath, defConnectionCSVPath),
	}

	return migrations.Config{
		LogLevel:     mainflux0130.Env(envLogLevel, defLogLevel),
		ThingsConfig: thConfig,
		UsersURL:     mainflux0130.Env(envUsersURL, defUsersURL),
		ThingsURL:    mainflux0130.Env(envThingsURL, defThingsURL),
		UserIdentity: mainflux0130.Env(envUserIdentity, defUserIdentity),
		UserSecret:   mainflux0130.Env(envUserSecret, defUserSecret),
	}
}
