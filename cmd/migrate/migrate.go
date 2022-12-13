package migrate

import (
	"fmt"
	"os"

	"github.com/autarcenergy/weather-data-poc/internal/parser"
	"github.com/gocarina/gocsv"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	_ "github.com/lib/pq"
)

const (
	defaultConnectionString = "postgresql://username:password@localhost:5432/autarc?sslmode=disable"
	defaultImportPath       = "files/output.csv"
)

func NewCmd() *cobra.Command {
	var connStr, importPath string
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Insert the results into a database",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			db, err := sqlx.ConnectContext(ctx, "postgres", connStr)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %v", err)
			}
			defer db.Close()

			records, err := getRecords(importPath)
			if err != nil {
				return err
			}

			for _, r := range records {
				query := "INSERT INTO weather_data (zip, mean_annual_temperature, norm_outside_temperature, height, zone, place) VALUES (" + fmt.Sprintf("'%v', %v, %v, %v, %v, '%v'", r.ZIP, r.MeanAnnualTemperature, r.NormOutsideTemperature, r.Height, r.Zone, r.Place) + ");"
				db.MustExec(query)
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&connStr, "connStr", defaultConnectionString, "Database connection string")
	cmd.Flags().StringVar(&importPath, "path", defaultImportPath, "Input CSV file")

	return cmd
}

func getRecords(path string) ([]*parser.Record, error) {
	f, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("getRecords: cannot open file, %v", err)
	}
	defer f.Close()

	r := []*parser.Record{}
	if err := gocsv.UnmarshalFile(f, &r); err != nil { // Load clients from file
		return nil, fmt.Errorf("getRecords: failed to UnmarshalFile, %v", err)
	}

	return r, nil
}
