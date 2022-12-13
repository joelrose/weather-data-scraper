package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
)

const (
	inputFileName  = "files/results.html"
	outputFileName = "files/output.csv"
)

func initLogger() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
}

type Record struct {
	Postleitzahl           string  `csv:"zip"`
	Jahresmitteltemperatur float32 `csv:"meanAnnualTemperature"`
	NormAußentemperatur    float32 `csv:"normOutsideTemperature"`
	Hoehe                  int     `csv:"height"`
	Klimazone              int     `csv:"zone"`
	Ort                    string  `csv:"place"`
}

func cleanCelsius(s string) string {
	return strings.Replace(s, " °C", "", 1)
}

func cleanMeter(s string) string {
	return strings.Replace(s, " m", "", 1)
}

func main() {
	initLogger()
	log.Info("Running parser")

	html, err := os.ReadFile(inputFileName)
	if err != nil {
		log.WithError(err).Error("cannot read file")
	}

	doc := soup.HTMLParse(string(html))
	elements := doc.FindAll("polygon")

	var sb strings.Builder
	sb.WriteString("INSERT INTO weather_data (zip, mean_annual_temperature, norm_outside_temperature, height, zone, place) VALUES")
	records := []Record{}

	seenElements := map[string]bool{}
	for _, element := range elements {
		attrs := element.Attrs()

		zip := attrs["zip"]
		if seenElements[zip] {
			continue
		} else {
			seenElements[zip] = true
		}

		jat, err := strconv.ParseFloat(cleanCelsius(attrs["aat"]), 16)
		if err != nil {
			log.WithError(err).Fatal("failed converting aat")
		}
		nat, err := strconv.ParseFloat(cleanCelsius(attrs["dot"]), 16)
		if err != nil {
			log.WithError(err).Fatal("failed converting dot")
		}
		height, err := strconv.Atoi(cleanMeter(attrs["alt"]))
		if err != nil {
			log.WithError(err).Fatal("failed converting alt")
		}
		zone, err := strconv.Atoi(attrs["zone"])
		if err != nil {
			log.WithError(err).Fatal("failed converting zone")
		}
		place := attrs["place"]

		sb.WriteString(fmt.Sprintf("('%v', %v, %v, %v, %v, '%v'),", zip, jat, nat, height, zone, place))

		records = append(records, Record{
			Postleitzahl:           zip,
			Jahresmitteltemperatur: float32(jat),
			NormAußentemperatur:    float32(nat),
			Hoehe:                  height,
			Klimazone:              zone,
			Ort:                    place,
		})
	}

	file, err := os.Create(outputFileName)
	defer func() {
		if err := file.Close(); err != nil {
			log.WithError(err).Error("failed to close csv file handle")
		}
	}()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	err = gocsv.Marshal(&records, file)
	if err != nil {
		log.WithError(err).Error("failed to write to csv")
	}

	sqlHandle, _ := os.Create("test.sql")

	sqlHandle.WriteString(sb.String())

	log.Infof("Parser executed successful, check the %s file for the results.", outputFileName)
}
