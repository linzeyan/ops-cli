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
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func initWhois() *cobra.Command {
	var flags struct {
		/* Bind flags */
		ns, expiry, registrar, days bool
	}
	var whoisCmd = &cobra.Command{
		Use:  CommandWhois + " domain",
		Args: cobra.ExactArgs(1),
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Short: "List domain name information",
		Run: func(_ *cobra.Command, args []string) {
			var w Whois
			if err := w.Request(args[0]); err != nil {
				logger.Info(err.Error(), common.DefaultField(args))
				printer.Error(err)
				return
			}
			if w.empty() {
				logger.Info(common.ErrResponse.Error())
				return
			}
			if flags.expiry {
				printer.Printf("%s\n", w.ExpiresDate)
			}
			if flags.ns {
				printer.Printf(rootOutputFormat, w.NameServers)
			}
			if flags.registrar {
				printer.Printf("%s\n", w.Registrar)
			}
			if flags.days {
				printer.Printf("%d\n", w.RemainDays)
			}
			if flags.expiry || flags.ns || flags.registrar || flags.days {
				return
			}
			printer.Printf(rootOutputFormat, w)
		},
		Example: common.Examples(`# Search domain
apple.com`, CommandWhois),
	}
	whoisCmd.Flags().BoolVarP(&flags.ns, "ns", "n", false, common.Usage("Only print Name Servers"))
	whoisCmd.Flags().BoolVarP(&flags.expiry, "expiry", "e", false, common.Usage("Only print expiry time"))
	whoisCmd.Flags().BoolVarP(&flags.registrar, "registrar", "r", false, common.Usage("Only print Registrar"))
	whoisCmd.Flags().BoolVarP(&flags.days, "days", "d", false, common.Usage("Only print the remaining days"))
	return whoisCmd
}

type Whois struct {
	Registrar   string   `json:"registrar" yaml:"registrar"`
	CreatedDate string   `json:"createdDate" yaml:"createdDate"`
	ExpiresDate string   `json:"expiresDate" yaml:"expiresDate"`
	UpdatedDate string   `json:"updatedDate" yaml:"updatedDate"`
	RemainDays  int      `json:"remainDays" yaml:"remainDays"`
	NameServers []string `json:"nameServers" yaml:"nameServers"`
}

func (w *Whois) Request(domain string) error {
	conn, err := net.Dial("tcp", net.JoinHostPort("whois.verisign-grs.com", "43"))
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	if conn != nil {
		defer conn.Close()
	}
	_, err = conn.Write([]byte(domain + "\n"))
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	result, err := io.ReadAll(conn)
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	replace := strings.ReplaceAll(string(result), ": ", ";")
	replace1 := strings.ReplaceAll(replace, "\r\n", ",")
	split := strings.Split(replace1, ",")
	var ns []string
	var calErr error
	/* Filter field. */
	for i := range split {
		v := strings.Split(split[i], ";")
		if strings.Contains(split[i], "Updated Date") {
			w.UpdatedDate, err = w.ParseTime(v[1])
		}
		if strings.Contains(split[i], "Creation Date") {
			w.CreatedDate, err = w.ParseTime(v[1])
		}
		if strings.Contains(split[i], "Registry Expiry Date") {
			w.ExpiresDate, err = w.ParseTime(v[1])
			w.RemainDays, calErr = w.CalculateDays(v[1])
		}
		if strings.Contains(split[i], "Registrar") {
			if strings.TrimSpace(v[0]) == "Registrar" {
				w.Registrar = v[1]
			}
		}
		if strings.Contains(split[i], "Name Server") {
			ns = append(ns, v[1])
		}
		if err != nil {
			logger.Debug(err.Error())
			printer.Error(err)
		}
		if calErr != nil {
			logger.Debug(calErr.Error())
			printer.Error(calErr)
			err = calErr
		}
	}
	w.NameServers = ns
	return err
}

/* Convert time to RFC3339 format. */
func (w *Whois) ParseTime(t string) (string, error) {
	/* 1997-09-15T04:00:00Z */
	s, err := time.Parse("2006-01-02T15:04:05Z", t)
	if err != nil {
		logger.Debug(err.Error())
		return "", err
	}
	return s.Local().Format(time.RFC3339), err
}

/* Convert time to days. */
func (w *Whois) CalculateDays(t string) (int, error) {
	s, err := time.Parse("2006-01-02T15:04:05Z", t)
	if err != nil {
		logger.Debug(err.Error())
		return 0, err
	}
	return int(s.Local().Sub(common.TimeNow.Local()).Hours() / 24), err
}

func (w *Whois) empty() bool {
	if w.CreatedDate != w.ExpiresDate || w.UpdatedDate != w.Registrar {
		return false
	}
	return true
}

func (w Whois) String() string {
	var name []string
	rt := reflect.TypeOf(w)
	for _, f := range reflect.VisibleFields(rt) {
		if f.IsExported() {
			name = append(name, f.Name)
		}
	}
	var s strings.Builder
	f := reflect.ValueOf(&w).Elem()
	for _, v := range name {
		temp := fmt.Sprintf("%s\t%v\n", v, f.FieldByName(v).Interface())
		s.WriteString(temp)
	}
	return s.String()
}
