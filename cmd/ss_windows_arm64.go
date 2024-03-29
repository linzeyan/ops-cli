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

import "github.com/spf13/cobra"

func initSs() *cobra.Command {
	var ssCmd = &cobra.Command{
		GroupID: getGroupID(CommandSs),
		Use:     CommandSs,
		Short:   "Displays sockets informations",
		Run: func(_ *cobra.Command, _ []string) {
			printer.Printf(NotImplemented)
		},
	}
	return ssCmd
}
