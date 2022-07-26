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
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var geoipCmd = &cobra.Command{
	Use:   "geoip",
	Short: "Print IP geographic information",
	Run: func(cmd *cobra.Command, args []string) {
		var out rootOutput
		var err error
		if len(args) == 1 {
			var r GeoIPSingle
			out, err = r.Request(args[0])
			if err != nil {
				log.Println(err)
				return
			}
			outputDefaultJson(out)
			return
		}
		if len(args) > 1 {
			var r GeoIPBatch
			out, err = r.Request(args)
			if err != nil {
				log.Println(err)
				return
			}
			outputDefaultJson(out)
			return
		}
		cmd.Help()
	},
	Example: Examples(`# Print IP geographic information
ops-cli geoip 1.1.1.1

# Print multiple IP geographic information
ops-cli geoip 1.1.1.1 8.8.8.8`),
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

func (g GeoIPSingle) Json() {
	out, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

func (g GeoIPSingle) Yaml() {
	out, err := yaml.Marshal(g)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

func (g GeoIPSingle) String() {
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

func (GeoIPSingle) Request(geoipInput string) (*GeoIPSingle, error) {
	apiUrl := fmt.Sprintf("http://ip-api.com/json/%s?fields=continent,countryCode,country,regionName,city,district,query,isp,org,as,asname,currency,timezone,mobile,proxy,hosting", geoipInput)

	var client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
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
		content, err := ioutil.ReadAll(resp.Body)
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

func (g GeoIPBatch) Json() {
	out, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

func (g GeoIPBatch) Yaml() {
	out, err := yaml.Marshal(g)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

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
		content, err := ioutil.ReadAll(resp.Body)
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