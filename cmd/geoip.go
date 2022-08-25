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
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

var geoipCmd = &cobra.Command{
	Use:   common.CommandGeoip + " IP...",
	Short: "Print IP geographic information",
	Run: func(_ *cobra.Command, args []string) {
		var out any
		var err error
		switch len(args) {
		case 0:
			var resp []byte
			resp, err = common.HTTPRequestContent("https://myexternalip.com/raw", nil)
			out = map[string]string{"ip": string(resp)}
		case 1:
			var r geoIPSingle
			out, err = r.Request(args[0])
		default:
			var r geoIPBatch
			out, err = r.Request(args)
		}
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		OutputDefaultJSON(out)
	},
	Example: common.Examples(`# Print IP geographic information
1.1.1.1

# Print multiple IP geographic informations
1.1.1.1 8.8.8.8`, common.CommandGeoip),
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(geoipCmd)
}

type geoIPSingle struct {
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

func (geoIPSingle) Request(geoipInput string) (*geoIPSingle, error) {
	/* Valid IP */
	if !validator.ValidIP(geoipInput) {
		return nil, common.ErrInvalidIP
	}
	apiURL := fmt.Sprintf("http://ip-api.com/json/%s?fields=continent,countryCode,country,regionName,city,district,query,isp,org,as,asname,currency,timezone,mobile,proxy,hosting", geoipInput)

	content, err := common.HTTPRequestContent(apiURL, nil)
	if err != nil {
		return nil, err
	}
	var data geoIPSingle
	if err = Encoder.JSONMarshaler(content, &data); err != nil {
		return nil, err
	}
	return &data, err
}

type geoIPBatch []geoIPSingle

func (geoIPBatch) Request(geoipBatch []string) (*geoIPBatch, error) {
	var ips = `[`
	/* Valid IP and combine args */
	for i := range geoipBatch {
		switch {
		case validator.ValidIP(geoipBatch[i]):
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
	content, err := common.HTTPRequestContent(apiURL, strings.NewReader(ips), http.MethodPost)
	if err != nil {
		return nil, err
	}
	var data geoIPBatch
	if err = Encoder.JSONMarshaler(content, &data); err != nil {
		return nil, err
	}
	return &data, err
}
