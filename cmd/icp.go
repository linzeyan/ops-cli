/*
Copyright © 2022 ZeYanLin <zeyanlin@outlook.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/linzeyan/icp"
	"github.com/spf13/cobra"
)

var icpCmd = &cobra.Command{
	Use:   "icp",
	Short: "Check ICP status",
	Args:  cobra.OnlyValidArgs,
	Run: func(_ *cobra.Command, _ []string) {
		icp.ReadConf()
		result := icp.Check()
		fmt.Println(icp.Domain, result)
	},
	Example: Examples(`# Print the ICP status of the domain
ops-cli icp -d apple.com

# Print the ICP status of the domain and specify the configuration path
ops-cli icp -c ~/.env -d baidu.com`),
}

func init() {
	rootCmd.AddCommand(icpCmd)

	icpCmd.Flags().StringVarP(&icp.ConfigFile, "config", "c", "", "Specify config file")
	icpCmd.Flags().StringVarP(&icp.Domain, "domain", "d", "", "Specify domain name")
	icpCmd.MarkFlagRequired("domain")
}
