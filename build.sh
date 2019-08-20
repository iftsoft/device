#!/bin/bash

pushd "$(dirname "$0")" > /dev/null || exit

cd ./example/client/ || exit
go build .
cd ../..

cd ./example/server/ || exit
go build .
cd ../..

popd > /dev/null
