#!/bin/bash

pushd "$(dirname "$0")" > /dev/null || exit

GOOS=linux GOARCH=amd64 go build .

popd > /dev/null || exit
