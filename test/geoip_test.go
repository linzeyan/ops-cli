package test_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

const geoipSingle = `{
  "continent": "Oceania",
  "country": "Australia",
  "countryCode": "AU",
  "regionName": "Queensland",
  "city": "South Brisbane",
  "district": "",
  "timezone": "Australia/Brisbane",
  "currency": "AUD",
  "isp": "Cloudflare, Inc",
  "org": "APNIC and Cloudflare DNS Resolver project",
  "as": "AS13335 Cloudflare, Inc.",
  "asname": "CLOUDFLARENET",
  "mobile": false,
  "proxy": false,
  "hosting": true,
  "query": "1.1.1.1"
}

`

const geoipBatch = `[
  {
    "continent": "Oceania",
    "country": "Australia",
    "countryCode": "AU",
    "regionName": "Queensland",
    "city": "South Brisbane",
    "district": "",
    "timezone": "Australia/Brisbane",
    "currency": "AUD",
    "isp": "Cloudflare, Inc",
    "org": "APNIC and Cloudflare DNS Resolver project",
    "as": "AS13335 Cloudflare, Inc.",
    "asname": "CLOUDFLARENET",
    "mobile": false,
    "proxy": false,
    "hosting": true,
    "query": "1.1.1.1"
  },
  {
    "continent": "North America",
    "country": "United States",
    "countryCode": "US",
    "regionName": "Virginia",
    "city": "Ashburn",
    "district": "",
    "timezone": "America/New_York",
    "currency": "USD",
    "isp": "Google LLC",
    "org": "Google Public DNS",
    "as": "AS15169 Google LLC",
    "asname": "GOOGLE",
    "mobile": false,
    "proxy": false,
    "hosting": true,
    "query": "8.8.8.8"
  }
]

`

func TestGeoip(t *testing.T) {
	const subCommand = "geoip"
	testCases := []struct {
		input    []string
		expected string
	}{
		// {[]string{runCommand, mainGo, subCommand, "999.999.1.1"}, ""},
		{[]string{runCommand, mainGo, subCommand, "1.1.1.1"}, geoipSingle},
		{[]string{runCommand, mainGo, subCommand, "1.1.1.1", "8.8.8.8"}, geoipBatch},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			got, err := exec.Command(mainCommand, testCases[i].input...).Output()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, testCases[i].expected, string(got))
		})
	}
}
