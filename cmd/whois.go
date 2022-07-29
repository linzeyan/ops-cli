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
	"net"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var whoisCmd = &cobra.Command{
	Use:   "whois",
	Short: "List domain name information",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
			return
		}
		var whoisDomain = args[0]
		var resp *whoisResponse
		var err error
		var req whoisVerisign
		resp, err = req.Request(whoisDomain)
		if err != nil {
			log.Println(err)
			return
		}
		if resp == nil {
			log.Println("response is empty")
			return
		}
		outputDefaultString(resp)
	},
	Example: Examples(`# Search domain
ops-cli whois apple.com`),
}

var whoisNameServer, whoisExpiry, whoisRegistrar bool

func init() {
	rootCmd.AddCommand(whoisCmd)

	whoisCmd.Flags().BoolVarP(&whoisNameServer, "ns", "n", false, "Only print Name Servers")
	whoisCmd.Flags().BoolVarP(&whoisExpiry, "expiry", "e", false, "Only print expiry time")
	whoisCmd.Flags().BoolVarP(&whoisRegistrar, "registrar", "r", false, "Only print Registrar")
}

type whoisResponse struct {
	Registrar   string   `json:"registrar" yaml:"registrar"`
	CreatedDate string   `json:"createdDate" yaml:"createdDate"`
	ExpiresDate string   `json:"expiresDate" yaml:"expiresDate"`
	UpdatedDate string   `json:"updatedDate" yaml:"updatedDate"`
	NameServers []string `json:"nameServers" yaml:"nameServers"`
}

func (r whoisResponse) String() {
	if whoisExpiry {
		fmt.Println(r.ExpiresDate)
		return
	}
	if whoisNameServer {
		fmt.Println(r.NameServers)
		return
	}
	if whoisRegistrar {
		fmt.Println(r.Registrar)
		return
	}
	f := reflect.ValueOf(&r).Elem()
	t := f.Type()
	for i := 0; i < f.NumField(); i++ {
		fmt.Printf("%s\t%v\n", t.Field(i).Name, f.Field(i).Interface())
	}
}

func (r whoisResponse) JSON() {
	out, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

func (r whoisResponse) YAML() {
	out, err := yaml.Marshal(r)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

type whoisVerisign struct{}

func (w whoisVerisign) Request(domain string) (*whoisResponse, error) {
	conn, err := net.Dial("tcp", "whois.verisign-grs.com:43")
	if err != nil {
		return nil, err
	}
	if conn != nil {
		defer conn.Close()
	}
	_, err = conn.Write([]byte(domain + "\n"))
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, err
	}

	replace := strings.ReplaceAll(string(result), ": ", ";")
	replace1 := strings.ReplaceAll(replace, "\r\n", ",")
	split := strings.Split(replace1, ",")
	var ns []string
	var r whoisResponse
	for i := range split {
		if strings.Contains(split[i], "Updated Date") {
			v := strings.Split(split[i], ";")
			r.UpdatedDate = v[1]
		}
		if strings.Contains(split[i], "Creation Date") {
			v := strings.Split(split[i], ";")
			r.CreatedDate = v[1]
		}
		if strings.Contains(split[i], "Registry Expiry Date") {
			v := strings.Split(split[i], ";")
			r.ExpiresDate = v[1]
		}
		if strings.Contains(split[i], "Registrar") {
			v := strings.Split(split[i], ";")
			if strings.TrimSpace(v[0]) == "Registrar" {
				r.Registrar = v[1]
			}
		}
		if strings.Contains(split[i], "Name Server") {
			v := strings.Split(split[i], ";")
			ns = append(ns, v[1])
		}
	}
	r.NameServers = ns
	return &r, nil
}
