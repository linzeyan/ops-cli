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
	"io/fs"
	"runtime"
)

const (
	CommandArping     = "arping"
	CommandAudio      = "audio"
	CommandBase32     = "base32"
	CommandBase64     = "base64"
	CommandBase32Hex  = CommandBase32 + CommandHex
	CommandBase32Std  = CommandBase32 + "std"
	CommandBase64Std  = CommandBase64 + "std"
	CommandBase64URL  = CommandBase64 + "url"
	CommandBootstrap  = "bootstrap-token"
	CommandCalculate  = "calculate"
	CommandCert       = "cert"
	CommandConvert    = "convert"
	CommandCPU        = "cpu"
	CommandCsv        = "csv"
	CommandCsv2JSON   = CommandCsv + "2" + CommandJSON
	CommandCsv2Toml   = CommandCsv + "2" + CommandToml
	CommandCsv2XML    = CommandCsv + "2" + CommandXML
	CommandCsv2Yaml   = CommandCsv + "2" + CommandYaml
	CommandDate       = "date"
	CommandDf         = "df"
	CommandDig        = "dig"
	CommandDiscord    = "discord"
	CommandDisk       = "disk"
	CommandDoc        = "doc"
	CommandDos2Unix   = "dos2unix"
	CommandEncode     = "encode"
	CommandEncrypt    = "encrypt"
	CommandFile       = "file"
	CommandFree       = "free"
	CommandGenerate   = "generate"
	CommandGeoip      = "geoip"
	CommandHash       = "hash"
	CommandHex        = "hex"
	CommandHost       = "host"
	CommandICP        = "icp"
	CommandID         = "id"
	CommandIP         = "ip"
	CommandJSON       = "json"
	CommandJSON2Csv   = CommandJSON + "2" + CommandCsv
	CommandJSON2Toml  = CommandJSON + "2" + CommandToml
	CommandJSON2XML   = CommandJSON + "2" + CommandXML
	CommandJSON2Yaml  = CommandJSON + "2" + CommandYaml
	CommandLINE       = "line"
	CommandLoad       = "load"
	CommandLowercase  = "lowercase"
	CommandMan        = "man"
	CommandMarkdown   = "markdown"
	CommandMemory     = "memory"
	CommandMTR        = "mtr"
	CommandNetmask    = "netmask"
	CommandNetwork    = "network"
	CommandNumber     = "number"
	CommandOTP        = "otp"
	CommandPhoto      = "photo"
	CommandPing       = "ping"
	CommandPs         = "ps"
	CommandQrcode     = "qrcode"
	CommandRandom     = "random"
	CommandRead       = "read"
	CommandReadlink   = "readlink"
	CommandRedis      = "redis"
	CommandReST       = "rest"
	CommandSign       = "sign"
	CommandSlack      = "slack"
	CommandSs         = "ss"
	CommandSSH        = "ssh-keygen"
	CommandSSL        = "ssl"
	CommandStat       = "stat"
	CommandString     = "string"
	CommandSymbol     = "symbol"
	CommandSystem     = "system"
	CommandTCPing     = "tcping"
	CommandTelegram   = "telegram"
	CommandText       = "text"
	CommandToml       = "toml"
	CommandToml2Csv   = CommandToml + "2" + CommandCsv
	CommandToml2JSON  = CommandToml + "2" + CommandJSON
	CommandToml2XML   = CommandToml + "2" + CommandXML
	CommandToml2Yaml  = CommandToml + "2" + CommandYaml
	CommandTraceroute = "traceroute"
	CommandTree       = "tree"
	CommandUpdate     = "update"
	CommandUppercase  = "uppercase"
	CommandURL        = "url"
	CommandVersion    = "version"
	CommandVideo      = "video"
	CommandVoice      = "voice"
	CommandWhois      = "whois"
	CommandWiFi       = "wifi"
	CommandWsping     = "wsping"
	CommandXML        = "xml"
	CommandXML2Csv    = CommandXML + "2" + CommandCsv
	CommandXML2JSON   = CommandXML + "2" + CommandJSON
	CommandXML2Toml   = CommandXML + "2" + CommandToml
	CommandXML2Yaml   = CommandXML + "2" + CommandYaml
	CommandYaml       = "yaml"
	CommandYaml2Csv   = CommandYaml + "2" + CommandCsv
	CommandYaml2JSON  = CommandYaml + "2" + CommandJSON
	CommandYaml2Toml  = CommandYaml + "2" + CommandToml
	CommandYaml2XML   = CommandYaml + "2" + CommandXML
)

const (
	EncryptModeCFB = "CFB"
	EncryptModeCTR = "CTR"
	EncryptModeGCM = "GCM"
	EncryptModeOFB = "OFB"
)

const (
	FileModeROwner fs.FileMode = 0600
	FileModeRAll   fs.FileMode = 0644
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
	TypeBinary  = "binary"
	TypeOctal   = "octal"
	TypeDecimal = "decimal"
	TypeHex     = "hex"
	TypeCisco   = "cisco"
)

const (
	IndentTwoSpaces = "  "

	NotImplemented = "not implemented on " + PlatformS
	PlatformS      = runtime.GOOS + "/" + runtime.GOARCH
	PlatformU      = runtime.GOOS + "_" + runtime.GOARCH

	TCP  = "tcp"
	TCP6 = "tcp6"
	UDP  = "udp"
	UDP6 = "udp6"
	All  = "all"
	IPv4 = "ipv4"
	IPv6 = "ipv6"

	inetLocalhost   = "localhost"
	inetLocalhostIP = "127.0.0.1"

	keyFileExtension  = ".key"
	tempFileExtension = ".temp"

	mtrStatHeader = "Loss%   Snt   Last   Avg  Best  Wrst StDev"
)
