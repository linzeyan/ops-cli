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
	var flags TreeOptions
	var treeCmd = &cobra.Command{
		Use:   CommandTree,
		Short: "Show the contents of the giving directory as a tree",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveFilterDirs
		},
		Run: func(_ *cobra.Command, args []string) {
			printer.Printf(NotImplemented)
		},
	}
	treeCmd.Flags().BoolVarP(&flags.All, "all", "a", false, "List all files")
	treeCmd.Flags().BoolVarP(&flags.Change, "change", "c", false, "Print the date of last modification")
	treeCmd.Flags().BoolVarP(&flags.Dirs, "dirs", "d", false, "List only directories")
	treeCmd.Flags().BoolVarP(&flags.Full, "full", "f", false, "Print full path for each file")
	treeCmd.Flags().BoolVarP(&flags.Perm, "perm", "p", false, "Print file permission")
	treeCmd.Flags().BoolVarP(&flags.Mode, "mode", "m", false, "Print file mode")
	treeCmd.Flags().BoolVarP(&flags.Size, "size", "s", false, "Print the size for each file")
	treeCmd.Flags().IntVarP(&flags.Limit, "limit", "l", 30, "Specify directories depth")
	treeCmd.Flags().BoolVarP(&flags.GID, "gid", "g", false, "Print group owner for each file")
	treeCmd.Flags().BoolVarP(&flags.UID, "uid", "u", false, "Print owner for each file")
	treeCmd.Flags().BoolVarP(&flags.Links, "links", "", false, "Print links for each file")
	treeCmd.Flags().BoolVarP(&flags.Inodes, "inodes", "", false, "Print inode number for each file")
	treeCmd.Flags().BoolVarP(&flags.Device, "device", "", false, "Print device ID number for each file")
	return treeCmd
}

type TreeOptions struct {
	All    bool
	Change bool
	Dirs   bool
	Full   bool
	Limit  int
	Perm   bool
	Mode   bool
	Size   bool
	UID    bool
	GID    bool
	Inodes bool
	Device bool
	Links  bool
}
