package migrate

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	mf13log "github.com/mainflux/mainflux/logger"
	mf14sdk "github.com/mainflux/mainflux/pkg/sdk/go/0140"
	mf13thingspostgres "github.com/mainflux/mainflux/things/postgres"
	mf13userspostgres "github.com/mainflux/mainflux/users/postgres"
	"github.com/mainflux/migrations"
	"github.com/mainflux/migrations/migrate/things"
	"github.com/mainflux/migrations/migrate/users"
)

const (
	version13 = "0.13.0"
	version14 = "0.14.0"
	importOp  = "import"
	exportOp  = "export"
)

func Migrate(cfg migrations.Config, logger mf13log.Logger) {
	switch cfg.Operation {
	case importOp:
		switch cfg.ToVersion {
		case version14:
			Import14(cfg, logger)
		}
	case exportOp:
		switch cfg.FromVersion {
		case version13:
			Export13(cfg, logger)
		}
	}
}

func Export13(cfg migrations.Config, logger mf13log.Logger) {
	logger.Info(fmt.Sprintf("starting export from version %s", version13))

	usersDB := connectToUsersDB(cfg.UsersConfig.DBConfig)
	defer usersDB.Close()

	usersDatabase := mf13userspostgres.NewDatabase(usersDB)
	logger.Debug("connected to users database")

	thingsDB := connectToThingsDB(cfg.ThingsConfig.DBConfig)
	defer thingsDB.Close()

	thingsDatabase := mf13thingspostgres.NewDatabase(thingsDB)
	logger.Debug("connected to things database")

	us, err := users.MFRetrieveUsers(context.Background(), usersDatabase)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
	logger.Debug("retrieved users from database")
	if err := users.UsersToCSV(cfg.UsersConfig.UsersCSVPath, us.Users); err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
	logger.Debug("written users to csv file")

	ths, err := things.MFRetrieveThings(context.Background(), thingsDatabase)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
	logger.Debug("retrieved things from database")
	if err := things.ThingsToCSV(cfg.ThingsConfig.ThingsCSVPath, ths.Things); err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
	logger.Debug("written things to csv file")
	channels, err := things.MFRetrieveChannels(context.Background(), thingsDatabase)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
	logger.Debug("retrieved channels from database")
	if err := things.ChannelsToCSV(cfg.ThingsConfig.ChannelsCSVPath, channels.Channels); err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
	logger.Debug("written channels to csv file")
	connections, err := things.MFRetrieveConnections(context.Background(), thingsDatabase)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
	logger.Debug("retrieved connections from database")
	if err := things.ConnectionsToCSV(cfg.ThingsConfig.ConnectionsCSVPath, connections.Connections); err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
	logger.Debug("written connections to csv file")
	logger.Info(fmt.Sprintf("finished exporting from version %s", version13))
}

func Import14(cfg migrations.Config, logger mf13log.Logger) {
	logger.Info(fmt.Sprintf("starting importing to version %s", version14))
	sdkConf := mf14sdk.Config{
		ThingsURL:       cfg.ThingsURL,
		UsersURL:        cfg.UsersURL,
		MsgContentType:  mf14sdk.CTJSONSenML,
		TLSVerification: false,
	}

	sdk := mf14sdk.NewSDK(sdkConf)
	user := mf14sdk.User{
		Credentials: mf14sdk.Credentials{
			Identity: cfg.UserIdentity,
			Secret:   cfg.UserSecret,
		},
	}
	token, err := sdk.CreateToken(user)
	if err != nil {
		log.Panic(fmt.Errorf("failed to create token: %v", err))
	}
	logger.Debug("created user token")
	if err := things.CreateThings(sdk, cfg.UsersConfig.UsersCSVPath, cfg.ThingsConfig.ThingsCSVPath, token.AccessToken); err != nil {
		log.Panic(err)
	}
	logger.Debug("created things")
	if err := things.CreateChannels(sdk, cfg.UsersConfig.UsersCSVPath, cfg.ThingsConfig.ChannelsCSVPath, token.AccessToken); err != nil {
		log.Panic(err)
	}
	logger.Debug("created channels")
	if err := things.CreateConnections(sdk, cfg.ThingsConfig.ConnectionsCSVPath, token.AccessToken); err != nil {
		log.Panic(err)
	}
	logger.Debug("created connections")
	logger.Info(fmt.Sprintf("finished importing to version %s", version14))
}

func connectToThingsDB(dbConfig mf13thingspostgres.Config) *sqlx.DB {
	db, err := mf13thingspostgres.Connect(dbConfig)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to connect to things postgres: %s", err))
		os.Exit(1)
	}
	return db
}

func connectToUsersDB(dbConfig mf13userspostgres.Config) *sqlx.DB {
	db, err := mf13userspostgres.Connect(dbConfig)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to connect to users postgres: %s", err))
		os.Exit(1)
	}
	return db
}
