# ops-cli

Try to collect useful tools for ops.

## Installation

### Go Install

```bash
go install github.com/linzeyan/ops-cli@latest
```

### Go Build

```bash
go build . -trimpath -ldflags='-s -w'
```

### [Download Page](https://github.com/linzeyan/ops-cli/releases/latest)


## Usage

```bash
OPS useful tools

Usage:
  ops-cli [flags]
  ops-cli [command]

Available Commands:
  cert        Check tls cert
  completion  Generate the autocompletion script for the specified shell
  dig         Resolve domain name
  geoip       Print IP geographic information
  help        Help about any command
  icp         Check ICP status
  otp         Calculate passcode
  ping        Send ICMP echo packets to host
  qrcode      Read or output QR Code
  random      Generate random string
  url         Expand shorten url
  version     Print version information
  whois       List domain name information

Flags:
  -h, --help   help for ops-cli
  -j, --json   Output JSON format
  -y, --yaml   Output YAML format
```
