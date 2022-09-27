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
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

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
	treeCmd.Flags().BoolVarP(&treeFlag.change, "change", "c", false, "Print the date of last modification")
	treeCmd.Flags().BoolVarP(&treeFlag.dirs, "dirs", "d", false, "List only directories")
	treeCmd.Flags().BoolVarP(&treeFlag.full, "full", "f", false, "Print full path for each file")
	treeCmd.Flags().BoolVarP(&treeFlag.perm, "perm", "p", false, "Print file permission")
	treeCmd.Flags().BoolVarP(&treeFlag.mode, "mode", "m", false, "Print file mode")
	treeCmd.Flags().BoolVarP(&treeFlag.size, "size", "s", false, "Print the size for each file")
	treeCmd.Flags().IntVarP(&treeFlag.limit, "limit", "l", 30, "Specify directories depth")
	treeCmd.Flags().BoolVarP(&treeFlag.gid, "gid", "g", false, "Print group owner for each file")
	treeCmd.Flags().BoolVarP(&treeFlag.uid, "uid", "u", false, "Print owner for each file")
	treeCmd.Flags().BoolVarP(&treeFlag.inodes, "inodes", "", false, "Print inode number for each file")
	treeCmd.Flags().BoolVarP(&treeFlag.device, "device", "", false, "Print device ID number for each file")
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
	Size     string        `json:"size"`
	ModTime  string        `json:"modTime"`
	UID      string        `json:"uid"`
	GID      string        `json:"gid"`
	Inode    string        `json:"inode"`
	Devide   string        `json:"device"`
	Contents *[]TreeFormat `json:"contents"`

	layers int
}

type TreeFlag struct {
	all    bool
	change bool
	dirs   bool
	full   bool
	limit  int
	perm   bool
	mode   bool
	size   bool
	uid    bool
	gid    bool
	inodes bool
	device bool

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
		uid, gid, inode, device, err := t.getInfo(f)
		if err != nil {
			PrintString(err)
			return
		}
		output := TreeFormat{
			Type:     t.stat.FileType(f),
			Path:     dirName,
			Name:     v,
			Perm:     fmt.Sprintf("%#o", f.Mode().Perm()),
			Mode:     f.Mode().String(),
			Size:     common.ByteSize(f.Size()).String(),
			ModTime:  f.ModTime().Format(time.ANSIC),
			UID:      uid,
			GID:      gid,
			Inode:    inode,
			Devide:   device,
			Contents: new([]TreeFormat),
			layers:   1,
		}

		err = t.iterate(&output)
		if err != nil {
			PrintString(err)
			return
		}
		t.Print("", output)
		t.summary()
	}
}

/* getInfo returns uid, gid, inodes, device and error. */
func (t *TreeFlag) getInfo(fileinfo fs.FileInfo) (string, string, string, string, error) {
	err := Encoder.JSONMarshaler(fileinfo.Sys(), &t.stat)
	if err != nil {
		return "", "", "", "", err
	}
	uid, err := user.LookupId(fmt.Sprintf(`%d`, t.stat.UID))
	if err != nil {
		return "", "", "", "", err
	}
	gid, err := user.LookupGroupId(fmt.Sprintf(`%d`, t.stat.GID))
	if err != nil {
		return "", "", "", "", err
	}
	return uid.Name, gid.Name, fmt.Sprintf("%d", t.stat.Ino), fmt.Sprintf("%d", t.stat.Dev), err
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
		uid, gid, inode, device, err := t.getInfo(fi)
		if err != nil {
			return err
		}
		temp := &TreeFormat{
			Type:     t.stat.FileType(fi),
			Path:     filepath.Join(trees.Path, f.Name()),
			Name:     f.Name(),
			Perm:     fmt.Sprintf("%#o", fi.Mode().Perm()),
			Mode:     fi.Mode().String(),
			Size:     common.ByteSize(fi.Size()).String(),
			ModTime:  fi.ModTime().Format(time.ANSIC),
			UID:      uid,
			GID:      gid,
			Inode:    inode,
			Devide:   device,
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
	switch {
	case rootOutputJSON:
		PrintJSON(output)
		return
	case rootOutputYAML:
		PrintYAML(output)
		return
	}

	var p []string
	if t.uid {
		p = append(p, output.UID)
	}
	if t.gid {
		p = append(p, output.GID)
	}
	if t.change {
		p = append(p, output.ModTime)
	}
	if t.mode {
		p = append(p, output.Mode)
	}
	if t.perm {
		p = append(p, output.Perm)
	}
	if t.size {
		p = append(p, output.Size)
	}
	if t.inodes {
		p = append(p, output.Inode)
	}
	if t.device {
		p = append(p, output.Devide)
	}
	if len(p) != 0 {
		fmt.Printf("%v ", p)
	}

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
