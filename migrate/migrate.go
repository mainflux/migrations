package migrate

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go/0140"
	thingsPostgres "github.com/mainflux/mainflux/things/postgres"
	"github.com/mainflux/migrations"
	"github.com/mainflux/migrations/migrate/things"
)

const (
	version13 = "0.13.0"
	version14 = "0.14.0"
	importOp  = "import"
	exportOp  = "export"
)

func Migrate(cfg migrations.Config) {
	switch cfg.Operation {
	case importOp:
		switch cfg.ToVersion {
		case version14:
			Import14(cfg)
		}
	case exportOp:
		switch cfg.FromVersion {
		case version13:
			Export13(cfg)
		}
	}
}

func Export13(cfg migrations.Config) {
	db := connectToThingsDB(cfg.ThingsConfig.DBConfig)
	defer db.Close()

	database := thingsPostgres.NewDatabase(db)

	ths, err := things.MFRetrieveThings(context.Background(), database)
	if err != nil {
		log.Fatal(err)
	}
	if err := things.ThingsToCSV(cfg.ThingsConfig.ThingsCSVPath, ths.Things); err != nil {
		log.Fatal(err)
	}
	channels, err := things.MFRetrieveChannels(context.Background(), database)
	if err != nil {
		log.Fatal(err)
	}
	if err := things.ChannelsToCSV(cfg.ThingsConfig.ChannelsCSVPath, channels.Channels); err != nil {
		log.Fatal(err)
	}
	connections, err := things.MFRetrieveConnections(context.Background(), database)
	if err != nil {
		log.Fatal(err)
	}
	if err := things.ConnectionsToCSV(cfg.ThingsConfig.ConnectionsCSVPath, connections.Connections); err != nil {
		log.Fatal(err)
	}
}

func Import14(cfg migrations.Config) {
	sdkConf := mfsdk.Config{
		ThingsURL:       cfg.ThingsURL,
		UsersURL:        cfg.UsersURL,
		MsgContentType:  mfsdk.CTJSONSenML,
		TLSVerification: false,
	}

	sdk := mfsdk.NewSDK(sdkConf)
	user := mfsdk.User{
		Credentials: mfsdk.Credentials{
			Identity: cfg.UserIdentity,
			Secret:   cfg.UserSecret,
		},
	}
	token, err := sdk.CreateToken(user)
	if err != nil {
		log.Panic(fmt.Errorf("failed to create token: %v", err))
	}
	if err := things.CreateThings(sdk, cfg.ThingsConfig.ThingsCSVPath, token.AccessToken); err != nil {
		log.Panic(err)
	}
	if err := things.CreateChannels(sdk, cfg.ThingsConfig.ChannelsCSVPath, token.AccessToken); err != nil {
		log.Panic(err)
	}
	if err := things.CreateConnections(sdk, cfg.ThingsConfig.ConnectionsCSVPath, token.AccessToken); err != nil {
		log.Panic(err)
	}
}

func connectToThingsDB(dbConfig thingsPostgres.Config) *sqlx.DB {
	db, err := thingsPostgres.Connect(dbConfig)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to connect to things postgres: %s", err))
		os.Exit(1)
	}
	return db
}
