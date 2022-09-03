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
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:   CommandStat + " file...",
	Short: "Display file informations",
	RunE: func(_ *cobra.Command, args []string) error {
		var err error
		if len(args) == 0 {
			f, err := os.Getwd()
			if err != nil {
				return err
			}
			stat, err := os.Lstat(f)
			if err != nil {
				return err
			}
			var s FileStat
			return s.String(f, stat)
		}
		for _, v := range args {
			f, err := filepath.Abs(v)
			if err != nil {
				return err
			}
			stat, err := os.Lstat(f)
			if err != nil {
				return err
			}
			var s FileStat
			err = s.String(v, stat)
			if err != nil {
				return err
			}
		}
		return err
	},
	DisableFlagsInUseLine: true,
	DisableFlagParsing:    true,
}

func init() {
	rootCmd.AddCommand(statCmd)
}

type FileStat struct {
	Dev           int32    `json:"Dev"`
	Mode          uint16   `json:"Mode"`
	Nlink         uint16   `json:"Nlink"`
	Ino           uint64   `json:"Ino"`
	UID           uint32   `json:"Uid"`
	GID           uint32   `json:"Gid"`
	Rdev          int32    `json:"Rdev"`
	PadCgo0       [4]byte  `json:"Pad_cgo_0"`
	Atimespec     Timespec `json:"Atimespec"`
	Mtimespec     Timespec `json:"Mtimespec"`
	Ctimespec     Timespec `json:"Ctimespec"`
	Birthtimespec Timespec `json:"Birthtimespec"`
	Size          int64    `json:"Size"`
	Blocks        int64    `json:"Blocks"`
	Blksize       int32    `json:"Blksize"`
	Flags         uint32   `json:"Flags"`
	Gen           uint32   `json:"Gen"`
	Lspare        int32    `json:"Lspare"`
	Qspare        [2]int64 `json:"Qspare"`
}

func (f *FileStat) String(path string, stat fs.FileInfo) error {
	var err error
	err = Encoder.JSONMarshaler(stat.Sys(), f)
	if err != nil {
		return err
	}
	var out string
	out = fmt.Sprintf(`  File: "%s"`, path)
	out += fmt.Sprintf("\n  Size: %s", common.ByteSize(f.Size).String())
	out += fmt.Sprintf("\t\tFileType: %d", stat.Mode())
	out += fmt.Sprintf("\n  Mode: (%04o/%s)", stat.Mode().Perm(), stat.Mode().String())
	uid, err := user.LookupId(fmt.Sprintf(`%d`, f.UID))
	if err != nil {
		return err
	}
	out += fmt.Sprintf("\tUid: (%d/%s)", f.UID, uid.Username)
	gid, err := user.LookupGroupId(fmt.Sprintf(`%d`, f.GID))
	if err != nil {
		return err
	}
	out += fmt.Sprintf("\tGid: (%d/%s)", f.GID, gid.Name)
	out += fmt.Sprintf("\nBlocks: %d", f.Blocks)
	out += fmt.Sprintf("\tBlock Size: %d", f.Blksize)
	// out += fmt.Sprintf("\nDevice: %d,%s", f.Rdev, f.Rdev)
	out += fmt.Sprintf("\tInode: %d", f.Ino)
	out += fmt.Sprintf("\tLinks: %d", f.Nlink)
	out += fmt.Sprintf("\nAccess: %s", time.Unix(f.Atimespec.SEC, f.Atimespec.Nsec).Local().Format(time.ANSIC))
	out += fmt.Sprintf("\nModify: %s", time.Unix(f.Mtimespec.SEC, f.Mtimespec.Nsec).Local().Format(time.ANSIC))
	out += fmt.Sprintf("\nChange: %s", time.Unix(f.Ctimespec.SEC, f.Ctimespec.Nsec).Local().Format(time.ANSIC))
	out += fmt.Sprintf("\n Birth: %s", time.Unix(f.Birthtimespec.SEC, f.Birthtimespec.Nsec).Local().Format(time.ANSIC))
	PrintString(out)
	return err
}

type Timespec struct {
	SEC  int64 `json:"Sec"`
	Nsec int64 `json:"Nsec"`
}
