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

// icpCmd represents the icp command
var icpCmd = &cobra.Command{
	Use:   "icp",
	Short: "Check ICP status",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// icpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// icpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	icpCmd.Flags().StringVarP(&icp.ConfigFile, "config", "c", "", "Specify config file")
	icpCmd.Flags().StringVarP(&icp.Domain, "domain", "d", "", "Specify domain name")
	icpCmd.MarkFlagRequired("domain")
}
