package parser

import (
	"reflect"
	"testing"
)

func TestMustParse(t *testing.T) {
	cases := []struct {
		input    string
		expected []Record
	}{
		{
			`<polygon zip="12345" aat="10 °C" dot="11 °C" alt="100 m" zone="1" place="Place 1" />`,
			[]Record{
				{
					ZIP:                    "12345",
					MeanAnnualTemperature:  10,
					NormOutsideTemperature: 11,
					Height:                 100,
					Zone:                   1,
					Place:                  "Place 1",
				},
			},
		},
		{
			`<polygon zip="54321" aat="-11 °C" dot="-12 °C" alt="-200 m" zone="-2" place="Place 2" />`,
			[]Record{
				{
					ZIP:                    "54321",
					MeanAnnualTemperature:  -11,
					NormOutsideTemperature: -12,
					Height:                 -200,
					Zone:                   -2,
					Place:                  "Place 2",
				},
			},
		},
	}

	for _, c := range cases {
		actual := MustParse(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("MustParse(%q) == %v, want %v", c.input, actual, c.expected)
		}
	}
}

func TestCleanCelsius(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"-10 °C", "-10"},
		{"-10.5 °C", "-10.5"},
		{"-11", "-11"},
		{"-11.5", "-11.5"},
		{"", ""},
		{"12", "12"},
		{"12.5", "12.5"},
		{"13 °C", "13"},
		{"13.5 °C", "13.5"},
	}

	for _, c := range cases {
		actual := cleanCelsius(c.input)
		if actual != c.expected {
			t.Errorf("cleanCelsius(%q) == %q, want %q", c.input, actual, c.expected)
		}
	}
}

func TestCleanMeter(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"-98 m", "-98"},
		{"-99", "-99"},
		{"", ""},
		{"100", "100"},
		{"101 m", "101"},
	}

	for _, c := range cases {
		actual := cleanMeter(c.input)
		if actual != c.expected {
			t.Errorf("cleanMeter(%q) == %q, want %q", c.input, actual, c.expected)
		}
	}
}
