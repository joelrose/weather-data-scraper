package main

import (
	"os"

	"github.com/anaskhan96/soup"
	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
)

const (
	inputFileName  = "files/results.html"
	outputFileName = "output.csv"
)

func initLogger() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

type Record struct {
	Postleitzahl           string `csv:"postleitzahl"`
	Jahresmitteltemperatur string `csv:"jahresmitteltemperatur"`
	NormAußentemperatur    string `csv:"normAußentemperatur"`
	Hoehe                  string `csv:"hoehe"`
	Klimazone              string `csv:"klimazone"`
	Ort                    string `csv:"ort"`
}

func main() {
	initLogger()

	html, err := os.ReadFile(inputFileName)
	if err != nil {
		log.WithError(err).Error("cannot read file")
	}

	doc := soup.HTMLParse(string(html))
	elements := doc.FindAll("polygon")

	records := []Record{}
	for _, element := range elements {
		attrs := element.Attrs()

		records = append(records, Record{
			Postleitzahl:           attrs["zip"],
			Jahresmitteltemperatur: attrs["aat"],
			NormAußentemperatur:    attrs["dot"],
			Hoehe:                  attrs["alt"],
			Klimazone:              attrs["zone"],
			Ort:                    attrs["place"],
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
}
