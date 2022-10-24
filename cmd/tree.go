//go:build darwin || linux

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
	"syscall"
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func initTree() *cobra.Command {
	var flags TreeOptions
	var treeCmd = &cobra.Command{
		Use:   CommandTree,
		Short: "Show the contents of the giving directory as a tree",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveFilterDirs
		},
		Run: func(_ *cobra.Command, args []string) {
			if flags.Limit < 1 {
				printer.Error(common.ErrInvalidArg)
				return
			}
			if len(args) == 0 {
				args = append(args, ".")
			}
			var t Tree
			t.Walk(args, &flags)
		},
	}
	treeCmd.Flags().BoolVarP(&flags.All, "all", "a", false, common.Usage("List all files"))
	treeCmd.Flags().BoolVarP(&flags.Change, "change", "c", false, common.Usage("Print the date of last modification"))
	treeCmd.Flags().BoolVarP(&flags.Dirs, "dirs", "d", false, common.Usage("List only directories"))
	treeCmd.Flags().BoolVarP(&flags.Full, "full", "f", false, common.Usage("Print full path for each file"))
	treeCmd.Flags().BoolVarP(&flags.Perm, "perm", "p", false, common.Usage("Print file permission"))
	treeCmd.Flags().BoolVarP(&flags.Mode, "mode", "m", false, common.Usage("Print file mode"))
	treeCmd.Flags().BoolVarP(&flags.Size, "size", "s", false, common.Usage("Print the size for each file"))
	treeCmd.Flags().IntVarP(&flags.Limit, "limit", "l", 30, common.Usage("Specify directories depth"))
	treeCmd.Flags().BoolVarP(&flags.GID, "gid", "g", false, common.Usage("Print group owner for each file"))
	treeCmd.Flags().BoolVarP(&flags.UID, "uid", "u", false, common.Usage("Print owner for each file"))
	treeCmd.Flags().BoolVarP(&flags.Links, "links", "", false, common.Usage("Print links for each file"))
	treeCmd.Flags().BoolVarP(&flags.Inodes, "inodes", "", false, common.Usage("Print inode number for each file"))
	treeCmd.Flags().BoolVarP(&flags.Device, "device", "", false, common.Usage("Print device ID number for each file"))
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

type Tree struct {
	Type     string  `json:"type"`
	Path     string  `json:"path"`
	Name     string  `json:"name"`
	Perm     string  `json:"perm"`
	Mode     string  `json:"mode"`
	Links    string  `json:"links"`
	UID      string  `json:"uid"`
	GID      string  `json:"gid"`
	Size     string  `json:"size"`
	ModTime  string  `json:"modTime"`
	Inode    string  `json:"inode"`
	Devide   string  `json:"device"`
	Contents *[]Tree `json:"contents"`

	layers      int
	dirN, fileN int
	stat        FileStat
}

func (t *Tree) Walk(args []string, opt *TreeOptions) {
	for _, v := range args {
		dirName, err := filepath.Abs(v)
		if err != nil {
			printer.Error(err)
			return
		}
		f, err := os.Lstat(dirName)
		if err != nil {
			printer.Error(err)

			return
		}
		uid, gid, links, inode, device, err := t.getInfo(f)
		if err != nil {
			printer.Error(err)

			return
		}
		*t = Tree{
			Type:     t.stat.FileType(f),
			Path:     dirName,
			Name:     v,
			Perm:     fmt.Sprintf("%#o", f.Mode().Perm()),
			Mode:     f.Mode().String(),
			Size:     common.ByteSize(f.Size()).String(),
			ModTime:  f.ModTime().Format(time.ANSIC),
			UID:      uid,
			GID:      gid,
			Links:    links,
			Inode:    inode,
			Devide:   device,
			Contents: new([]Tree),
			layers:   1,
		}

		err = t.iterate(t, opt)
		if err != nil {
			printer.Error(err)

			return
		}
		t.print("", *t, opt)
		t.summary()
	}
}

/* getInfo returns uid, gid, inodes, device and error. */
func (t *Tree) getInfo(fileinfo fs.FileInfo) (string, string, string, string, string, error) {
	s, ok := fileinfo.Sys().(*syscall.Stat_t)
	if !ok {
		return "", "", "0", "", "", common.ErrResponse
	}
	uid, err := user.LookupId(fmt.Sprintf(`%d`, s.Uid))
	if err != nil {
		return "", "", "0", "", "", err
	}
	gid, err := user.LookupGroupId(fmt.Sprintf(`%d`, s.Gid))
	if err != nil {
		return "", "", "0", "", "", err
	}
	return uid.Username, gid.Name,
		fmt.Sprintf("%d", s.Nlink),
		fmt.Sprintf("%d", s.Ino),
		fmt.Sprintf("%d", s.Dev), err
}

func (t *Tree) iterate(trees *Tree, opt *TreeOptions) error {
	files, err := os.ReadDir(trees.Path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !opt.All {
			if strings.HasPrefix(f.Name(), ".") {
				continue
			}
		}
		if opt.Dirs {
			if !f.IsDir() {
				continue
			}
		}
		fi, err := f.Info()
		if err != nil {
			return err
		}
		uid, gid, links, inode, device, err := t.getInfo(fi)
		if err != nil {
			return err
		}
		temp := &Tree{
			Type:     t.stat.FileType(fi),
			Path:     filepath.Join(trees.Path, f.Name()),
			Name:     f.Name(),
			Perm:     fmt.Sprintf("%#o", fi.Mode().Perm()),
			Mode:     fi.Mode().String(),
			Size:     common.ByteSize(fi.Size()).String(),
			ModTime:  fi.ModTime().Format(time.ANSIC),
			Links:    links,
			UID:      uid,
			GID:      gid,
			Inode:    inode,
			Devide:   device,
			Contents: &[]Tree{},
			layers:   trees.layers + 1,
		}

		if trees.layers > opt.Limit {
			continue
		}

		if f.IsDir() {
			t.dirN++
			*trees.Contents = append(*trees.Contents, *temp)
			err = t.iterate(temp, opt)
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

func (t *Tree) print(prefix string, output Tree, opt *TreeOptions) {
	if rootOutputFormat != "" {
		printer.Printf(rootOutputFormat, output)
		return
	}
	var p []string
	if opt.Mode {
		p = append(p, output.Mode)
	}
	if opt.Perm {
		p = append(p, output.Perm)
	}
	if opt.Links {
		p = append(p, output.Links)
	}
	if opt.UID {
		p = append(p, output.UID)
	}
	if opt.GID {
		p = append(p, output.GID)
	}
	if opt.Size {
		p = append(p, output.Size)
	}
	if opt.Change {
		p = append(p, output.ModTime)
	}
	if opt.Inodes {
		p = append(p, output.Inode)
	}
	if opt.Device {
		p = append(p, output.Devide)
	}
	if len(p) != 0 {
		fmt.Printf("%v ", p)
	}

	if opt.Full {
		printer.Printf("%s\n", output.Path)
	} else {
		printer.Printf("%s\n", output.Name)
	}

	for i, v := range *output.Contents {
		if i == len(*output.Contents)-1 {
			printer.Printf("%s%s", prefix, treePerfixEnd)
			t.print(prefix+treePerfixEmpty, v, opt)
		} else {
			printer.Printf("%s%s", prefix, treePerfixFile)
			t.print(prefix+treePerfixLayer, v, opt)
		}
	}
}

func (t *Tree) summary() {
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
	printer.Printf(out, t.dirN, t.fileN)
	t.dirN, t.fileN = 0, 0
}
