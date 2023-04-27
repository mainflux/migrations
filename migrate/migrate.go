package migrate

import (
	"context"
	"fmt"
	"log"
	"net/http"

	// For profiling.
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/jmoiron/sqlx"
	"github.com/julz/prettyprogress"
	"github.com/julz/prettyprogress/ui"
	"github.com/mainflux/mainflux/logger"
	thingspostgres "github.com/mainflux/mainflux/things/postgres" // for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0
	userspostgres "github.com/mainflux/mainflux/users/postgres"   // for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0
	"github.com/mainflux/migrations"
	exportthings "github.com/mainflux/migrations/migrate/things/export"
	importthings "github.com/mainflux/migrations/migrate/things/import"
	exportusers "github.com/mainflux/migrations/migrate/users/export"
)

const (
	importOp        = "import"
	exportOp        = "export"
	refreshDuration = 200 * time.Millisecond

	version10 = "0.10.0"
	version11 = "0.11.0"
	version12 = "0.12.0"
	version13 = "0.13.0"
	version14 = "0.14.0"
)

type retrieveQueries struct {
	users       string
	things      string
	channels    string
	connections string
}

var (
	version10Queries = retrieveQueries{
		users:       "SELECT email, password, metadata FROM users LIMIT :limit OFFSET :offset;",
		things:      "SELECT id, owner, name, key, metadata FROM things LIMIT :limit OFFSET :offset;",
		channels:    "SELECT id, owner, name, metadata FROM channels LIMIT :limit OFFSET :offset;",
		connections: "SELECT channel_id, channel_owner, thing_id, thing_owner FROM connections LIMIT :limit OFFSET :offset;",
	}

	version11Queries = version10Queries // same as version 10

	version12Queries = retrieveQueries{
		users:       "SELECT id, email, password, metadata FROM users LIMIT :limit OFFSET :offset;",
		things:      version10Queries.things,
		channels:    version10Queries.channels,
		connections: version10Queries.connections,
	}

	version13Queries = version12Queries // same as version 12
)

var retrieveSQLQueries = map[string]retrieveQueries{
	version10: version10Queries,
	version11: version11Queries,
	version12: version12Queries,
	version13: version13Queries,
}

func Migrate(ctx context.Context, cfg migrations.Config, logger logger.Logger) {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	switch cfg.Operation {
	case importOp:
		switch cfg.ToVersion {
		case version14:
			Import(ctx, cfg, logger)
		default:
			logger.Error("we only support importing to version 0.14.0")
		}

	case exportOp:
		if sqlStatements, ok := retrieveSQLQueries[cfg.FromVersion]; ok {
			cfg.UsersRetrievalSQL = sqlStatements.users
			cfg.ThingsRetrievalSQL = sqlStatements.things
			cfg.ChannelsRetrievalSQL = sqlStatements.channels
			cfg.ConnectionsRetrievalSQL = sqlStatements.connections
			Export(ctx, cfg, logger)
		}
	}
}

