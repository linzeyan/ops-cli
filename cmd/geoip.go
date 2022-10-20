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
	"net"
	"net/http"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

func initGeoip() *cobra.Command {
	var geoipCmd = &cobra.Command{
		Use:   CommandGeoip + " IP...",
		Short: "Print IP geographic information",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(_ *cobra.Command, args []string) error {
			var out any
			var err error
			switch len(args) {
			case 0:
				var resp []byte
				resp, err = common.HTTPRequestContent("https://myexternalip.com/raw")
				out = map[string]string{"ip": string(resp)}
			case 1:
				var r GeoIP
				out, err = r.Request(args[0])
			default:
				var r GeoIPList
				out, err = r.Request(args)
			}
			if err != nil {
				return err
			}
			OutputDefaultJSON(out)
			return err
		},
		Example: common.Examples(`# Print IP geographic information
1.1.1.1

# Print multiple IP geographic informations
1.1.1.1 8.8.8.8`, CommandGeoip),
		DisableFlagsInUseLine: true,
	}
	return geoipCmd
}

type GeoIP struct {
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

func (GeoIP) Request(ip string) (*GeoIP, error) {
	/* Valid IP */
	if !validator.IsIP(ip) {
		return nil, common.ErrInvalidIP
	}
	apiURL := fmt.Sprintf("http://ip-api.com/json/%s?fields=continent,countryCode,country,regionName,city,district,query,isp,org,as,asname,currency,timezone,mobile,proxy,hosting", ip)

	content, err := common.HTTPRequestContent(apiURL)
	if err != nil {
		return nil, err
	}
	var data GeoIP
	if err = Encoder.JSONMarshaler(content, &data); err != nil {
		return nil, err
	}
	return &data, err
}

type GeoIPList []GeoIP

func (GeoIPList) Request(inputs []string) (*GeoIPList, error) {
	var ips = `[`
	/* Valid IP and combine args */
	for i := range inputs {
		switch {
		case validator.IsIP(inputs[i]):
			ips += fmt.Sprintf(`"%s", `, inputs[i])
		default:
			ip, err := net.LookupIP(inputs[i])
			if err != nil {
				return nil, err
			}
			ips += fmt.Sprintf(`"%s", `, ip[0])
		}
	}
	ips = strings.TrimRight(ips, `, `) + `]`

	apiURL := "http://ip-api.com/batch?fields=continent,countryCode,country,regionName,city,district,query,isp,org,as,asname,currency,timezone,mobile,proxy,hosting"
	content, err := common.HTTPRequestContent(apiURL, common.HTTPConfig{Method: http.MethodPost, Body: ips})
	if err != nil {
		return nil, err
	}
	var data GeoIPList
	if err = Encoder.JSONMarshaler(content, &data); err != nil {
		return nil, err
	}
	return &data, err
}
