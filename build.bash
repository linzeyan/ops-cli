#!/usr/bin/env bash

build() {
    local os="$(go env GOOS)"
    local arch="$(go env GOARCH)"
    export CGO_ENABLED=0
    go build -a -trimpath -ldflags="-X 'github.com/linzeyan/ops-cli/cmd.appVersion=$(git describe --tags)' -X 'github.com/linzeyan/ops-cli/cmd.appBuildTime=$(date)' -X 'github.com/linzeyan/ops-cli/cmd.appCommit=$(git rev-parse HEAD)' -X 'github.com/linzeyan/ops-cli/cmd.appPlatform=${os}/${arch}'" -o $(basename ${PWD})_${os}_${arch} .
}

release() {
    export CGO_ENABLED=0 GOOS=darwin GOARCH=amd64
    go build -a -trimpath -ldflags="-X 'github.com/linzeyan/ops-cli/cmd.appVersion=$(git describe --tags)' -X 'github.com/linzeyan/ops-cli/cmd.appBuildTime=$(date)' -X 'github.com/linzeyan/ops-cli/cmd.appCommit=$(git rev-parse HEAD)' -X 'github.com/linzeyan/ops-cli/cmd.appPlatform=${GOOS}/${GOARCH}'" -o $(basename ${PWD})_${GOOS}_${GOARCH} .

    export CGO_ENABLED=0 GOOS=linux GOARCH=amd64
    go build -a -trimpath -ldflags="-X 'github.com/linzeyan/ops-cli/cmd.appVersion=$(git describe --tags)' -X 'github.com/linzeyan/ops-cli/cmd.appBuildTime=$(date)' -X 'github.com/linzeyan/ops-cli/cmd.appCommit=$(git rev-parse HEAD)' -X 'github.com/linzeyan/ops-cli/cmd.appPlatform=${GOOS}/${GOARCH}'" -o $(basename ${PWD})_${GOOS}_${GOARCH} .

    export CGO_ENABLED=0 GOOS=linux GOARCH=arm64
    go build -a -trimpath -ldflags="-X 'github.com/linzeyan/ops-cli/cmd.appVersion=$(git describe --tags)' -X 'github.com/linzeyan/ops-cli/cmd.appBuildTime=$(date)' -X 'github.com/linzeyan/ops-cli/cmd.appCommit=$(git rev-parse HEAD)' -X 'github.com/linzeyan/ops-cli/cmd.appPlatform=${GOOS}/${GOARCH}'" -o $(basename ${PWD})_${GOOS}_${GOARCH} .

    export CGO_ENABLED=0 GOOS=windows GOARCH=amd64
    go build -a -trimpath -ldflags="-X 'github.com/linzeyan/ops-cli/cmd.appVersion=$(git describe --tags)' -X 'github.com/linzeyan/ops-cli/cmd.appBuildTime=$(date)' -X 'github.com/linzeyan/ops-cli/cmd.appCommit=$(git rev-parse HEAD)' -X 'github.com/linzeyan/ops-cli/cmd.appPlatform=${GOOS}/${GOARCH}'" -o $(basename ${PWD})_${GOOS}_${GOARCH} .
}

version() {
    echo "$(git describe --tags)" >version.txt
    echo "BuildTime: $(date)" >>version.txt
    echo "GitCommit: $(git rev-parse HEAD)" >>version.txt
    echo "Platform:  ${GOOS}/${GOARCH}" >>version.txt
}

generate() {
    go generate cmd/*
    export CGO_ENABLED=0
    go build -a -trimpath -o $(basename ${PWD}) .
}

$1
