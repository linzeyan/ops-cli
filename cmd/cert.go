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
	"log"
	"time"

	"github.com/linzeyan/tlsCheck"
	"github.com/spf13/cobra"
)

// certCmd represents the cert command
var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Check tls cert",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, _ []string) {
		if host != "domain:port" {
			outputByHost()
			return
		} else if crt != "" {
			outputByFile()
			return
		}
		cmd.Help()
	},
	Example: `ops-cli cert -d www.google.com:443
ops-cli cert -d www.google.com:443 -i -n -p
ops-cli cert -f example.com.crt
ops-cli cert -f example.com.crt -i -n`,
}

var crt, host string
var ip, expiry, dns, issuer bool

func init() {
	rootCmd.AddCommand(certCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// certCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// certCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	certCmd.Flags().StringVarP(&crt, "file", "f", "", "Specify .crt file path")
	certCmd.Flags().StringVarP(&host, "domain", "d", "domain:port", "Specify domain and host port")
	certCmd.MarkFlagsMutuallyExclusive("file", "domain")
	certCmd.Flags().BoolVarP(&ip, "ip", "p", false, "Print IP")
	certCmd.Flags().BoolVarP(&expiry, "expiry", "e", true, "Print expiry time")
	certCmd.Flags().BoolVarP(&dns, "dns", "n", false, "Print DNS names")
	certCmd.Flags().BoolVarP(&issuer, "issuer", "i", false, "Print issuer")
}

func outputByHost() {
	conn, err := tlsCheck.CheckByHost(host)
	if err != nil {
		log.Println(err)
		return
	}
	cert := conn.ConnectionState().PeerCertificates[0]
	if expiry {
		fmt.Println("Expiry time:", cert.NotAfter.Local().Format(time.RFC3339))
	}
	if ip {
		fmt.Println("Server IP:", conn.RemoteAddr().String())
	}
	if dns {
		fmt.Println("DNS:", cert.DNSNames)
	}
	if issuer {
		fmt.Println("Issuer:", cert.Issuer.String())
	}
}

func outputByFile() {
	cert, err := tlsCheck.CheckByFile(crt)
	if err != nil {
		log.Println(err)
		return
	}
	if cert == nil {
		log.Println(cert)
		return
	}
	if expiry {
		fmt.Println("Expiry time:", cert[0].NotAfter.Local().Format(time.RFC3339))
	}
	if dns {
		fmt.Println("DNS:", cert[0].DNSNames)
	}
	if issuer {
		fmt.Println("Issuer:", cert[0].Issuer.String())
	}
}
