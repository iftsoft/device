#!/bin/bash
set -e

pushd "$(dirname "$0")" > /dev/null || exit

export DOCKER_BUILD_BASE_DIR="$(dirname $0)"
IMAGE_NAME="device-builder"
TAG_NAME="latest"

docker build \
    -t "$IMAGE_NAME":"$TAG_NAME" \
    -f "Dockerfile" .
container=$(docker create "$IMAGE_NAME":"$TAG_NAME")
mkdir -p build/
function cleanup() {
    docker rm $container
}
trap cleanup EXIT

docker cp $container:/gopath/src/github.com/iftsoft/device/build/. build/

popd > /dev/null
