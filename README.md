# ops-cli

[![https://github.com/linzeyan/ops-cli/actions?query=workflow:golangci-lint](https://github.com/linzeyan/ops-cli/workflows/golangci-lint/badge.svg?branch=main)](https://github.com/linzeyan/ops-cli/actions?query=workflow:golangci-lint)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/linzeyan/ops-cli/goreleaser?style=flat-square)](https://github.com/linzeyan/ops-cli/actions?query=workflow%3Agoreleaser)
[![Release](https://img.shields.io/github/release/linzeyan/ops-cli.svg?style=flat-square)](https://github.com/linzeyan/ops-cli/releases/latest)
[![Software License](https://img.shields.io/github/license/linzeyan/ops-cli?style=flat-square)](./LICENSE)

Try to collect useful tools for ops.

## Installation

### Go Install

```bash
go install github.com/linzeyan/ops-cli@latest
```

### Go Build

```bash
go build -trimpath -ldflags='-s -w' .
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
  LINE        Send message to LINE
  Slack       Send message to Slack
  Telegram    Send message to Telegram
  cert        Check tls cert expiry time
  completion  Generate the autocompletion script for the specified shell
  convert     Convert data format
  dig         Resolve domain name
  doc         Generate documentation
  geoip       Print IP geographic information
  help        Help about any command
  icp         Check ICP status
  otp         Calculate passcode or generate secret
  qrcode      Read or generate QR Code
  random      Generate random string
  url         Expand shorten url
  version     Print version information
  whois       List domain name information

Flags:
      --config string   Specify config path (toml)
  -h, --help            help for ops-cli
  -j, --json            Output JSON format
  -y, --yaml            Output YAML format
```
