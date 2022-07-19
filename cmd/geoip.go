/*
Copyright © 2022 ZeYanLin <zeyanlin@outlook.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var geoipCmd = &cobra.Command{
	Use:   "geoip",
	Short: "Print IP geographic information",
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		if geoipInput != "" {
			geoipRequestSingle()
			return
		}
		if geoipBatch != nil {
			geoipRequestBatch()
			return
		}
		cmd.Help()
	},
	Example: Examples(`# Print IP geographic information
ops-cli geoip -s 1.1.1.1

# Print multiple IP geographic information
ops-cli geoip -b 1.1.1.1 -b 8.8.8.8`),
}

var geoipInput string
var geoipBatch []string

func init() {
	rootCmd.AddCommand(geoipCmd)

	geoipCmd.Flags().StringVarP(&geoipInput, "source", "s", "", "Specify IP or domain")
	geoipCmd.Flags().StringArrayVarP(&geoipBatch, "batch", "b", nil, "Enter multiple IPs or domains")
}

type GeoIPSingle struct {
	Continent   string `json:"continent"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	RegionName  string `json:"regionName"`
	City        string `json:"city"`
	District    string `json:"district"`
	Timezone    string `json:"timezone"`
	Currency    string `json:"currency"`
	ISP         string `json:"isp"`
	Org         string `json:"org"`
	As          string `json:"as"`
	Asname      string `json:"asname"`
	Mobile      bool   `json:"mobile"`
	Proxy       bool   `json:"proxy"`
	Hosting     bool   `json:"hosting"`
	Query       string `json:"query"`
}

func geoipRequestSingle() error {
	apiUrl := fmt.Sprintf("http://ip-api.com/json/%s?fields=continent,countryCode,country,regionName,city,district,query,isp,org,as,asname,currency,timezone,mobile,proxy,hosting", geoipInput)

	var client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var data GeoIPSingle
		err = json.Unmarshal(content, &data)
		if err != nil {
			return err
		}
		out, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	}
	return err
}

type GeoIPBatch []GeoIPSingle

func geoipRequestBatch() error {
	apiUrl := "http://ip-api.com/batch?fields=continent,countryCode,country,regionName,city,district,query,isp,org,as,asname,currency,timezone,mobile,proxy,hosting"

	var client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	var ips string = `[`
	for i := range geoipBatch {
		ips = ips + fmt.Sprintf(`"%s", `, geoipBatch[i])
	}
	ips = strings.TrimRight(ips, `, `) + `]`
	req, err := http.NewRequest(http.MethodPost, apiUrl, strings.NewReader(ips))
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var data GeoIPBatch
		err = json.Unmarshal(content, &data)
		if err != nil {
			return err
		}
		out, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	}
	return err
}
