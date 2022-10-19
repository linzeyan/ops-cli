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

func initTree() *cobra.Command {
	var treeCmd = &cobra.Command{
		Use:   CommandTree,
		Short: "Show the contents of the giving directory as a tree",
		Run: func(_ *cobra.Command, args []string) {
			PrintString(NotImplemented)
		},
	}
	return treeCmd
}
