#!/usr/bin/env bash

rm -rf completion dist doc

doc() {
    mkdir doc
    go run . doc man
}

completion() {
    mkdir completion
    go run . completion zsh &>completion/ops-cli.zsh
    go run . completion bash &>completion/ops-cli.bash
    go run . completion fish &>completion/ops-cli.fish
}

$1
