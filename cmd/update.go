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
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func initUpdate() *cobra.Command {
	var updateCmd = &cobra.Command{
		Use: CommandUpdate,
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Short: fmt.Sprintf("Update %s to the latest release", common.RepoName),
		Run: func(_ *cobra.Command, _ []string) {
			updater := NewUpdater(common.RepoOwner, common.RepoName, appVersion)
			if !updater.Upgrade {
				logger.Warn("", common.DefaultField(updater))
				printer.Printf("up-to-date\n")
				return
			}

			printer.Printf("Update...\n")
			printer.Printf("%s\n", common.Usage("==> Downloading file from GitHub "+updater.Repository.DownloadLink))
			err := updater.Download()
			if err != nil {
				logger.Error(err.Error())
				return
			}

			printer.Printf("Upgrading %s %s -> %s", common.RepoName, appVersion, updater.Repository.ReleaseTag)
			printer.Printf("%s\n", common.Usage("==> Cleanup..."))
			err = updater.Rename()
			if err != nil {
				logger.Warn(err.Error())
			}
			printer.Printf("Update completed in %s", time.Since(common.TimeNow))
		},

		DisableFlagsInUseLine: true,
	}
	return updateCmd
}

type version struct {
	Ver   string
	Major int
	Minor int
	Patch int
}

type repository struct {
	GithubUsername string
	Repository     string
	DownloadLink   string
	DownloadPath   string
	ExtractPath    string
	ReleaseTag     string
	ExtractFunc    func() error
}

/* fetchLatestVersion fetchs the latest release version tag, and return struct. */
func (r *repository) fetchLatestVersion(username, repo string) *repository {
	const githubURL = "https://github.com/%s/%s/releases/"
	urlBase := fmt.Sprintf(githubURL, username, repo)
	getTagURL := fmt.Sprintf("%s%s", urlBase, "latest")
	tagURL, err := common.HTTPRequestRedirectURL(getTagURL)
	if err != nil {
		logger.Debug(err.Error())
		tagURL = ""
	}
	tag := path.Base(tagURL)
	var extension string
	if common.IsWindows() {
		extension = "zip"
		r.ExtractFunc = r.UnZip
	} else {
		extension = "tar.gz"
		r.ExtractFunc = r.UnGzip
	}
	downloadLink := fmt.Sprintf("%s%s/%s/%s_%s_%s.%s",
		urlBase, "download", tag, repo, tag, PlatformU, extension)

	r.GithubUsername = username
	r.Repository = repo
	r.ReleaseTag = tag
	r.DownloadLink = downloadLink
	r.DownloadPath = path.Base(downloadLink)

	return r
}

/* Sanitize archive file pathing from G305: File traversal when extracting zip/tar archive. */
func (r *repository) sanitizeExtractPath(filePath string, destination string) (string, error) {
	destpath := filepath.Join(destination, filePath)
	if !strings.HasPrefix(destpath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return destination, common.ErrIllegalPath
	}
	return destpath, nil
}

/* G110: Potential DoS vulnerability via decompression bomb. */
func (*repository) copy(dst io.Writer, src io.Reader) error {
	for {
		_, err := io.CopyN(dst, src, 1024)
		if err != nil {
			logger.Debug(err.Error(), common.NewField("src", src), common.NewField("dst", dst))
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
}

/* Extract gzip and untar. */
func (r *repository) UnGzip() error {
	f, err := os.Open(r.DownloadPath)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(r.DownloadPath))
		return err
	}
	defer f.Close()
	/* Decompress */
	ungzip, err := gzip.NewReader(f)
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	defer ungzip.Close()

	/* Untar */
	reader := tar.NewReader(ungzip)

	for {
		header, err := reader.Next()
		logger.Debug(err.Error())
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if header == nil {
			continue
		}
		downloadPath, err := filepath.Abs(r.DownloadPath)
		if err != nil {
			logger.Debug(err.Error(), common.DefaultField(r.DownloadPath))
			return err
		}
		paths, err := r.sanitizeExtractPath(header.Name, filepath.Dir(downloadPath))
		if err != nil {
			logger.Debug(err.Error())
			return err
		}
		_, err = os.Stat(filepath.Dir(paths))
		if err != nil {
			logger.Debug(err.Error())
			if mkErr := os.MkdirAll(filepath.Dir(paths), os.ModePerm); mkErr != nil {
				logger.Debug(mkErr.Error())
				return mkErr
			}
		}
		if header.FileInfo().IsDir() {
			continue
		}
		if strings.Contains(filepath.Base(paths), r.Repository) {
			r.ExtractPath = paths
		}

		file, err := os.OpenFile(paths, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, header.FileInfo().Mode())
		if err != nil {
			logger.Debug(err.Error())
			return err
		}
		defer file.Close()
		err = r.copy(file, reader)
		if err != nil {
			logger.Debug(err.Error())
			return err
		}
	}
}

