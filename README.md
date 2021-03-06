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
  ops-cli [flags]
  ops-cli [command]

Available Commands:
  cert        Check tls cert expiry time
  completion  Generate the autocompletion script for the specified shell
  dig         Resolve domain name
  doc         Generate documentation
  geoip       Print IP geographic information
  help        Help about any command
  icp         Check ICP status
  otp         Calculate passcode
  qrcode      Read or generate QR Code
  random      Generate random string
  slack       Send message to slack
  telegram    Send message to telegram
  url         Expand shorten url
  version     Print version information
  whois       List domain name information

Flags:
      --config string   Specify config
  -h, --help            help for ops-cli
  -j, --json            Output JSON format
  -y, --yaml            Output YAML format
```
