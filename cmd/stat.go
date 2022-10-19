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
	"syscall"
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func initStat() *cobra.Command {
	var statCmd = &cobra.Command{
		Use:   CommandStat + " path...",
		Short: "Display file informations",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var err error
			for _, v := range args {
				var s FileStat
				err = s.String(v)
				if err != nil {
					return err
				}
			}
			return err
		},
		DisableFlagsInUseLine: true,
	}
	return statCmd
}

type FileStat struct{}

func (f *FileStat) String(path string) error {
	var err error
	stat, err := os.Lstat(path)
	if err != nil {
		return err
	}
	err = Encoder.JSONMarshaler(stat.Sys(), f)
	if err != nil {
		return err
	}
	var out string
	out = fmt.Sprintf(`  File: "%s"`, path)
	out += fmt.Sprintf("\n  Size: %s", common.ByteSize(stat.Sys().(*syscall.Stat_t).Size).String())
	out += fmt.Sprintf("\t\tBlocks: %d", stat.Sys().(*syscall.Stat_t).Blocks)
	out += fmt.Sprintf("\tIO Block: %d", stat.Sys().(*syscall.Stat_t).Blksize)
	out += fmt.Sprintf("\tFileType: %s", f.FileType(stat))
	out += fmt.Sprintf("\n  Mode: (%#o/%s)", stat.Mode().Perm(), stat.Mode())
	uid, err := user.LookupId(fmt.Sprintf(`%d`, stat.Sys().(*syscall.Stat_t).Uid))
	if err != nil {
		return err
	}
	out += fmt.Sprintf("\tUid: (%5d/%8s)", stat.Sys().(*syscall.Stat_t).Uid, uid.Username)
	gid, err := user.LookupGroupId(fmt.Sprintf(`%d`, stat.Sys().(*syscall.Stat_t).Gid))
	if err != nil {
		return err
	}
	out += fmt.Sprintf("\tGid: (%5d/%8s)", stat.Sys().(*syscall.Stat_t).Gid, gid.Name)
	out += fmt.Sprintf("\nDevice: %d", stat.Sys().(*syscall.Stat_t).Dev)
	out += fmt.Sprintf("\tInode: %d", stat.Sys().(*syscall.Stat_t).Ino)
	out += fmt.Sprintf("\tLinks: %d", stat.Sys().(*syscall.Stat_t).Nlink)
	out += fmt.Sprintf("\nAccess: %s", time.Unix(stat.Sys().(*syscall.Stat_t).Atimespec.Unix()).Local().Format(time.ANSIC))
	out += fmt.Sprintf("\nModify: %s", time.Unix(stat.Sys().(*syscall.Stat_t).Mtimespec.Unix()).Local().Format(time.ANSIC))
	out += fmt.Sprintf("\nChange: %s", time.Unix(stat.Sys().(*syscall.Stat_t).Ctimespec.Unix()).Local().Format(time.ANSIC))
	out += fmt.Sprintf("\n Birth: %s", time.Unix(stat.Sys().(*syscall.Stat_t).Birthtimespec.Unix()).Local().Format(time.ANSIC))
	PrintString(out)
	return err
}

func (f *FileStat) FileType(stat fs.FileInfo) string {
	switch stat.Mode() & fs.ModeType {
	case fs.ModeDir: // d
		return "Directory"
	case fs.ModeAppend: // a
		return "append-only"
	case fs.ModeExclusive: // l
		return "exclusive use"
	case fs.ModeTemporary: // T
		return "temporary file"
	case fs.ModeSymlink: // L
		return "Symbolic Link"
	case fs.ModeDevice: // D
		return "Block Device"
	case fs.ModeNamedPipe: // p
		return "named pipe"
	case fs.ModeSocket: // S
		return "Socket"
	case fs.ModeSetuid: // u
		return "setuid"
	case fs.ModeSetgid: // g
		return "setgid"
	case fs.ModeCharDevice, fs.ModeCharDevice | fs.ModeDevice: // c, Dc
		return "Character Device"
	case fs.ModeSticky: // t
		return "sticky"
	case fs.ModeIrregular: // ?
		return "Non-regular file"
	case fs.ModePerm, fs.ModeType:
		return "unknown"
	default:
		return "Regular File"
	}
}
