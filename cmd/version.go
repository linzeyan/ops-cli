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
	"reflect"
	"runtime"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

var (
	appVersion   = "dev"
	appBuildTime = "unknown"
	appCommit    = "unknown"
)

func initVersion() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   CommandVersion,
		Short: "Print version information",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(_ *cobra.Command, _ []string) {
			var v = Version{
				Version: appVersion,
				Commit:  appCommit,
				Date:    appBuildTime,
				Runtime: fmt.Sprintf("%s %s", runtime.Version(), PlatformS),
			}
			printer.Printf(rootOutputFormat, v)
		},
	}
	return versionCmd
}

type Version struct {
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	Commit  string `json:"commit,omitempty" yaml:"commit,omitempty"`
	Date    string `json:"date,omitempty" yaml:"date,omitempty"`
	Runtime string `json:"runtime,omitempty" yaml:"runtime,omitempty"`
}

func (r Version) String() string {
	var ver strings.Builder
	f := reflect.ValueOf(&r).Elem()
	t := f.Type()
	ver.WriteString(fmt.Sprintf("%-10s\t%v\n", "App", common.RepoName))
	for i := 0; i < f.NumField(); i++ {
		_, _ = ver.WriteString(fmt.Sprintf("%-10s\t%v\n", t.Field(i).Name, f.Field(i).Interface()))
		// f.Field(i).Type()
	}
	ver.WriteString("Copyright © 2022 ZeYanLin <zeyanlin@outlook.com>\n")
	ver.WriteString(fmt.Sprintf("Source available at https://github.com/%s/%s", common.RepoOwner, common.RepoName))
	return ver.String()
}
