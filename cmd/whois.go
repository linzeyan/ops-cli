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
	"io/ioutil"
	"log"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var whoisCmd = &cobra.Command{
	Use:   "whois",
	Short: "List domain name information",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
			return
		}
		var resp *whoisResponse
		var err error
		resp, err = wf.Request(args[0])
		if err != nil {
			log.Println(err)
			return
		}
		if resp == nil {
			log.Println("response is empty")
			return
		}
		OutputDefaultString(resp)
	},
	Example: Examples(`# Search domain
ops-cli whois apple.com`),
}

var wf whoisFlag

func init() {
	rootCmd.AddCommand(whoisCmd)

	whoisCmd.Flags().BoolVarP(&wf.ns, "ns", "n", false, "Only print Name Servers")
	whoisCmd.Flags().BoolVarP(&wf.expiry, "expiry", "e", false, "Only print expiry time")
	whoisCmd.Flags().BoolVarP(&wf.registrar, "registrar", "r", false, "Only print Registrar")
	whoisCmd.Flags().BoolVarP(&wf.days, "days", "d", false, "Only print the remaining days")
}

type whoisResponse struct {
	Registrar   string   `json:"registrar" yaml:"registrar"`
	CreatedDate string   `json:"createdDate" yaml:"createdDate"`
	ExpiresDate string   `json:"expiresDate" yaml:"expiresDate"`
	UpdatedDate string   `json:"updatedDate" yaml:"updatedDate"`
	RemainDays  int      `json:"remainDays" yaml:"remainDays"`
	NameServers []string `json:"nameServers" yaml:"nameServers"`
}

func (r whoisResponse) String() {
	if wf.expiry {
		fmt.Println(r.ExpiresDate)
		return
	}
	if wf.ns {
		fmt.Println(r.NameServers)
		return
	}
	if wf.registrar {
		fmt.Println(r.Registrar)
		return
	}
	if wf.days {
		fmt.Println(r.RemainDays)
		return
	}
	f := reflect.ValueOf(&r).Elem()
	t := f.Type()
	for i := 0; i < f.NumField(); i++ {
		fmt.Printf("%s\t%v\n", t.Field(i).Name, f.Field(i).Interface())
	}
}

func (r whoisResponse) JSON() { PrintJSON(r) }

func (r whoisResponse) YAML() { PrintYAML(r) }

type whoisFlag struct {
	/* Bind flags */
	ns, expiry, registrar, days bool
}

func (w whoisFlag) Request(domain string) (*whoisResponse, error) {
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
			r.UpdatedDate = w.ParseTime(v[1])
		}
		if strings.Contains(split[i], "Creation Date") {
			v := strings.Split(split[i], ";")
			r.CreatedDate = w.ParseTime(v[1])
		}
		if strings.Contains(split[i], "Registry Expiry Date") {
			v := strings.Split(split[i], ";")
			r.ExpiresDate = w.ParseTime(v[1])
			r.RemainDays = w.CalculateDays(v[1])
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

func (w whoisFlag) ParseTime(t string) string {
	/* 1997-09-15T04:00:00Z */
	s, err := time.Parse("2006-01-02T03:04:05Z", t)
	if err != nil {
		log.Println(err)
		return ""
	}
	return s.Local().Format(time.RFC3339)
}

func (w whoisFlag) CalculateDays(t string) int {
	s, err := time.Parse("2006-01-02T03:04:05Z", t)
	if err != nil {
		log.Println(err)
		return 0
	}
	return int(s.Local().Sub(rootNow.Local()).Hours() / 24)
}