/* Unzip. */
func (r *repository) UnZip() error {
	unzip, err := zip.OpenReader(r.DownloadPath)
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	defer unzip.Close()
	for _, f := range unzip.File {
		downloadPath, err := filepath.Abs(r.DownloadPath)
		if err != nil {
			logger.Debug(err.Error(), common.DefaultField(r.DownloadPath))
			return err
		}
		paths, err := r.sanitizeExtractPath(f.Name, filepath.Dir(downloadPath))
		if err != nil {
			logger.Debug(err.Error())
			return err
		}
		dir := filepath.Dir(paths)
		_, err = os.Stat(dir)
		if err != nil {
			logger.Debug(err.Error())
			if mkErr := os.MkdirAll(dir, os.ModePerm); mkErr != nil {
				logger.Debug(mkErr.Error())
				return mkErr
			}
		}
		if f.FileInfo().IsDir() {
			continue
		}
		if strings.Contains(filepath.Base(paths), r.Repository) {
			r.ExtractPath = paths
		}

		dstFile, err := os.OpenFile(paths, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, f.Mode())
		if err != nil {
			logger.Debug(err.Error())
			return err
		}
		defer dstFile.Close()
		zipFile, err := f.Open()
		if err != nil {
			logger.Debug(err.Error())
			return err
		}
		defer zipFile.Close()
		err = r.copy(dstFile, zipFile)
		if err != nil {
			logger.Debug(err.Error())
			return err
		}
	}
	return err
}

func NewRepository(username, repo string) *repository {
	var r repository
	return r.fetchLatestVersion(username, repo)
}

type Updater struct {
	Upgrade        bool
	ExecutablePath string
	Repository     *repository
}

/* split splits string version into int, and return struct. */
func (*Updater) split(tag string) *version {
	if strings.Contains(tag, "dev") {
		return &version{
			Ver:   tag,
			Major: 0,
			Minor: 0,
			Patch: 0,
		}
	}
	replace := strings.Replace(tag, "v", "", 1)
	s := strings.Split(replace, ".")
	if len(s) != 3 {
		return &version{
			Ver:   s[0],
			Major: 0,
			Minor: 0,
			Patch: 0,
		}
	}
	major, err := strconv.Atoi(s[0])
	if err != nil {
		logger.Fatal(err.Error())
	}
	minor, err := strconv.Atoi(s[1])
	if err != nil {
		logger.Fatal(err.Error())
	}
	patch, err := strconv.Atoi(s[2])
	if err != nil {
		logger.Fatal(err.Error())
	}
	return &version{
		Ver:   tag,
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

/* compare compares latest version to current version, returns 1 if newer, 0 if the same, -1 if older. */
func (*Updater) compare(current, latest *version) int {
	if current.Major < latest.Major {
		return -1
	}
	if current.Major > latest.Major {
		return 1
	}
	if current.Minor < latest.Minor {
		return -1
	}
	if current.Minor > latest.Minor {
		return 1
	}
	if current.Patch < latest.Patch {
		return -1
	}
	if current.Patch > latest.Patch {
		return 1
	}
	return 0
}

/* Get current version and return struct. */
func (u *Updater) parseVersion(releaseTag string) *version {
	return u.split(releaseTag)
}

/* Get current version to compare with latest version, and if need to update fetch the executable file path. */
func (u *Updater) init(ver string) *Updater {
	current := u.parseVersion(ver)
	latest := u.parseVersion(u.Repository.ReleaseTag)
	compare := u.compare(current, latest)
	if compare < 0 && current.Ver != "dev" {
		u.Upgrade = true
		/* Get executable path. */
		execute, err := os.Executable()
		if err != nil {
			logger.Fatal(err.Error())
		}
		/* Get the real path. */
		realPath, err := filepath.EvalSymlinks(execute)
		if err != nil {
			logger.Fatal(err.Error())
		}
		u.ExecutablePath = realPath
	}
	return u
}

/* Fetch the release file. */
func (u *Updater) Download() error {
	var err error
	resp, err := common.HTTPRequestContent(u.Repository.DownloadLink)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(u.Repository.DownloadLink))
		return err
	}
	return os.WriteFile(u.Repository.DownloadPath, resp, FileModeRAll)
}

/* Decompress, replace original file, and remove compress files ...etc. */
func (u *Updater) Rename() error {
	/* Decompress */
	err := u.Repository.ExtractFunc()
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	/* Replace file. */
	err = os.Rename(u.Repository.ExtractPath, u.ExecutablePath)
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	/* Remove useless files. */
	err = os.Remove(u.Repository.DownloadPath)
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	extractDir := filepath.Dir(u.Repository.ExtractPath)
	return os.RemoveAll(extractDir)
}

func NewUpdater(username, repo, tag string) *Updater {
	var u Updater
	u.Repository = NewRepository(username, repo)
	return u.init(tag)
}
