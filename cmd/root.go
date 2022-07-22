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
	_ "embed"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "ops-cli",
	Short:   "OPS useful tools",
	Version: rootOutputVersion(),
	Run:     func(cmd *cobra.Command, _ []string) { cmd.Help() },
}

func Examples(s string) string {
	c := color.New(color.FgYellow)
	return c.Sprintf(`%s`, s)
}

/*
//go:generate /usr/local/bin/bash ../build.bash version
//go:embed version.txt
var version string
*/

var appVersion, appBuildTime, appCommit, appPlatform string
var rootOutputJson, rootOutputYaml bool

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&rootOutputJson, "json", "j", false, "Output JSON format")
	rootCmd.PersistentFlags().BoolVarP(&rootOutputYaml, "yaml", "y", false, "Output YAML format")
}

func rootOutputVersion() string {
	if appVersion != "" {
		return fmt.Sprintf("%s\nBuildTime: %s\nGitCommit: %s\nPlatform:  %s", appVersion, appBuildTime, appCommit, appPlatform)
	}
	return "v0.0.9"
}
