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
				t.Error(testCases[i].input, err)
			}
			assert.Equal(t, testCases[i].expected, string(got))
		})
	}
}

func TestBinaryGeoip(t *testing.T) {
	const subCommand = "geoip"
	host := "8.8.8.8"
	args := []string{"-j", "-y"}

	for i := range args {
		t.Run(args[i], func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, host, args[i]).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
	t.Run("batch", func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand, "1.2.3.4", "1.1.1.1", "8.8.4.4").Run(); err != nil {
			t.Error(err)
		}
	})
}
