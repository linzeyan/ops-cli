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
	"os"
	"path/filepath"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func init() {
	var treeFlag TreeFlag
	var treeCmd = &cobra.Command{
		Use:   "tree",
		Short: "Show the contents of the giving directory as a tree",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveFilterDirs
		},
		Run: treeFlag.Run,
	}
	rootCmd.AddCommand(treeCmd)
	treeCmd.Flags().BoolVarP(&treeFlag.all, "all", "a", false, "List all files")
	treeCmd.Flags().BoolVarP(&treeFlag.dirs, "dirs", "d", false, "List only directories")
	treeCmd.Flags().BoolVarP(&treeFlag.full, "full", "f", false, "Print full path for each file")
	treeCmd.Flags().IntVarP(&treeFlag.limit, "limit", "l", 30, "Specify directories depth")
}

const (
	lastDirPrefix = "    "
	layerPrefix   = "│   "
	filePrefix    = "├── "
	endPrefix     = "└── "
)

type TreeFormat struct {
	Type     string        `json:"type"`
	Path     string        `json:"path"`
	Name     string        `json:"name"`
	Perm     string        `json:"perm"`
	Mode     string        `json:"mode"`
	Contents *[]TreeFormat `json:"contents"`

	layers int
}

type TreeFlag struct {
	all   bool
	dirs  bool
	full  bool
	limit int

	dirN, fileN int
	stat        FileStat
}

func (t *TreeFlag) Run(cmd *cobra.Command, args []string) {
	if t.limit < 1 {
		PrintString(common.ErrInvalidArg)
		return
	}

	if len(args) == 0 {
		args = append(args, ".")
	}

	for _, v := range args {
		dirName, err := filepath.Abs(v)
		if err != nil {
			PrintString(err)
			return
		}
		f, err := os.Lstat(dirName)
		if err != nil {
			PrintString(err)
			return
		}

		output := TreeFormat{
			Type:     t.stat.FileType(f),
			Path:     dirName,
			Name:     v,
			Contents: new([]TreeFormat),
			layers:   1,
		}

		err = t.iterate(&output)
		if err != nil {
			PrintString(err)
			return
		}
		switch {
		default:
			t.Print("", output)
		case rootOutputJSON:
			PrintJSON(output)
		case rootOutputYAML:
			PrintYAML(output)
		}
		t.summary()
	}
}

func (t *TreeFlag) iterate(trees *TreeFormat) error {
	files, err := os.ReadDir(trees.Path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !t.all {
			if strings.HasPrefix(f.Name(), ".") {
				continue
			}
		}
		if t.dirs {
			if !f.IsDir() {
				continue
			}
		}
		fi, err := f.Info()
		if err != nil {
			return err
		}
		temp := &TreeFormat{
			Type:     t.stat.FileType(fi),
			Path:     filepath.Join(trees.Path, f.Name()),
			Name:     f.Name(),
			Perm:     fmt.Sprintf("%#o", fi.Mode().Perm()),
			Mode:     fi.Mode().String(),
			Contents: &[]TreeFormat{},
			layers:   trees.layers + 1,
		}

		if trees.layers > t.limit {
			continue
		}

		if f.IsDir() {
			t.dirN++
			*trees.Contents = append(*trees.Contents, *temp)
			err = t.iterate(temp)
			if err != nil {
				return err
			}
		} else {
			t.fileN++
			*trees.Contents = append(*trees.Contents, *temp)
		}
	}
	return err
}

func (t *TreeFlag) Print(prefix string, output TreeFormat) {
	if t.full {
		PrintString(output.Path)
	} else {
		PrintString(output.Name)
	}

	for i, v := range *output.Contents {
		if i == len(*output.Contents)-1 {
			fmt.Print(prefix + endPrefix)
			t.Print(prefix+lastDirPrefix, v)
		} else {
			fmt.Print(prefix + filePrefix)
			t.Print(prefix+layerPrefix, v)
		}
	}
}

func (t *TreeFlag) summary() {
	out := "\n%d "
	switch {
	default:
		out += "directories"
	case t.dirN == 1:
		out += "directory"
	}
	switch {
	default:
		out += ", %d files\n"
	case t.fileN == 1:
		out += ", %d file\n"
	}
	PrintString(fmt.Sprintf(out, t.dirN, t.fileN))
	t.dirN, t.fileN = 0, 0
}
