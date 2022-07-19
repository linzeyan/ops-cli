# ops-cli

Try to collect tools.
## Installation

### Go Install

```bash
go install github.com/linzeyan/ops-cli@latest
```
### Go Build

```bash
go build .
```
or
```bash
bash ./build.bash build
```

## Usage

```bash
OPS useful tools

Usage:
  ops-cli [flags]
  ops-cli [command]

Available Commands:
  cert        Check tls cert
  completion  Generate the autocompletion script for the specified shell
  geoip       Print IP geographic information
  help        Help about any command
  icp         Check ICP status
  otp         Calculate passcode
  ping        Send ICMP echo packets to host
  qrcode      Read or output QR Code
  random      Generate random string
  url         Expand shorten url
  whois       List domain name information

Flags:
  -h, --help      help for ops-cli
  -v, --version   version for ops-cli
```
