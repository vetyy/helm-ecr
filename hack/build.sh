#!/usr/bin/env sh

projectRoot="$1"
pkg="$2"

if [ ! -e "${GOPATH}/src/${pkg}" ]; then
    mkdir -p $(dirname "${GOPATH}/src/${pkg}")
    ln -sfn "${projectRoot}" "${GOPATH}/src/${pkg}"
fi

version="$(cat plugin.yaml | grep "version" | cut -d '"' -f 2)"

cd "${GOPATH}/src/${pkg}"
go build -o bin/helmecr -ldflags "-X main.version=${version}" ./cmd/helmecr
