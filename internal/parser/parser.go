package parser

import (
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
	log "github.com/sirupsen/logrus"
)

type Record struct {
	ZIP                    string  `csv:"zip"`
	MeanAnnualTemperature  float32 `csv:"meanAnnualTemperature"`
	NormOutsideTemperature float32 `csv:"normOutsideTemperature"`
	Height                 int     `csv:"height"`
	Zone                   int     `csv:"zone"`
	Place                  string  `csv:"place"`
}

func MustParse(input string) []Record {
	doc := soup.HTMLParse(string(input))
	elements := doc.FindAll("polygon")

	records := []Record{}
	seenElements := map[string]bool{}
	for _, element := range elements {
		attrs := element.Attrs()

		record := mustCleanData(attrs)
		if seenElements[record.ZIP] {
			continue
		} else {
			seenElements[record.ZIP] = true
		}

		records = append(records, record)
	}

	return records
}

func mustCleanData(attrs map[string]string) Record {
	zip := attrs["zip"]

	jat, err := strconv.ParseFloat(cleanCelsius(attrs["aat"]), 32)
	if err != nil {
		log.WithError(err).Fatal("failed converting aat")
	}
	nat, err := strconv.ParseFloat(cleanCelsius(attrs["dot"]), 32)
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

	return Record{
		ZIP:                    zip,
		MeanAnnualTemperature:  float32(jat),
		NormOutsideTemperature: float32(nat),
		Height:                 height,
		Zone:                   zone,
		Place:                  place,
	}
}

func cleanCelsius(s string) string {
	return strings.Replace(s, " °C", "", 1)
}

func cleanMeter(s string) string {
	return strings.Replace(s, " m", "", 1)
}
