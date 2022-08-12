#!/usr/bin/env bash

set -e

runCommand="./ops-cli"
assests="test/assets/"
blackHole="/dev/null"
fileDst="/tmp/temp/"
mkdir -p ${fileDst}
trap "rm -f ${runCommand};rm -rf ${fileDst}" EXIT

__build() {
    CGO_ENABLED=0 go build -trimpath -ldflags '-w -s' .
}

__cert() {
    local -r subCommand="cert"
    local -r testHost="www.google.com"
    local -r testArgs="days dns expiry ip issuer"
    local testCommand
    for i in $(echo ${testArgs}); do
        testCommand="${runCommand} ${subCommand} ${testHost} --${i}"
        ${testCommand} >${blackHole}
        sleep 1
    done
}

__convert() {
    local -r subCommand="convert"
    local -r testSubCommand="yaml2json yaml2toml"
    local -r testSrc="${assests}proxy.yaml"
    local testCommand
    for i in $(echo ${testSubCommand}); do
        testCommand="${runCommand} ${subCommand} ${i} -i ${testSrc} -o ${fileDst}${i}.txt"
        ${testCommand}
        sleep 1
    done
}

__dig() {
    local -r subCommand="dig"
    local -r testHost="google.com"
    local -r testArgs="A AAAA CNAME NS"
    local testCommand
    for i in $(echo ${testArgs}); do
        testCommand="${runCommand} ${subCommand} ${testHost} @1.1.1.1 ${i}"
        ${testCommand} >${blackHole}
        sleep 1
    done
}

__doc() {
    local -r subCommand="doc"
    local -r testSubCommand="man markdown rest yaml"
    local -r testDst="${fileDst}docs"
    local testCommand
    for i in $(echo ${testSubCommand}); do
        testCommand="${runCommand} ${subCommand} ${i} -d ${testDst}"
        ${testCommand}
        sleep 1
    done
}

__geoip() {
    local -r subCommand="geoip"
    local -r testHost="9.9.9.9 1.1.1.1"
    local -r testArgs="j y"
    local testCommand
    for i in $(echo ${testArgs}); do
        for j in $(echo ${testHost}); do
            testCommand="${runCommand} ${subCommand} ${j} -${i}"
            ${testCommand} >${blackHole}
            sleep 1
        done
    done
}

__otp() {
    local -r subCommand="otp"
    local -r testArgs1="calculate 6BDR T7AT RRCZ V5IS FLOH AHQL YF4Z ORG7"
    local -r testArgs2="calculate T7L756M2FEL6CHISIXVSGT4VUDA4ZLIM -p 15 -d 7"
    local -r testArgs3="generate"
    local testCommand
    testCommand="${runCommand} ${subCommand} ${testArgs1}"
    ${testCommand} >${blackHole}
    sleep 1
    testCommand="${runCommand} ${subCommand} ${testArgs2}"
    ${testCommand} >${blackHole}
    sleep 1
    testCommand="${runCommand} ${subCommand} ${testArgs3}"
    ${testCommand} >${blackHole}
}

__qrcode() {
    local -r subCommand="qrcode"
    local -r testArgs1="text https://www.google.com -o ${fileDst}out.png"
    local -r testArgs2="otp --otp-account my@gmail.com --otp-secret fqowefilkjfoqwie --otp-issuer aws -o ${fileDst}otp.png"
    local -r testArgs3="wifi --wifi-type WPA --wifi-pass your_password --wifi-ssid your_wifi_ssid -o ${fileDst}wifi.png"
    local -r testArgs4="read ${assests}example.png"
    local testCommand
    testCommand="${runCommand} ${subCommand} ${testArgs1}"
    ${testCommand} >${blackHole}
    sleep 1
    testCommand="${runCommand} ${subCommand} ${testArgs2}"
    ${testCommand} >${blackHole}
    sleep 1
    testCommand="${runCommand} ${subCommand} ${testArgs3}"
    ${testCommand} >${blackHole}
}

__random() {
    local -r subCommand="random"
    local -r testSubCommand="lowercase uppercase number symbol"
    local testCommand
    for i in $(echo ${testSubCommand}); do
        testCommand="${runCommand} ${subCommand} ${i} -l 50"
        ${testCommand} >${blackHole}
        sleep 1
    done
    testCommand="${runCommand} ${subCommand}"
    ${testCommand} >${blackHole}
    sleep 1
    testCommand="${runCommand} ${subCommand} -s 10"
    ${testCommand} >${blackHole}
    sleep 1
    testCommand="${runCommand} ${subCommand} -s 10 -o 10 -l 32"
    ${testCommand} >${blackHole}
    sleep 1
    testCommand="${runCommand} ${subCommand} -s 10 -o 10 -u 10 -l 64"
    ${testCommand} >${blackHole}
    sleep 1
    testCommand="${runCommand} ${subCommand} -s 10 -o 10 -u 10 -n 2 -l 39"
    ${testCommand} >${blackHole}
}

__system() {
    local -r subCommand="system"
    local -r testSubCommand="cpu disk host load memory network"
    local testCommand
    for i in $(echo ${testSubCommand}); do
        testCommand="${runCommand} ${subCommand} ${i}"
        ${testCommand} >${blackHole}
        sleep 1
    done
    testCommand="${runCommand} ${subCommand} network -a"
    ${testCommand} >${blackHole}
    sleep 1
    testCommand="${runCommand} ${subCommand} network -i"
    ${testCommand} >${blackHole}
}

__url() {
    local -r subCommand="url"
    local -r testArgs="https://goo.gl/maps/b37Aq3Anc7taXQDd9 https://reurl.cc/7peeZl https://bit.ly/3gk7w5x"
    local testCommand
    for i in $(echo ${testArgs}); do
        testCommand="${runCommand} ${subCommand} ${i}"
        ${testCommand} >${blackHole}
        sleep 1
    done
}

__version() {
    local -r subCommand="version"
    local testCommand
    testCommand="${runCommand} ${subCommand}"
    ${testCommand} >${blackHole}
    sleep 1
    testCommand="${runCommand} ${subCommand} -c"
    ${testCommand} >${blackHole}
}

__whois() {
    local -r subCommand="whois"
    local -r testHost="apple.com"
    local -r testArgs="d e n r"
    local testCommand
    for i in $(echo ${testArgs}); do
        testCommand="${runCommand} ${subCommand} ${testHost} -${i}"
        ${testCommand} >${blackHole}
        sleep 1
    done
    testCommand="${runCommand} ${subCommand} google.com"
    ${testCommand} >${blackHole}
}

__all() {
    # build binary
    __build
    # test subcommand
    __cert
    __convert
    __dig
    __doc
    __geoip
    __otp
    __qrcode
    __random
    __system
    __url
    __version
    __whois
}

__$1
