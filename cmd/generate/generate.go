package generate

import (
	"os"

	"github.com/autarcenergy/weather-data-scraper/internal/parser"
	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	defaultExportPath = "files/output.csv"
	defaultImportPath = "files/results.html"
)

func NewCmd() *cobra.Command {
	var importPath, exportPath string
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate and migrate the results into a database",
		Run: func(cmd *cobra.Command, args []string) {
			html, err := os.ReadFile(importPath)
			if err != nil {
				log.WithError(err).Fatal("cannot read file")
			}

			records := parser.MustParse(string(html))

			file, err := os.Create(exportPath)
			defer func() {
				if err := file.Close(); err != nil {
					log.WithError(err).Error("failed to close csv file handle")
				}
			}()
			if err != nil {
				log.WithError(err).Fatal("failed to open file")
			}

			err = gocsv.Marshal(&records, file)
			if err != nil {
				log.WithError(err).Error("failed to write to csv")
			}
		},
	}
	cmd.Flags().StringVar(&exportPath, "exportPath", defaultExportPath, "Default export path")
	cmd.Flags().StringVar(&importPath, "importPath", defaultImportPath, "Default import path")

	return cmd
}
