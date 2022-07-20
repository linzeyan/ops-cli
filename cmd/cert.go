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
	"time"

	"github.com/linzeyan/tlsCheck"
	"github.com/spf13/cobra"
)

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Check tls cert",
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		if certHost != "domain:port" {
			certOutputByHost()
			return
		} else if certCrt != "" {
			certOutputByFile()
			return
		}
		cmd.Help()
	},
	Example: Examples(`# Print certificate expiration time
ops-cli cert -d www.google.com:443

# Print certificate expiration time, DNS, IP and issuer
ops-cli cert -d www.google.com:443 -i -n -p

# Print certificate expiration time
ops-cli cert -f example.com.crt

# Print certificate expiration time, DNS and issuer
ops-cli cert -f example.com.crt -i -n`),
}

var certCrt, certHost string
var certIP, certExpiry, certDNS, certIssuer bool

func init() {
	rootCmd.AddCommand(certCmd)

	certCmd.Flags().StringVarP(&certCrt, "file", "f", "", "Specify .crt file path")
	certCmd.Flags().StringVarP(&certHost, "domain", "d", "domain:port", "Specify domain and host port")
	certCmd.MarkFlagsMutuallyExclusive("file", "domain")
	certCmd.Flags().BoolVarP(&certIP, "ip", "p", false, "Print IP")
	certCmd.Flags().BoolVarP(&certExpiry, "expiry", "e", true, "Print expiry time")
	certCmd.Flags().BoolVarP(&certDNS, "dns", "n", false, "Print DNS names")
	certCmd.Flags().BoolVarP(&certIssuer, "issuer", "i", false, "Print issuer")
}

func certOutputByHost() {
	conn, err := tlsCheck.CheckByHost(certHost)
	if err != nil {
		log.Println(err)
		return
	}
	cert := conn.ConnectionState().PeerCertificates[0]
	if certExpiry {
		fmt.Println("Expiry time:", cert.NotAfter.Local().Format(time.RFC3339))
	}
	if certIP {
		fmt.Println("Server IP:", conn.RemoteAddr().String())
	}
	if certDNS {
		fmt.Println("DNS:", cert.DNSNames)
	}
	if certIssuer {
		fmt.Println("Issuer:", cert.Issuer.String())
	}
}

func certOutputByFile() {
	cert, err := tlsCheck.CheckByFile(certCrt)
	if err != nil {
		log.Println(err)
		return
	}
	if cert == nil {
		log.Println(cert)
		return
	}
	if certExpiry {
		fmt.Println("Expiry time:", cert[0].NotAfter.Local().Format(time.RFC3339))
	}
	if certDNS {
		fmt.Println("DNS:", cert[0].DNSNames)
	}
	if certIssuer {
		fmt.Println("Issuer:", cert[0].Issuer.String())
	}
}
