package migrate

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	mf13log "github.com/mainflux/mainflux/logger"
	mf14sdk "github.com/mainflux/mainflux/pkg/sdk/go/0140"
	mf13postgres "github.com/mainflux/mainflux/things/postgres"
	"github.com/mainflux/migrations"
	"github.com/mainflux/migrations/migrate/things"
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
	db := connectToThingsDB(cfg.ThingsConfig.DBConfig)
	defer db.Close()

	database := mf13postgres.NewDatabase(db)
	logger.Debug("connected to things database")

	ths, err := things.MFRetrieveThings(context.Background(), database)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
	logger.Debug("retrieved things from database")
	if err := things.ThingsToCSV(cfg.ThingsConfig.ThingsCSVPath, ths.Things); err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
	logger.Debug("written things to csv file")
	channels, err := things.MFRetrieveChannels(context.Background(), database)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
	logger.Debug("retrieved channels from database")
	if err := things.ChannelsToCSV(cfg.ThingsConfig.ChannelsCSVPath, channels.Channels); err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
	logger.Debug("written channels to csv file")
	connections, err := things.MFRetrieveConnections(context.Background(), database)
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
	if err := things.CreateThings(sdk, cfg.ThingsConfig.ThingsCSVPath, token.AccessToken); err != nil {
		log.Panic(err)
	}
	logger.Debug("created things")
	if err := things.CreateChannels(sdk, cfg.ThingsConfig.ChannelsCSVPath, token.AccessToken); err != nil {
		log.Panic(err)
	}
	logger.Debug("created channels")
	if err := things.CreateConnections(sdk, cfg.ThingsConfig.ConnectionsCSVPath, token.AccessToken); err != nil {
		log.Panic(err)
	}
	logger.Debug("created connections")
	logger.Info(fmt.Sprintf("finished importing to version %s", version14))
}

func connectToThingsDB(dbConfig mf13postgres.Config) *sqlx.DB {
	db, err := mf13postgres.Connect(dbConfig)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to connect to things postgres: %s", err))
		os.Exit(1)
	}
	return db
}
