# ops-cli

[![https://github.com/linzeyan/ops-cli/actions?query=workflow:test](https://github.com/linzeyan/ops-cli/workflows/test/badge.svg?branch=main)](https://github.com/linzeyan/ops-cli/actions?query=workflow:test)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/linzeyan/ops-cli/release?style=flat-square)](https://github.com/linzeyan/ops-cli/actions?query=workflow:release)
[![Release](https://img.shields.io/github/release/linzeyan/ops-cli.svg?style=flat-square)](https://github.com/linzeyan/ops-cli/releases/latest)
[![Software License](https://img.shields.io/github/license/linzeyan/ops-cli?style=flat-square)](./LICENSE)

Try to collect useful tools for ops.

## Table of contents

- [Installation](#installation)
  - [Go Install](#go-install)
  - [Homebrew](#homebrew)
- [Upgrade](#upgrade)
  - [Brew](#brew)
  - [Others](#others)
- [Usage](#usage)

## Installation

### Docker

```bash
docker pull zeyanlin/ops-cli
```

### Go Install

```bash
go install github.com/linzeyan/ops-cli@latest
```

### Homebrew

```bash
brew tap linzeyan/tools
brew install ops-cli
```

## Upgrade

### Brew

```bash
brew upgrade ops-cli
```

### Others

```bash
ops-cli update
```

### [Download Page](https://github.com/linzeyan/ops-cli/releases/latest)

## Usage

```bash
OPS useful tools

Usage:
  ops-cli
  ops-cli [command]

Available Commands:
  Discord     Send message to Discord
  LINE        Send message to LINE
  Slack       Send message to Slack
  Telegram    Send message to Telegram
  arping      Discover and probe hosts in a network using the ARP protocol
  cert        Check tls cert expiry time
  completion  Generate the autocompletion script for the specified shell
  convert     Convert data format, support csv, json, toml, xml, yaml
  date        Print date time
  df          Display free disk spaces
  dig         Resolve domain name
  doc         Generate documentation
  dos2unix    Convert file eol to unix style
  encode      Encode and decode string or file
  encrypt     Encrypt or decrypt
  free        Display free memory spaces
  geoip       Print IP geographic information
  hash        Hash string or file
  help        Help about any command
  icp         Check ICP status
  ip          View interfaces configuration
  mtr         Combined traceroute and ping
  netmask     Print IP/Mask pair, list address ranges
  otp         Calculate passcode or generate secret
  ping        Send ICMP ECHO_REQUEST packets to network hosts
  qrcode      Read or generate QR Code
  random      Generate random string
  readlink    Get symlink information
  redis       Opens a connection to a Redis server
  ss          Displays sockets informations
  ssh-keygen  Generate SSH keypair
  ssl         Genreate self-sign certificate
  stat        Display file informations
  system      Display system informations
  tcping      Connect to a port of a host
  traceroute  Print the route packets trace to network host
  tree        Show the contents of the giving directory as a tree
  update      Update ops-cli to the latest release
  url         Get url content or expand shorten url or download
  version     Print version information
  whois       List domain name information
  wsping      Connect to a websocket server

Flags:
      --config string   Specify config path
  -h, --help            help for ops-cli
      --output string   Output format, can be json/yaml
```

### `arping`

```bash
→ sudo ops-cli arping 192.168.181.1 -c
online
```

### `cert`

```bash
→ ops-cli cert www.google.com
{
  "expiryTime": "2022-11-07T16:25:22+08:00",
  "days": 63,
  "issuer": "CN=GTS CA 1C3,O=Google Trust Services LLC,C=US",
  "serverIp": "142.251.42.228:443",
  "dns": [
    "www.google.com"
  ]
}
```

### `convert`

```bash
→ ops-cli convert toml2yaml -i .config.reference.toml
→ cat .config.reference.yaml
discord:
  channel_id: channelID
  token: token
encrypt:
  key: ""
line:
  access_token: Channel Access Token
  id: ""
  secret: Channel Secret
slack:
  channel_id: CHANNEL
  token: token
telegram:
  chat_id: "12345678"
  token: token:token
west:
  account: account
  api_key: apikey
```

### `date`

```bash
→ ops-cli date -s
1663222044

→ ops-cli date -M
1663822818051885
```

```bash
→ ops-cli date --format '01-02-2006'
09-15-2022
```

### `df`

```bash
→ ops-cli df
Filesystem      Size      Used      Avail     Use%     Mounted on                   FsType  iUsed    iFree       iUse%
/dev/disk1s1s1  465.63GB  169.81GB  295.82GB  36.47%   /                            apfs    501138   3101858760  0.02%
devfs           191.50KB  191.50KB  0.00B     100.00%  /dev                         devfs   664      0           100.00%
/dev/disk1s5    465.63GB  169.81GB  295.82GB  36.47%   /System/Volumes/VM           apfs    3        3101858760  0.00%
/dev/disk1s3    465.63GB  169.81GB  295.82GB  36.47%   /System/Volumes/Preboot      apfs    4005     3101858760  0.00%
/dev/disk1s6    465.63GB  169.81GB  295.82GB  36.47%   /System/Volumes/Update       apfs    493      3101858760  0.00%
/dev/disk1s2    465.63GB  169.81GB  295.82GB  36.47%   /System/Volumes/Data         apfs    1988242  3101858760  0.06%
map auto_home   0.00B     0.00B     0.00B     0.00%    /System/Volumes/Data/home    autofs  0        0           0.00%
/dev/disk1s1    465.63GB  169.81GB  295.82GB  36.47%   /System/Volumes/Update/mnt1  apfs    502050   3101858760  0.02%
```

### `dig`

```bash
→ ops-cli dig @1.1.1.1 tw.yahoo.com CNAME
NAME            TTL     CLASS   TYPE    RECORD
tw.yahoo.com.   20      IN      CNAME   fp-ycpi.g03.yahoodns.net.
```

### `Discord`

```bash
ops-cli Discord text 'hello' --config ~/.config/.myconfig
```

### `dos2unix`

```bash
→ ops-cli dos2unix /tmp/abc.com/*
Converting file /tmp/abc.com/abc.com.crt to Unix format...
Converting file /tmp/abc.com/abc.com.key to Unix format...
```

### `encode`

```bash
→ ops-cli encode base64std 'https://github.com'
aHR0cHM6Ly9naXRodWIuY29t
→ ops-cli encode base64std -d 'aHR0cHM6Ly9naXRodWIuY29t'
https://github.com
```

### `encrypt`

```bash
→ ops-cli encrypt string 'https://github.com' --key '0123456789012345'
SmNEHlJ1QUw6yLyzcTQ1uibhg4SnTWuOkwo5c4A69JtVgw==
→ ops-cli encrypt string 'SmNEHlJ1QUw6yLyzcTQ1uibhg4SnTWuOkwo5c4A69JtVgw==' --key '0123456789012345' -d
https://github.com
```

### `free`

```bash
→ ops-cli free
         total     used       free  available    use%
Mem:   16.00GB  10.22GB   838.58MB     5.78GB  63.85%
Swap:   3.00GB   2.01GB  1012.00MB             67.06%
```

### `geoip`

```bash
→ ops-cli geoip 1.1.1.1
{
  "continent": "Oceania",
  "country": "Australia",
  "countryCode": "AU",
  "regionName": "Queensland",
  "city": "South Brisbane",
  "district": "",
  "timezone": "Australia/Brisbane",
  "currency": "AUD",
  "isp": "Cloudflare, Inc",
  "org": "APNIC and Cloudflare DNS Resolver project",
  "as": "AS13335 Cloudflare, Inc.",
  "asname": "CLOUDFLARENET",
  "mobile": false,
  "proxy": false,
  "hosting": true,
  "query": "1.1.1.1"
}
```

### `hash`

```bash
→ ops-cli hash -l main.go
MD5: 1b48671ec88f2b498820450f802097a6
SHA1: b2828f7ce1eb4872e5543617f9bf2b7ef28e7c61
SHA256: 74d680e9a561929551611bcecf9fc1704c75a8491e9aec00065adbfaefa36905
SHA512: 9aa41fc10c66de39e30dcf7e35be7c15878322df8e318922dc8623f140464c51b964953d950d8dc1ef4b0d4b614f01c065ef50d0c0e686cdd21ff4580fec3ce7
```

```bash
→ ops-cli hash sha512 'This is a string.'
0145c77435b886e43fcfa5b8a6e2c5a9f1c216f694a65e75354f9679174551b7a0151b72f5497d58845bc5033f39f3249ee087cdb602680edc3fdeda8a18ff9b
```

### `icp`

```bash
→ ops-cli icp apple.com --config ~/.config/.myconfig
domain: apple.com
icp: 京ICP备10214630号
icpstatus: 已备案
```

### `ip`

```bash
→ ops-cli ip en0
6: en0: <UP,BROADCAST,MULTICAST> mtu 1500
        ether 14:7d:da:aa:46:53
        inet 192.168.181.74/24
        inet6 fe80::18d6:da54:513a:20e5/64
        RX packets 62046457  bytes 45060208660 (41.97GB)
        RX errors 0  dropped 0
        TX packets 38729722  bytes 26083432718 (24.29GB)
        TX errors 0  dropped 15935
```

### `LINE`

```bash
ops-cli LINE text 'hello' --config ~/.config/.myconfig
```

### `mtr`

```bash
ops-cli mtr 1.1.1.1

                                                Packets               Pings
 Host                                          Loss%   Snt   Last   Avg  Best  Wrst StDev

 1. 192.168.181.1                               0.0%    16    1.9   2.5   1.2  10.0   2.4
 2. 61.220.168.254                              0.0%    16   12.3   8.4   3.3  13.1   2.7
 3. 168.95.83.214                               0.0%    16   24.0   6.6   2.5  32.9   8.4
 4. 220.128.27.94                               0.0%    16    3.6   3.8   2.8   5.3 0.656
 5. 220.128.25.181                              0.0%    16    4.7   5.6   3.0   8.6   1.5
 6. 220.128.4.77                                0.0%    16    3.3   3.3   2.7   4.7 0.513
 7. 210.242.214.45                              0.0%    16    5.8   6.7   4.0  19.7   3.6
 8. 1.1.1.1                                     0.0%    16    3.0   5.4   3.0  26.6   5.5
```

### `netmask`

```bash
→ ops-cli netmask -b 192.168.0.0/16
11000000 10101000 00000000 00000000 / 11111111 11111111 00000000 00000000
```

```bash
→ ops-cli netmask -r 172.16.0.0/12
172.16.0.0 -> 172.31.255.255 (1048576)
```

```bash
→ ops-cli netmask -c 172.16.9.1-192.168.3.1
172.16.9.1/32
172.16.9.2/31
172.16.9.4/30
172.16.9.8/29
172.16.9.16/28
172.16.9.32/27
172.16.9.64/26
172.16.9.128/25
172.16.10.0/23
172.16.12.0/22
172.16.16.0/20
172.16.32.0/19
172.16.64.0/18
172.16.128.0/17
172.17.0.0/16
172.18.0.0/15
172.20.0.0/14
172.24.0.0/13
172.32.0.0/11
172.64.0.0/10
172.128.0.0/9
173.0.0.0/8
174.0.0.0/7
176.0.0.0/4
192.0.0.0/9
192.128.0.0/11
192.160.0.0/13
192.168.0.0/23
192.168.2.0/24
192.168.3.0/31
```

### `otp`

```bash
→ ops-cli otp calculate 6BDR T7AT RRCZ V5IS FLOH AHQL YF4Z ORG7
631843
```

### `ping`

```bash
→ sudo ops-cli ping www.google.com -c 2
Password:
PING www.google.com (172.217.163.36): 24 data bytes
32 bytes from 172.217.163.36: icmp_seq=0 ttl=57 time=3.417797ms
32 bytes from 172.217.163.36: icmp_seq=1 ttl=57 time=4.2143ms

--- www.google.com ping statistics ---
2 packets transmitted, 2 packets received, 0.00% packet loss
round-trip min/avg/max/stddev = 3.417797ms/3.816048ms/4.2143ms/398.251µs
```

### `qrcode`

```bash
→ ops-cli qrcode read ~/Downloads/qrcode-two.png
HTTPS://MAGICLEN.ORG
```

### `random`

```bash
→ ops-cli random bootstrap-token
7fa086.d40039e9efc249ec
```

```bash
→ ops-cli random base64 -l 500
TKxrKqeIoY/FrGI4h2vNnalc0Ga0h3CoKtWVBluWF+Lu4CSF078oUPstVF0CTzBD
OQdtIGSNUEepInGU0vH4BWw9yo2iHJib9Ti9EAOGJs/caiy7QSQhJ2c5oamHyFu9
a0Li9ULAm+IsaY0od1xN0QDn3uwnGST/1x2dzDPzaHpU53P7IEWtb1Zioirk8VQZ
ZBHd3e/VNye93i0p0FpVxCi1Q+MhdBJzjx8F9f0arT4brHi0TruAgXd3TlYjjl7A
6F7iL9AfsUshtL2tKCF3EcmPv1UtCm8S7Hh6lYltzXrmCTCqBKI+doNp2/yPioHQ
IrpMxjepApa4IxMFapm5LSCBI8UrF4avpTi2cnD8wJfj/RAUlf/8VtbirUh8rtot
tdR715KHaNRXHpMRwmkDWHTpOtlCotBmKUEkQXJ5ZZyWkp7kkC6F6ADIrtcoa/7h
fR/ivHlj+PtIyB4em5y0Nq2nTuajHRjZbZTC7akTKx070UWH4uNajkP/Eq82E6Kn
QbqZnLicRpJj60TUiAFuA3ohjqFynoozRUWjmI8ZmBD6J4o+2viDUHTLRCdIuVeV
Garvm8Mi95i4Zt8ZMjOhcI6W9zzuMb9bJZqBVsl8GYANHo35X2sfR6burD94We7O
JkiKqH3a33NP9AA2xpVGrcQrBdo=
```

### `readlink`

```bash
→ ops-cli readlink /tmp
private/tmp
```

### `redis`

```bash
→ ops-cli redis 'set name Joe'
"OK"
→ ops-cli redis 'get name'
"Joe"
```

### `Slack`

```bash
ops-cli Slack text 'hello' --config ~/.config/.myconfig
```

### `ss`

```bash
ops-cli ss
Proto Local Address Foreign Address  State PID/Program name
  tcp    0.0.0.0:80       0.0.0.0:0 LISTEN         25/nginx
 tcp6         :::80            :::0 LISTEN         25/nginx
```

### `ssh-keygen`

```bash
→ ops-cli ssh-keygen --bits 4096 -f /tmp/rsa
/tmp/rsa generated
/tmp/rsa.pub generated
```

### `ssl`

```bash
→ ops-cli ssl generate
→ ls
ca.crt     ca.key     root.crt   root.key   server.crt server.key
```

### `stat`

```bash
→ ops-cli stat docker.sock
  File: "docker.sock"
  Size: 36.00B          Blocks: 0       IO Block: 4096  FileType: Symbolic Link
  Mode: (0755/Lrwxr-xr-x)       Uid: (    0/    root)   Gid: (    1/  daemon)
Device: 16777221        Inode: 37003970 Links: 1
Access: Wed Aug 10 13:25:51 2022
Modify: Wed Aug 10 13:25:51 2022
Change: Wed Aug 10 13:25:51 2022
 Birth: Wed Aug 10 13:25:51 2022
```

### `system`

```bash
→ ops-cli system cpu
{
  "VendorID": "GenuineIntel",
  "ModelName": "Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz",
  "Cores": 4,
  "CacheSize": 256,
  "GHz": 2
}
```

### `tcping`

```bash
→ ops-cli tcping 1.1.1.1 80
tcp response from 1.1.1.1 (1.1.1.1) port 80 [open] 5.311941ms
```

### `Telegram`

```bash
ops-cli Telegram text 'hello' --config ~/.config/.myconfig
```

### `traceroute`

```bash
→ sudo ops-cli traceroute 1.1.1.1
traceroute to 1.1.1.1 (1.1.1.1), 64 hops max, 24 byte packets
1     192.168.181.1    4.951536ms 2.309237ms 3.206471ms
2     61.220.168.254   11.807506ms 3.370591ms 15.020509ms
3     168.95.83.214    10.43289ms 9.700833ms 9.874369ms
4     220.128.27.94    9.995535ms 4.686122ms 9.784072ms
5     220.128.25.181   11.213021ms 7.17189ms 11.335182ms
6     220.128.4.77     9.642501ms 7.084146ms 5.372547ms
7     210.242.214.45   8.191182ms 4.458661ms 11.622122ms
8     1.1.1.1          11.124769ms 10.348664ms 3.23984ms
```

### `tree`

```bash
→ ops-cli tree cmd/common -gp
cmd/common
├── [-rw-r--r-- 20  ]  bytes.go
├── [-rw-r--r-- 20  ]  common.go
├── [-rw-r--r-- 20  ]  config.go
├── [-rw-r--r-- 20  ]  constant.go
├── [-rw-r--r-- 20  ]  http.go
├── [-rw-r--r-- 20  ]  printer.go
├── [-rw-r--r-- 20  ]  qrcode.go
├── [-rw-r--r-- 20  ]  validator.go
└── [-rw-r--r-- 20  ]  variable.go

0 directories, 9 files
```

### `update`

```bash
→ ops-cli update
Update...
==> Downloading file from GitHub
Upgrading ops-cli v0.10.0 -> v0.10.1
==> Cleanup...
Update completed
```

### `url`

```bash
→ ops-cli url -e https://bit.ly/3gk7w5x
https://zh.wikipedia.org/zh-tw/%E5%8F%8D%E5%90%91%E4%BB%A3%E7%90%86#:~:text=%E5%8F%8D%E5%90%91%E4%BB%A3%E7%90%86%E5%9C%A8%E9%9B%BB%E8%85%A6,%E4%BC%BA%E6%9C%8D%E5%99%A8%E5%8F%A2%E9%9B%86%E7%9A%84%E5%AD%98%E5%9C%A8%E3%80%82
```

### `version`

```bash
→ ops-cli version
App        ops-cli
Version    v0.8.4
Commit     e4b96dfeb732a81440969877a6bb5fdef17d5d09
Date       2022-09-21T01:40:44Z
Runtime    go1.18.5 darwin/amd64
Copyright © 2022 ZeYanLin <zeyanlin@outlook.com>
Source available at https://github.com/linzeyan/ops-cli
→ ops-cli version --output json
{
  "version": "v0.8.4",
  "commit": "e4b96dfeb732a81440969877a6bb5fdef17d5d09",
  "date": "2022-09-21T01:40:44Z",
  "runtime": "go1.18.5 darwin/amd64"
}
```

### `whois`

```bash
→ ops-cli whois apple.com --output json
{
  "registrar": "CSC Corporate Domains",
  "createdDate": "1987-02-19T13:00:00+08:00",
  "expiresDate": "2023-02-20T13:00:00+08:00",
  "updatedDate": "2022-02-16T14:15:06+08:00",
  "remainDays": 168,
  "nameServers": [
    "A.NS.APPLE.COM",
    "B.NS.APPLE.COM",
    "C.NS.APPLE.COM",
    "D.NS.APPLE.COM"
  ]
}
```

```bash
→ ops-cli whois google.com --expiry
2028-09-14T12:00:00+08:00
```

### `wsping`

```bash
→ ops-cli wsping wss://wss.example.com
Connect success
```
