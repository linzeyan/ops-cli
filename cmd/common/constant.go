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

package common

import "os"

const (
	_            = iota
	KiB ByteSize = 1 << (10 * iota)
	MiB
	GiB
	TiB
	PiB
	EiB
	ZiB
	YiB
)

const (
	FileModeROwner os.FileMode = 0600
	FileModeRAll   os.FileMode = 0644
)

const (
	HashMd5        = "md5"
	HashSha1       = "sha1"
	HashSha224     = "sha224"
	HashSha256     = "sha256"
	HashSha384     = "sha384"
	HashSha512     = "sha512"
	HashSha512_224 = "sha512_224"
	HashSha512_256 = "sha512_256"
)

const (
	IndentTwoSpaces = "  "

	RepoOwner = "linzeyan"
	RepoName  = "ops-cli"

	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"
)
