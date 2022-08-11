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
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

var geoipCmd = &cobra.Command{
	Use:   "geoip [IP1] [IP2]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Print IP geographic information",
	Run: func(_ *cobra.Command, args []string) {
		var out interface{}
		var err error
		switch len(args) {
		case 1:
			var r GeoIPSingle
			out, err = r.Request(args[0])
		default:
			var r GeoIPBatch
			out, err = r.Request(args)
		}
		if err != nil {
			log.Println(err)
			return
		}
		OutputDefaultJSON(out)
	},
	Example: Examples(`# Print IP geographic information
ops-cli geoip 1.1.1.1

# Print multiple IP geographic informations
ops-cli geoip 1.1.1.1 8.8.8.8`),
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(geoipCmd)
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

func (GeoIPSingle) Request(geoipInput string) (*GeoIPSingle, error) {
	/* Valid IP */
	if !ValidIP(geoipInput) {
		return nil, errors.New("not a valid IP")
	}
	apiURL := fmt.Sprintf("http://ip-api.com/json/%s?fields=continent,countryCode,country,regionName,city,district,query,isp,org,as,asname,currency,timezone,mobile,proxy,hosting", geoipInput)

	var client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	req, err := http.NewRequestWithContext(rootContext, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var data GeoIPSingle
		err = json.Unmarshal(content, &data)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}
	return nil, err
}

type GeoIPBatch []GeoIPSingle

func (g GeoIPBatch) JSON() { PrintJSON(g) }

func (g GeoIPBatch) YAML() { PrintYAML(g) }

func (g GeoIPBatch) String() {
	var s strings.Builder
	f := reflect.ValueOf(&g).Elem()
	t := f.Type()
	for i := 0; i < f.NumField(); i++ {
		_, err := s.WriteString(fmt.Sprintf("%-10s\t%v\n", t.Field(i).Name, f.Field(i).Interface()))
		if err != nil {
			log.Println(err)
			return
		}
	}
	fmt.Println(s.String())
}

func (GeoIPBatch) Request(geoipBatch []string) (*GeoIPBatch, error) {
	var ips = `[`
	/* Valid IP and combine args */
	for i := range geoipBatch {
		switch {
		case ValidIP(geoipBatch[i]):
			ips += fmt.Sprintf(`"%s", `, geoipBatch[i])
		default:
			ip, err := net.LookupIP(geoipBatch[i])
			if err != nil {
				return nil, err
			}
			ips += fmt.Sprintf(`"%s", `, ip[0])
		}
	}
	ips = strings.TrimRight(ips, `, `) + `]`

	apiURL := "http://ip-api.com/batch?fields=continent,countryCode,country,regionName,city,district,query,isp,org,as,asname,currency,timezone,mobile,proxy,hosting"
	var client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	req, err := http.NewRequestWithContext(rootContext, http.MethodPost, apiURL, strings.NewReader(ips))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var data GeoIPBatch
		err = json.Unmarshal(content, &data)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}
	return nil, err
}
