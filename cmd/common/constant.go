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
	_           = iota
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

const (
	Discord  ConfigBlock = "discord"
	Encrypt  ConfigBlock = "encrypt"
	ICP      ConfigBlock = "west"
	LINE     ConfigBlock = "line"
	Slack    ConfigBlock = "slack"
	Telegram ConfigBlock = "telegram"
)

const (
	FileModeROwner os.FileMode = 0600
	FileModeRAll   os.FileMode = 0644
)

const (
	Base32    = "base32"
	Base64    = "base64"
	Hex       = "hex"
	Base32Hex = Base32 + Hex
	Base32Std = Base32 + "std"
	Base64Std = Base64 + "std"
	Base64URL = Base64 + "url"
)

const (
	CommandRoot     = "ops-cli"
	CommandCert     = "cert"
	CommandConvert  = "convert"
	CommandDig      = "dig"
	CommandDiscord  = "Discord"
	CommandDoc      = "doc"
	CommandDos2Unix = "dos2unix"
	CommandEncode   = "encode"
	CommandEncrypt  = "encrypt"
	CommandGeoip    = "geoip"
	CommandHash     = "hash"
	CommandIcp      = "icp"
	CommandLINE     = "LINE"
	CommandOtp      = "otp"
	CommandQrcode   = "qrcode"
	CommandRandom   = "random"
	CommandSlack    = "Slack"
	CommandSSH      = "ssh-keygen"
	CommandSystem   = "system"
	CommandTelegram = "Telegram"
	CommandUpdate   = "update"
	CommandURL      = "url"
	CommandVersion  = "version"
	CommandWhois    = "whois"

	SubCommandAudio     = "audio"
	SubCommandBootstrap = "bootstrap-token"
	SubCommandCalculate = "calculate"
	SubCommandExpand    = "expand"
	SubCommandFile      = "file"
	SubCommandGenerate  = "generate"
	SubCommandGet       = "get"
	SubCommandID        = "id"
	SubCommandLowercase = "lowercase"
	SubCommandMan       = "man"
	SubCommandMarkdown  = "markdown"
	SubCommandNumber    = "number"
	SubCommandPhoto     = "photo"
	SubCommandRead      = "read"
	SubCommandReST      = "rest"
	SubCommandString    = "string"
	SubCommandSymbol    = "symbol"
	SubCommandText      = "text"
	SubCommandUppercase = "uppercase"
	SubCommandVideo     = "video"
	SubCommandVoice     = "voice"
	SubCommandWiFi      = "wifi"
	SubCommandYaml      = "yaml"
	SubCommandYaml2JSON = "yaml2json"
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

	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"
)
