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

var icpCmd = &cobra.Command{
	Use:   common.CommandIcp + " [domain]",
	Args:  cobra.ExactArgs(1),
	Short: "Check ICP status",
	RunE: func(_ *cobra.Command, args []string) error {
		if (icpCmdGlobalVar.flags.Account == "" || icpCmdGlobalVar.flags.Key == "") && rootConfig != "" {
			v := common.Config(rootConfig, common.ICP)
			err := Encoder.JSONMarshaler(v, &icpCmdGlobalVar.flags)
			if err != nil {
				return err
			}
		}
		if icpCmdGlobalVar.flags.Account == "" || icpCmdGlobalVar.flags.Key == "" {
			return common.ErrInvalidToken
		}
		icpCmdGlobalVar.flags.domain = args[0]
		if err := icpCmdGlobalVar.Request(); err != nil {
			return err
		}
		OutputDefaultYAML(icpCmdGlobalVar)
		return nil
	},
	Example: common.Examples(`# Print the ICP status
-a account -k api_key google.com`, common.CommandIcp),
}

var icpCmdGlobalVar ICPResponse

func init() {
	rootCmd.AddCommand(icpCmd)

	icpCmd.Flags().StringVarP(&icpCmdGlobalVar.flags.Account, "account", "a", "", common.Usage("Enter the WEST account"))
	icpCmd.Flags().StringVarP(&icpCmdGlobalVar.flags.Key, "key", "k", "", common.Usage("Enter the WEST api key"))
	icpCmd.MarkFlagsRequiredTogether("account", "key")
}

type IcpFlags struct {
	Account string `json:"account"`
	Key     string `json:"api_key"`
	domain  string
}

type ICPResponse struct {
	DomainName string `json:"domain,omitempty" yaml:"domain,omitempty"`
	ICPCode    string `json:"icp,omitempty" yaml:"icp,omitempty"`
	ICPStatus  string `json:"icpstatus,omitempty" yaml:"icpstatus,omitempty"`

	flags IcpFlags
}

func (i *ICPResponse) requestURI() (string, error) {
	/* MD5 Hash */
	hashData := i.flags.Account + i.flags.Key + "domainname"
	sig, err := Hasher.Hash(common.HashAlgorithm(common.HashMd5), hashData)
	if err != nil {
		return "", err
	}
	rawCmd := fmt.Sprintf("domainname\r\ncheck\r\nentityname:icp\r\ndomains:%s\r\n.\r\n", i.flags.domain)
	/* URL Encoding */
	strCmd := url.QueryEscape(rawCmd)
	return fmt.Sprintf(`http://api.west263.com/api/?userid=%s&strCmd=%s&versig=%s`, i.flags.Account, strCmd, sig), nil
}

func (i *ICPResponse) Request() error {
	uri, err := i.requestURI()
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
	return Encoder.JSONMarshaler([]byte(match), i)
}
