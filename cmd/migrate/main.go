package main

import (
	"log"

	"github.com/mainflux/migrations/migrate"
	"github.com/spf13/cobra"
)

func main() {
	cfg := migrate.LoadConfig()

	var rootCmd = &cobra.Command{
		Use:   "migrations",
		Short: "migrations is migration tool for Mainflux",
		Long: `Tool for migrating from one version of mainflux to another.It migrates things, channels and thier connections.
				Complete documentation is available at https://docs.mainflux.io`,
		Run: func(cmd *cobra.Command, args []string) {
			migrate.Migrate(cfg)
		},
	}

	// Root Flags
	rootCmd.PersistentFlags().StringVarP(&cfg.FromVersion, "fromversion", "f", "0.13.0", "mainflux version you want to migrate from")
	rootCmd.PersistentFlags().StringVarP(&cfg.ToVersion, "toversion", "t", "0.14.0", "mainflux version you want to migrate to")
	rootCmd.PersistentFlags().StringVarP(&cfg.Operation, "operation", "o", "export", "export dataor import data to a new mainflux deployment")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
