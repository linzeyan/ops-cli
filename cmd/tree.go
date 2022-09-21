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
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var treeFlag TreeFlag
	var treeCmd = &cobra.Command{
		Use:   "tree",
		Short: "Show the contents of the giving directory as a tree",
		Run:   treeFlag.Run,
	}
	rootCmd.AddCommand(treeCmd)
}

type TreeFlag struct{}

func (t *TreeFlag) Run(cmd *cobra.Command, args []string) {
	const (
		layerPrefix = "│   "
		filePrefix  = "├── "
		endPrefix   = "└── "
	)
	var arg string
	if len(args) == 0 {
		arg = "."
	} else {
		arg = args[0]
	}

	var dirN, fileN int
	err := filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
		dir, file := filepath.Split(path)
		layer := strings.Count(dir, string(filepath.Separator))
		var prefix string
		for i := 1; i < layer; i++ {
			prefix += layerPrefix
		}
		if layer != 0 {
			prefix += filePrefix
			if d.IsDir() {
				dirN++
			} else {
				fileN++
			}
		}

		fmt.Println(prefix + file)
		return err
	})
	if err != nil {
		return
	}
	out := "\n%d "
	switch {
	default:
		out += "directories"
	case dirN == 1:
		out += "directory"
	}
	switch {
	default:
		out += ", %d files\n"
	case fileN == 1:
		out += ", %d file\n"
	}
	fmt.Printf(out, dirN, fileN)
}
