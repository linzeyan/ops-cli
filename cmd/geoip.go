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
	"github.com/spf13/cobra"
)

func initGeoip() *cobra.Command {
	var flags struct {
		country  bool
		region   bool
		city     bool
		timezone bool
		isp      bool
		org      bool
		as       bool
	}
	var geoipCmd = &cobra.Command{
		GroupID: getGroupID(CommandGeoip),
		Use:     CommandGeoip + " IP...",
		Short:   "Print IP geographic information",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(_ *cobra.Command, args []string) {
			switch l := len(args); l {
			case 0:
				var resp []byte
				resp, err := common.HTTPRequestContent("https://myexternalip.com/raw")
				if err != nil {
					logger.Info(err.Error())
					return
				}
				printer.Printf("%s", resp)
			default:
				var r GeoIP
				out, err := r.Request(args)
				if err != nil {
					logger.Info(err.Error())
					return
				}
				outFormat := printer.SetJSONAsDefaultFormat(rootOutputFormat)
				switch l {
				case 1:
					if flags.country {
						printer.Printf(outFormat, out[0].Country)
					}
					if flags.region {
						printer.Printf(outFormat, out[0].RegionName)
					}
					if flags.city {
						printer.Printf(outFormat, out[0].City)
					}
					if flags.timezone {
						printer.Printf(outFormat, out[0].Timezone)
					}
					if flags.isp {
						printer.Printf(outFormat, out[0].ISP)
					}
					if flags.org {
						printer.Printf(outFormat, out[0].Org)
					}
					if flags.as {
						printer.Printf(outFormat, out[0].As)
					}
					if flags.as || flags.city || flags.country || flags.isp || flags.org || flags.region || flags.timezone {
						return
					}
					printer.Printf(outFormat, out[0])
				default:
					printer.Printf(outFormat, out)
				}
			}
		},
		Example: common.Examples(`# Print IP geographic information
1.1.1.1

# Print multiple IP geographic informations
1.1.1.1 8.8.8.8`, CommandGeoip),
		DisableFlagsInUseLine: true,
	}
	geoipCmd.Flags().BoolVar(&flags.country, "country", false, common.Usage("Print country"))
	geoipCmd.Flags().BoolVar(&flags.region, "region", false, common.Usage("Print region"))
	geoipCmd.Flags().BoolVar(&flags.city, "city", false, common.Usage("Print city"))
	geoipCmd.Flags().BoolVar(&flags.timezone, "timezone", false, common.Usage("Print timezone"))
	geoipCmd.Flags().BoolVar(&flags.isp, "isp", false, common.Usage("Print ISP"))
	geoipCmd.Flags().BoolVar(&flags.org, "org", false, common.Usage("Print organization"))
	geoipCmd.Flags().BoolVar(&flags.as, "as", false, common.Usage("Print AS number"))
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

func (GeoIP) Request(inputs []string) ([]GeoIP, error) {
	var ips = `[`
	/* Valid IP and combine args */
	for i := range inputs {
		switch {
		case common.IsIP(inputs[i]):
			ips += fmt.Sprintf(`"%s", `, inputs[i])
		default:
			ip, err := net.LookupIP(inputs[i])
			if err != nil {
				logger.Debug(err.Error())
				return nil, err
			}
			ips += fmt.Sprintf(`"%s", `, ip[0])
		}
	}
	ips = strings.TrimRight(ips, `, `) + `]`

	apiURL := "http://ip-api.com/batch?fields=continent,countryCode,country,regionName,city,district,query,isp,org,as,asname,currency,timezone,mobile,proxy,hosting"
	content, err := common.HTTPRequestContent(apiURL, common.HTTPConfig{Method: http.MethodPost, Body: ips})
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	var data []GeoIP
	if err = Encoder.JSONMarshaler(content, &data); err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return data, err
}
