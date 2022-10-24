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
	"net/http"
	"net/url"
	"regexp"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func initICP() *cobra.Command {
	var flags struct {
		Account string `json:"account"`
		Key     string `json:"api_key"`
	}
	var icpCmd = &cobra.Command{
		Use:  CommandICP + " domain",
		Args: cobra.ExactArgs(1),
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Check ICP status",
		RunE: func(_ *cobra.Command, args []string) error {
			if rootConfig != "" {
				if err := ReadConfig(CommandICP, &flags); err != nil {
					return err
				}
			}
			if flags.Account == "" || flags.Key == "" {
				return common.ErrInvalidToken
			}

			var i ICP
			if err := i.Request(flags.Account, flags.Key, args[0]); err != nil {
				return err
			}
			printer.Printf(printer.SetYamlAsDefaultFormat(rootOutputFormat), i)
			return nil
		},
		Example: common.Examples(`# Print the ICP status
-a account -k api_key google.com`, CommandICP),
	}

	icpCmd.Flags().StringVarP(&flags.Account, "account", "a", "", common.Usage("Enter the WEST account"))
	icpCmd.Flags().StringVarP(&flags.Key, "key", "k", "", common.Usage("Enter the WEST api key"))
	icpCmd.MarkFlagsRequiredTogether("account", "key")
	return icpCmd
}

type ICP struct {
	DomainName string `json:"domain,omitempty" yaml:"domain,omitempty"`
	ICPCode    string `json:"icp,omitempty" yaml:"icp,omitempty"`
	ICPStatus  string `json:"icpstatus,omitempty" yaml:"icpstatus,omitempty"`
}

func (i *ICP) requestURI(account, key, domain string) (string, error) {
	/* MD5 Hash */
	hashData := account + key + "domainname"
	sig, err := Hasher.Hash(HashAlgorithm(HashMd5), hashData)
	if err != nil {
		return "", err
	}
	rawCmd := fmt.Sprintf("domainname\r\ncheck\r\nentityname:icp\r\ndomains:%s\r\n.\r\n", domain)
	/* URL Encoding */
	strCmd := url.QueryEscape(rawCmd)
	return fmt.Sprintf(`http://api.west263.com/api/?userid=%s&strCmd=%s&versig=%s`, account, strCmd, sig), nil
}

func (i *ICP) Request(account, key, domain string) error {
	uri, err := i.requestURI(account, key, domain)
	if err != nil {
		return err
	}
	resp, err := common.HTTPRequestContentGB18030(uri, nil, http.MethodPost)
	if err != nil {
		return err
	}
	/* Find String */
	re := regexp.MustCompile("{.*}")
	match := fmt.Sprintln(re.FindString(string(resp)))
	return Encoder.JSONMarshaler(match, i)
}
