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
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

var whoisCmd = &cobra.Command{
	Use:   "whois [domain]",
	Args:  cobra.ExactArgs(1),
	Short: "List domain name information",
	Run: func(_ *cobra.Command, args []string) {
		var resp *WhoisResponse
		var err error
		resp, err = whoisCmdGlobalVar.Request(args[0])
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if resp == nil {
			return
		}
		OutputInterfaceString(resp)
	},
	Example: common.Examples(`# Search domain
ops-cli whois apple.com`),
}

var whoisCmdGlobalVar WhoisFlag

func init() {
	rootCmd.AddCommand(whoisCmd)

	whoisCmd.Flags().BoolVarP(&whoisCmdGlobalVar.ns, "ns", "n", false, "Only print Name Servers")
	whoisCmd.Flags().BoolVarP(&whoisCmdGlobalVar.expiry, "expiry", "e", false, "Only print expiry time")
	whoisCmd.Flags().BoolVarP(&whoisCmdGlobalVar.registrar, "registrar", "r", false, "Only print Registrar")
	whoisCmd.Flags().BoolVarP(&whoisCmdGlobalVar.days, "days", "d", false, "Only print the remaining days")
}

type WhoisResponse struct {
	Registrar   string   `json:"registrar" yaml:"registrar"`
	CreatedDate string   `json:"createdDate" yaml:"createdDate"`
	ExpiresDate string   `json:"expiresDate" yaml:"expiresDate"`
	UpdatedDate string   `json:"updatedDate" yaml:"updatedDate"`
	RemainDays  int      `json:"remainDays" yaml:"remainDays"`
	NameServers []string `json:"nameServers" yaml:"nameServers"`
}

func (r WhoisResponse) String() {
	if whoisCmdGlobalVar.expiry {
		PrintString(r.ExpiresDate)
		return
	}
	if whoisCmdGlobalVar.ns {
		PrintJSON(r.NameServers)
		return
	}
	if whoisCmdGlobalVar.registrar {
		PrintString(r.Registrar)
		return
	}
	if whoisCmdGlobalVar.days {
		PrintString(r.RemainDays)
		return
	}
	f := reflect.ValueOf(&r).Elem()
	t := f.Type()
	for i := 0; i < f.NumField(); i++ {
		fmt.Printf("%s\t%v\n", t.Field(i).Name, f.Field(i).Interface())
	}
}

type WhoisFlag struct {
	/* Bind flags */
	ns, expiry, registrar, days bool
}

func (w WhoisFlag) Request(domain string) (*WhoisResponse, error) {
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
	result, err := io.ReadAll(conn)
	if err != nil {
		return nil, err
	}
	replace := strings.ReplaceAll(string(result), ": ", ";")
	replace1 := strings.ReplaceAll(replace, "\r\n", ",")
	split := strings.Split(replace1, ",")
	var ns []string
	var r WhoisResponse
	var calErr error
	/* Filter field. */
	for i := range split {
		v := strings.Split(split[i], ";")
		if strings.Contains(split[i], "Updated Date") {
			r.UpdatedDate, err = w.ParseTime(v[1])
		}
		if strings.Contains(split[i], "Creation Date") {
			r.CreatedDate, err = w.ParseTime(v[1])
		}
		if strings.Contains(split[i], "Registry Expiry Date") {
			r.ExpiresDate, err = w.ParseTime(v[1])
			r.RemainDays, calErr = w.CalculateDays(v[1])
		}
		if strings.Contains(split[i], "Registrar") {
			if strings.TrimSpace(v[0]) == "Registrar" {
				r.Registrar = v[1]
			}
		}
		if strings.Contains(split[i], "Name Server") {
			ns = append(ns, v[1])
		}
		if err != nil {
			log.Println(err)
		}
		if calErr != nil {
			log.Println(calErr)
			err = calErr
		}
	}
	r.NameServers = ns
	return &r, err
}

/* Convert time to RFC3339 format. */
func (w WhoisFlag) ParseTime(t string) (string, error) {
	/* 1997-09-15T04:00:00Z */
	s, err := time.Parse("2006-01-02T15:04:05Z", t)
	if err != nil {
		return "", err
	}
	return s.Local().Format(time.RFC3339), err
}

/* Convert time to days. */
func (w WhoisFlag) CalculateDays(t string) (int, error) {
	s, err := time.Parse("2006-01-02T15:04:05Z", t)
	if err != nil {
		return 0, err
	}
	return int(s.Local().Sub(common.TimeNow.Local()).Hours() / 24), err
}
