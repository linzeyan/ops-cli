#!/usr/bin/env bash

set -ex

VERSION="$(git describe --tags)"
PACKAGE=$(basename ${PWD})
TARGET="obj"

build() {
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0
    go build -a -trimpath -o ${TARGET}/${PACKAGE}_${VERSION}_${GOOS}_${GOARCH} cmd/main.go
    # Linux
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0
    go build -a -trimpath -o ${TARGET}/${PACKAGE}_${VERSION}_${GOOS}_${GOARCH} cmd/main.go
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0
    go build -a -trimpath -o ${TARGET}/${PACKAGE}_${VERSION}_${GOOS}_${GOARCH} cmd/main.go
    # Windows
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0
    go build -a -trimpath -o ${TARGET}/${PACKAGE}_${VERSION}_${GOOS}_${GOARCH} cmd/main.go
}

convert() {
    mkdir upx
    for i in $(ls ${TARGET}); do upx -9 -o upx/${i} ${TARGET}/${i}; done
    rm -rf ${TARGET}
}

clean() {
    go clean
    rm -rf ${TARGET}
}

$1