func Export(ctx context.Context, cfg migrations.Config, logger logger.Logger) {
	logger.Info(fmt.Sprintf("starting export from version %s", cfg.FromVersion))

	usersDB := connectToUsersDB(cfg.UsersConfig.DBConfig)
	defer usersDB.Close()

	usersDatabase := userspostgres.NewDatabase(usersDB)
	logger.Debug("connected to users database")

	thingsDB := connectToThingsDB(cfg.ThingsConfig.DBConfig)
	defer thingsDB.Close()

	thingsDatabase := thingspostgres.NewDatabase(thingsDB)
	logger.Debug("connected to things database")

	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	bullets := ui.AnimatedBulletSet
	bullets.Running = bullets.Running.WithColor(color.New(color.FgGreen))
	multiStep := prettyprogress.NewFancyMultistep(
		writer,
		prettyprogress.WithAnimationFrameTicker(time.NewTicker(refreshDuration).C),
		prettyprogress.WithBullets(
			bullets,
		),
	)

	usersStep := multiStep.AddStep("Retrieving Users", 0)
	thingsStep := multiStep.AddStep("Retrieving Things", 0)
	channelsStep := multiStep.AddStep("Retrieving Channels", 0)
	connStep := multiStep.AddStep("Retrieving Connections", 0)

	var wg sync.WaitGroup
	var numOfJobs = 4
	wg.Add(numOfJobs)

	go func() {
		defer wg.Done()

		if err := exportusers.RetrieveAndWriteUsers(ctx, cfg.UsersRetrievalSQL, usersDatabase, cfg.UsersConfig.UsersCSVPath); err != nil {
			logger.Error(err.Error())
		}
		usersStep.Complete("Finished Retrieveing Users")
	}()

	go func() {
		defer wg.Done()

		if err := exportthings.RetrieveAndWriteThings(ctx, cfg.ThingsRetrievalSQL, thingsDatabase, cfg.ThingsConfig.ThingsCSVPath); err != nil {
			logger.Error(err.Error())
		}
		thingsStep.Complete("Finished Retrieveing Things")
	}()

	go func() {
		defer wg.Done()

		if err := exportthings.RetrieveAndWriteChannels(ctx, cfg.ChannelsRetrievalSQL, thingsDatabase, cfg.ThingsConfig.ChannelsCSVPath); err != nil {
			logger.Error(err.Error())
		}
		channelsStep.Complete("Finished Retrieveing Channels")
	}()

	go func() {
		defer wg.Done()

		if err := exportthings.RetrieveAndWriteConnections(ctx, cfg.ConnectionsRetrievalSQL, thingsDatabase, cfg.ThingsConfig.ConnectionsCSVPath); err != nil {
			logger.Error(err.Error())
		}
		connStep.Complete("Finished Retrieveing Connection")
	}()

	wg.Wait()
	logger.Info(fmt.Sprintf("finished exporting from version %s", version13))
}

func Import(ctx context.Context, cfg migrations.Config, logger logger.Logger) {
	logger.Info(fmt.Sprintf("starting importing to version %s", cfg.ToVersion))

	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	bullets := ui.AnimatedBulletSet
	bullets.Running = bullets.Running.WithColor(color.New(color.FgGreen))
	multiStep := prettyprogress.NewFancyMultistep(
		writer,
		prettyprogress.WithAnimationFrameTicker(time.NewTicker(refreshDuration).C),
		prettyprogress.WithBullets(
			bullets,
		),
	)

	thingsStep := multiStep.AddStep("creating things", 0)
	channelsStep := multiStep.AddStep("creating channels", 0)
	connStep := multiStep.AddStep("creating connections", 0)

	sdk, token, err := importthings.InitSDK(cfg)
	if err != nil {
		logger.Error(err.Error())
	}

	if err := importthings.ReadAndCreateThings(ctx, sdk, cfg.UsersConfig.UsersCSVPath, cfg.ThingsConfig.ThingsCSVPath, token); err != nil {
		logger.Error(err.Error())
	}
	thingsStep.Complete("finished creating things")

	if err := importthings.ReadAndCreateChannels(ctx, sdk, cfg.UsersConfig.UsersCSVPath, cfg.ThingsConfig.ChannelsCSVPath, token); err != nil {
		logger.Error(err.Error())
	}
	channelsStep.Complete("finished creating channels")

	if err := importthings.ReadAndCreateConnections(ctx, sdk, cfg.ThingsConfig.ConnectionsCSVPath, token); err != nil {
		logger.Error(err.Error())
	}
	connStep.Complete("finished creating connections")

	logger.Info(fmt.Sprintf("finished importing to version %s", cfg.ToVersion))
}

func connectToThingsDB(dbConfig thingspostgres.Config) *sqlx.DB {
	db, err := thingspostgres.Connect(dbConfig)
	if err != nil {
		log.Fatalf("failed to connect to things postgres: %s", err)
	}

	return db
}

func connectToUsersDB(dbConfig userspostgres.Config) *sqlx.DB {
	db, err := userspostgres.Connect(dbConfig)
	if err != nil {
		log.Fatalf("failed to connect to users postgres: %s", err)
	}

	return db
}
