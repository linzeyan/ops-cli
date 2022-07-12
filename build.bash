#!/usr/bin/env bash

build() {
    export CGO_ENABLED=0 GOOS=darwin GOARCH=amd64
    go build -a -trimpath -ldflags="-X 'github.com/linzeyan/ops-cli/cmd.appVersion=$(git describe --tags)' -X 'github.com/linzeyan/ops-cli/cmd.appBuildTime=$(date)' -X 'github.com/linzeyan/ops-cli/cmd.appCommit=$(git rev-parse HEAD)' -X 'github.com/linzeyan/ops-cli/cmd.appPlatform=${GOOS}/${GOARCH}'" -o $(basename ${PWD})_${GOOS}_${GOARCH} .

    export CGO_ENABLED=0 GOOS=linux GOARCH=amd64
    go build -a -trimpath -ldflags="-X 'github.com/linzeyan/ops-cli/cmd.appVersion=$(git describe --tags)' -X 'github.com/linzeyan/ops-cli/cmd.appBuildTime=$(date)' -X 'github.com/linzeyan/ops-cli/cmd.appCommit=$(git rev-parse HEAD)' -X 'github.com/linzeyan/ops-cli/cmd.appPlatform=${GOOS}/${GOARCH}'" -o $(basename ${PWD})_${GOOS}_${GOARCH} .

    export CGO_ENABLED=0 GOOS=linux GOARCH=arm64
    go build -a -trimpath -ldflags="-X 'github.com/linzeyan/ops-cli/cmd.appVersion=$(git describe --tags)' -X 'github.com/linzeyan/ops-cli/cmd.appBuildTime=$(date)' -X 'github.com/linzeyan/ops-cli/cmd.appCommit=$(git rev-parse HEAD)' -X 'github.com/linzeyan/ops-cli/cmd.appPlatform=${GOOS}/${GOARCH}'" -o $(basename ${PWD})_${GOOS}_${GOARCH} .

    export CGO_ENABLED=0 GOOS=windows GOARCH=amd64
    go build -a -trimpath -ldflags="-X 'github.com/linzeyan/ops-cli/cmd.appVersion=$(git describe --tags)' -X 'github.com/linzeyan/ops-cli/cmd.appBuildTime=$(date)' -X 'github.com/linzeyan/ops-cli/cmd.appCommit=$(git rev-parse HEAD)' -X 'github.com/linzeyan/ops-cli/cmd.appPlatform=${GOOS}/${GOARCH}'" -o $(basename ${PWD})_${GOOS}_${GOARCH} .
}

$1
