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
)

type FileStat struct{}

func (f *FileStat) String(path string) error {
	var err error
	stat, err := os.Lstat(path)
	if err != nil {
		return err
	}
	s, ok := stat.Sys().(*syscall.Stat_t)
	if !ok {
		return common.ErrResponse
	}
	var out string
	out = fmt.Sprintf(`  File: "%s"`, path)
	out += fmt.Sprintf("\n  Size: %s", common.ByteSize(s.Size).String())
	out += fmt.Sprintf("\t\tBlocks: %d", s.Blocks)
	out += fmt.Sprintf("\tIO Block: %d", s.Blksize)
	out += fmt.Sprintf("\tFileType: %s", f.FileType(stat))
	out += fmt.Sprintf("\n  Mode: (%#o/%s)", stat.Mode().Perm(), stat.Mode())
	uid, err := user.LookupId(fmt.Sprintf(`%d`, s.Uid))
	if err != nil {
		return err
	}
	out += fmt.Sprintf("\tUid: (%5d/%8s)", s.Uid, uid.Username)
	gid, err := user.LookupGroupId(fmt.Sprintf(`%d`, s.Gid))
	if err != nil {
		return err
	}
	out += fmt.Sprintf("\tGid: (%5d/%8s)", s.Gid, gid.Name)
	out += fmt.Sprintf("\nDevice: %d", s.Dev)
	out += fmt.Sprintf("\tInode: %d", s.Ino)
	out += fmt.Sprintf("\tLinks: %d", s.Nlink)
	out += fmt.Sprintf("\nAccess: %s", time.Unix(s.Atim.Unix()).Local().Format(time.ANSIC))
	out += fmt.Sprintf("\nModify: %s", time.Unix(s.Mtim.Unix()).Local().Format(time.ANSIC))
	out += fmt.Sprintf("\nChange: %s", time.Unix(s.Ctim.Unix()).Local().Format(time.ANSIC))
	printer.Printf(out)
	return err
}

func (*FileStat) FileType(stat fs.FileInfo) string {
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
