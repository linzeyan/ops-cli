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

	"github.com/a8m/tree"
	"github.com/a8m/tree/ostree"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func initTree() *cobra.Command {
	var flags struct {
		/* Listing options. */
		all        bool
		dirs       bool
		links      bool
		full       bool
		level      int
		ignoreCase bool
		noReport   bool
		pattern    string
		ignore     string
		output     string
		/* File options. */
		quote   bool
		protect bool
		uid     bool
		gid     bool
		size    bool
		human   bool
		date    bool
		inodes  bool
		device  bool
		/* Sorting options. */
		version   bool
		modify    bool
		change    bool
		unsort    bool
		reverse   bool
		dirsFirst bool
		sort      string
		/* Graphics options. */
		indent bool
		color  bool
	}
	var treeCmd = &cobra.Command{
		Use:   CommandTree,
		Short: "Show the contents of the giving directory as a tree",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveFilterDirs
		},
		Run: func(_ *cobra.Command, args []string) {
			if flags.level < 1 {
				logger.Info("level < 1")
				printer.Error(common.ErrInvalidArg)
				return
			}
			if len(args) == 0 {
				args = append(args, ".")
			}

			var nd, nf int
			var outFile = os.Stdout
			var err error
			/* Output to file. */
			if flags.output != "" {
				outFile, err = os.Create(flags.output)
				if err != nil {
					logger.Info(err.Error())
					printer.Error(err)
					return
				}
				defer outFile.Close()
			}
			/* Check sort-type. */
			if flags.sort != "" {
				switch flags.sort {
				case "version", "mtime", "ctime", "name", "size":
				default:
					printer.Printf("sort type '%s' not valid, should be one of: name,version,size,mtime,ctime", flags.sort)
					return
				}
			}

			opts := &tree.Options{
				Fs:      new(ostree.FS),
				OutFile: outFile,

				All:        flags.all,
				DirsOnly:   flags.dirs,
				FullPath:   flags.full,
				DeepLevel:  flags.level,
				FollowLink: flags.links,
				Pattern:    flags.pattern,
				IPattern:   flags.ignore,
				IgnoreCase: flags.ignoreCase,

				ByteSize: flags.size,
				UnitSize: flags.human,
				FileMode: flags.protect,
				ShowUid:  flags.uid,
				ShowGid:  flags.gid,
				LastMod:  flags.date,
				Quotes:   flags.quote,
				Inodes:   flags.inodes,
				Device:   flags.device,

				NoSort:    flags.unsort,
				ReverSort: flags.reverse,
				DirSort:   flags.dirsFirst,
				VerSort:   flags.version || flags.sort == "version",
				ModSort:   flags.modify || flags.sort == "mtime",
				CTimeSort: flags.change || flags.sort == "ctime",
				NameSort:  flags.sort == "name",
				SizeSort:  flags.sort == "size",

				NoIndent: flags.indent,
				Colorize: flags.color,
			}

			for _, dir := range args {
				inf := tree.New(dir)
				d, f := inf.Visit(opts)
				nd, nf = nd+d, nf+f
				inf.Print(opts)
				if !flags.noReport {
					footer := "\n%d "
					switch {
					default:
						footer += "directories"
					case nd == 1:
						footer += "directory"
					}

					footer = fmt.Sprintf(footer, nd)
					if !flags.dirs {
						switch {
						default:
							footer += ", %d files\n"
						case nf == 1:
							footer += ", %d file\n"
						}
						footer = fmt.Sprintf(footer, nf)
					}
					fmt.Fprint(outFile, footer)
				}
			}
		},
	}
	treeCmd.Flags().BoolVarP(&flags.all, "all", "a", false, common.Usage("List all files"))
	treeCmd.Flags().BoolVarP(&flags.dirs, "dirs", "d", false, common.Usage("List only directories"))
	treeCmd.Flags().BoolVarP(&flags.links, "links", "l", false, common.Usage("Follow symbolic links"))
	treeCmd.Flags().BoolVarP(&flags.full, "full", "f", false, common.Usage("Print full path for each file"))
	treeCmd.Flags().IntVarP(&flags.level, "level", "L", 3, common.Usage("Specify directory depth"))
	treeCmd.Flags().BoolVar(&flags.ignoreCase, "ignore-case", false, common.Usage("Ignore case when pattern matching"))
	treeCmd.Flags().BoolVar(&flags.noReport, "noreport", false, common.Usage("Disable file/directory count"))
	treeCmd.Flags().StringVarP(&flags.pattern, "pattern", "P", "", common.Usage("List only match the pattern given"))
	treeCmd.Flags().StringVarP(&flags.ignore, "ignore", "I", "", common.Usage("Do not list that match the given pattern"))
	treeCmd.Flags().StringVarP(&flags.output, "outputfile", "o", "", common.Usage("Output to file instead of stdout"))

	treeCmd.Flags().BoolVarP(&flags.quote, "quote", "Q", false, common.Usage("Quote filenames with double quotes"))
	treeCmd.Flags().BoolVarP(&flags.protect, "protect", "p", false, common.Usage("Print the protections for each file"))
	treeCmd.Flags().BoolVarP(&flags.uid, "uid", "u", false, common.Usage("Displays file owner or UID number"))
	treeCmd.Flags().BoolVarP(&flags.gid, "gid", "g", false, common.Usage("Displays file group owner or GID number"))
	treeCmd.Flags().BoolVarP(&flags.size, "size", "s", false, common.Usage("Print the size in bytes of each file"))
	treeCmd.Flags().BoolVarP(&flags.human, "human", "h", false, common.Usage("Print the size in a more human readable way"))
	treeCmd.Flags().BoolVarP(&flags.date, "date", "D", false, common.Usage("Print the date of last modification or (-c) status change"))
	treeCmd.Flags().BoolVar(&flags.inodes, "inodes", false, common.Usage("Print inode number of each file"))
	treeCmd.Flags().BoolVar(&flags.device, "device", false, common.Usage("Print device ID number to which each file belongs"))

	treeCmd.Flags().BoolVarP(&flags.version, "version", "v", false, common.Usage("Sort files alphanumerically by version"))
	treeCmd.Flags().BoolVarP(&flags.modify, "modify", "t", false, common.Usage("Sort files by last modification time"))
	treeCmd.Flags().BoolVarP(&flags.change, "change", "c", false, common.Usage("Sort files by last status change time"))
	treeCmd.Flags().BoolVarP(&flags.unsort, "unsort", "U", false, common.Usage("Leave files unsorted"))
	treeCmd.Flags().BoolVarP(&flags.reverse, "reverse", "r", false, common.Usage("Reverse the order of the sort"))
	treeCmd.Flags().BoolVar(&flags.dirsFirst, "dirsfirst", false, common.Usage("List directories before files (-U disables)"))
	treeCmd.Flags().StringVar(&flags.sort, "sort", "", common.Usage("Select sort: name,version,size,mtime,ctime"))

	treeCmd.Flags().BoolVarP(&flags.indent, "indent", "i", false, common.Usage("Don't print indentation lines"))
	treeCmd.Flags().BoolVarP(&flags.color, "color", "C", false, common.Usage("Turn colorization on always"))

	return treeCmd
}
