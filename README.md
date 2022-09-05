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
