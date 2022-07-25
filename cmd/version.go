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
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(_ *cobra.Command, _ []string) {
		if appVersion == "" {
			appVersion = "v0.0.9"
		}
		var v = version{
			Version: appVersion,
			Commit:  appCommit,
			Date:    appBuildTime,
			Runtime: fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
		}

		if versionComplete {
			v.String()
			return
		}
		outputDefaultJson(v)
	},
}

var appVersion, appBuildTime, appCommit string
var versionComplete bool

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVarP(&versionComplete, "complete", "c", false, "Print version information completely")
}

type version struct {
	Version string `json:"Version,omitempty" yaml:"Version,omitempty"`
	Commit  string `json:"Commit,omitempty" yaml:"Commit,omitempty"`
	Date    string `json:"Date,omitempty" yaml:"Date,omitempty"`
	Runtime string `json:"Runtime,omitempty" yaml:"Runtime,omitempty"`
}

func (r version) Json() {
	out, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

func (r version) Yaml() {
	out, err := yaml.Marshal(r)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(out))
}

func (r version) String() {
	var ver strings.Builder
	f := reflect.ValueOf(&r).Elem()
	t := f.Type()
	ver.WriteString(fmt.Sprintf("%-10s\t%v\n", "App", "ops-cli"))
	for i := 0; i < f.NumField(); i++ {
		_, err := ver.WriteString(fmt.Sprintf("%-10s\t%v\n", t.Field(i).Name, f.Field(i).Interface()))
		//f.Field(i).Type()
		if err != nil {
			log.Println(err)
			return
		}
	}
	ver.WriteString("Copyright © 2022 ZeYanLin <zeyanlin@outlook.com>\n")
	ver.WriteString("Source available at https://github.com/linzeyan/ops-cli")
	fmt.Println(ver.String())
}
