package cmd

import (
	"fmt"
	"os"

	"github.com/autarcenergy/weather-data-scraper/cmd/generate"
	"github.com/autarcenergy/weather-data-scraper/cmd/migrate"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func initLogger() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
}

func Execute() {
	initLogger()

	cmd := &cobra.Command{
		Use:   "app",
		Short: "",
	}

	cmd.AddCommand(migrate.NewCmd())
	cmd.AddCommand(generate.NewCmd())

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "There was an error while executing your CLI '%v'", err)
		os.Exit(1)
	}
}
