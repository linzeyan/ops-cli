# ops-cli

[![https://github.com/linzeyan/ops-cli/actions?query=workflow:test](https://github.com/linzeyan/ops-cli/workflows/test/badge.svg?branch=main)](https://github.com/linzeyan/ops-cli/actions?query=workflow:test)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/linzeyan/ops-cli/release?style=flat-square)](https://github.com/linzeyan/ops-cli/actions?query=workflow:release)
[![Release](https://img.shields.io/github/release/linzeyan/ops-cli.svg?style=flat-square)](https://github.com/linzeyan/ops-cli/releases/latest)
[![Software License](https://img.shields.io/github/license/linzeyan/ops-cli?style=flat-square)](./LICENSE)

Try to collect useful tools for ops.

## Installation

### Go Install

```bash
go install github.com/linzeyan/ops-cli@latest
```

### Homebrew

```bash
brew tap linzeyan/tools
brew install ops-cli
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
  cert        Check tls cert expiry time
  completion  Generate the autocompletion script for the specified shell
  convert     Convert data format
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
  otp         Calculate passcode or generate secret
  qrcode      Read or generate QR Code
  random      Generate random string
  ssh-keygen  Generate SSH keypair
  stat        Display file informations
  system      Display system informations
  update      Update ops-cli to the latest release
  url         Get url content or expand shorten url or download
  version     Print version information
  whois       List domain name information

Flags:
      --config string   Specify config path (toml)
  -h, --help            help for ops-cli
  -j, --json            Output JSON format
  -y, --yaml            Output YAML format
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

### `otp`

```bash
→ ops-cli otp calculate 6BDR T7AT RRCZ V5IS FLOH AHQL YF4Z ORG7
631843
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

### `url`

```bash
→ ops-cli url -e https://bit.ly/3gk7w5x
https://zh.wikipedia.org/zh-tw/%E5%8F%8D%E5%90%91%E4%BB%A3%E7%90%86#:~:text=%E5%8F%8D%E5%90%91%E4%BB%A3%E7%90%86%E5%9C%A8%E9%9B%BB%E8%85%A6,%E4%BC%BA%E6%9C%8D%E5%99%A8%E5%8F%A2%E9%9B%86%E7%9A%84%E5%AD%98%E5%9C%A8%E3%80%82
```

### `whois`

```bash
→ ops-cli whois apple.com -j
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
