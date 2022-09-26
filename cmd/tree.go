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
	Name     string        `json:"name"`
	Contents *[]TreeFormat `json:"contents"`

	layers int
	prefix []string
}

type TreeFlag struct {
	all   bool
	limit int

	dirN, fileN int
	dirName     string
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
		t.dirName = v
		output := TreeFormat{
			Type:     t.stat.FileType(f),
			Name:     v,
			Contents: new([]TreeFormat),
		}

		err = t.iterate(v, output.Contents)
		if err != nil {
			PrintString(err)
			return
		}
		t.Print(output)
		// PrintJSON(output)
		t.summary()
	}
}

func (t *TreeFlag) iterate(arg string, contents *[]TreeFormat) error {
	files, err := os.ReadDir(arg)
	if err != nil {
		return err
	}
	n := len(files)
	for i := 0; i < n; i++ {
		f := files[i]
		if !t.all {
			if strings.HasPrefix(f.Name(), ".") {
				continue
			}
		}
		fullpath := filepath.Join(arg, f.Name())
		// need modify
		layer := strings.Count(fullpath, string(filepath.Separator))
		if t.dirName == "." {
			layer++
		}
		if layer > t.limit {
			continue
		}

		var prefix []string
		for j := 1; j < layer; j++ {
			prefix = append(prefix, layerPrefix)
		}
		if i == n-1 {
			prefix = append(prefix, endPrefix)
		} else {
			prefix = append(prefix, filePrefix)
		}

		fInfo, err := f.Info()
		if err != nil {
			return err
		}
		temp := &TreeFormat{
			Type:     t.stat.FileType(fInfo),
			Name:     f.Name(),
			Contents: &[]TreeFormat{},
			prefix:   prefix,
			layers:   layer,
		}

		if f.IsDir() {
			t.dirN++
			*contents = append(*contents, *temp)
			err = t.iterate(fullpath, temp.Contents)
			if err != nil {
				return err
			}
		} else {
			t.fileN++
			*contents = append(*contents, *temp)
		}
	}
	return err
}

func (t *TreeFlag) Print(output TreeFormat) {
	if output.Type == "Directory" {
		for i := 0; i < output.layers-1; i++ {
			if output.prefix[i] == layerPrefix {
				output.prefix[i] = lastDirPrefix
			}
		}
	}
	PrintString(common.SliceStringToString(output.prefix) + output.Name)

	for _, v := range *output.Contents {
		t.Print(v)
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
