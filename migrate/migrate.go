package migrate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"

	"github.com/gosuri/uilive"
	"github.com/jmoiron/sqlx"
	"github.com/julz/prettyprogress"
	mf13log "github.com/mainflux/mainflux/logger"
	mf14sdk "github.com/mainflux/mainflux/pkg/sdk/go/0140"
	mf13thingspostgres "github.com/mainflux/mainflux/things/postgres"
	mf13userspostgres "github.com/mainflux/mainflux/users/postgres"
	"github.com/mainflux/migrations"
	things13 "github.com/mainflux/migrations/migrate/things/v13"
	things14 "github.com/mainflux/migrations/migrate/things/v14"
	users13 "github.com/mainflux/migrations/migrate/users/v13"
)

const (
	version13 = "0.13.0"
	version14 = "0.14.0"
	importOp  = "import"
	exportOp  = "export"
)

func Migrate(cfg migrations.Config, logger mf13log.Logger) {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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

	usersDB := connectToUsersDB(cfg.UsersConfig13.DBConfig)
	defer usersDB.Close()

	usersDatabase := mf13userspostgres.NewDatabase(usersDB)
	logger.Debug("connected to users database")

	thingsDB := connectToThingsDB(cfg.ThingsConfig13.DBConfig)
	defer thingsDB.Close()

	thingsDatabase := mf13thingspostgres.NewDatabase(thingsDB)
	logger.Debug("connected to things database")

	w := uilive.New()
	w.Start()
	defer w.Stop()

	multiStep := prettyprogress.NewFancyMultistep(w)

	usersStep := multiStep.AddStep("Retrieving Users", 0)
	thingsStep := multiStep.AddStep("Retrieving Things", 0)
	channelsStep := multiStep.AddStep("Retrieving Channels", 0)
	connStep := multiStep.AddStep("Retrieving Connections", 0)

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		if err := users13.RetrieveAndWriteUsers(context.Background(), usersDatabase, cfg.UsersConfig13.UsersCSVPath); err != nil {
			logger.Error(err.Error())
		}
		usersStep.Complete("Finished Retrieveing Users")
		defer wg.Done()
	}()

	go func() {
		if err := things13.RetrieveAndWriteThings(context.Background(), thingsDatabase, cfg.ThingsConfig13.ThingsCSVPath); err != nil {
			logger.Error(err.Error())
		}
		thingsStep.Complete("Finished Retrieveing Things")
		defer wg.Done()
	}()

	go func() {
		if err := things13.RetrieveAndWriteChannels(context.Background(), thingsDatabase, cfg.ThingsConfig13.ChannelsCSVPath); err != nil {
			logger.Error(err.Error())
		}
		channelsStep.Complete("Finished Retrieveing Channels")
		defer wg.Done()
	}()

	go func() {
		if err := things13.RetrieveAndWriteConnections(context.Background(), thingsDatabase, cfg.ThingsConfig13.ConnectionsCSVPath); err != nil {
			logger.Error(err.Error())
		}
		connStep.Complete("Finished Retrieveing Connection")
		defer wg.Done()
	}()

	wg.Wait()
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
		logger.Error(fmt.Sprintf("failed to create token: %v", err))
	}
	logger.Debug("created user token")
	if err := things14.CreateThings(sdk, cfg.UsersConfig13.UsersCSVPath, cfg.ThingsConfig13.ThingsCSVPath, token.AccessToken); err != nil {
		logger.Error(err.Error())
	}
	logger.Debug("created things")
	if err := things14.CreateChannels(sdk, cfg.UsersConfig13.UsersCSVPath, cfg.ThingsConfig13.ChannelsCSVPath, token.AccessToken); err != nil {
		logger.Error(err.Error())
	}
	logger.Debug("created channels")
	if err := things14.CreateConnections(sdk, cfg.ThingsConfig13.ConnectionsCSVPath, token.AccessToken); err != nil {
		logger.Error(err.Error())
	}
	logger.Debug("created connections")
	logger.Info(fmt.Sprintf("finished importing to version %s", version14))
}

func connectToThingsDB(dbConfig mf13thingspostgres.Config) *sqlx.DB {
	db, err := mf13thingspostgres.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to things postgres: %s", err)
	}
	return db
}

func connectToUsersDB(dbConfig mf13userspostgres.Config) *sqlx.DB {
	db, err := mf13userspostgres.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to users postgres: %s", err)
	}
	return db
}
