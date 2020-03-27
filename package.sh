#!/usr/bin/env bash

# get current version
version=$(cat version.txt)
dir_name="Device_Package"

echo "Packaging executable and configs for Linux"
mkdir -p $dir_name
cp version.txt $dir_name/version.txt
# package binary
mkdir -p $dir_name/binary
cp build/linux/client $dir_name/binary/client
cp build/linux/server $dir_name/binary/server
# package scripts
cp -R install/config/ $dir_name/config/
# create archive
tar cfz "build/${dir_name}_linux_${version}.tar.gz" $dir_name/
rm -rf $dir_name/

echo "Packaging executable and configs for Mac OS"
mkdir -p $dir_name
cp version.txt $dir_name/version.txt
# package binary
mkdir -p $dir_name/binary
cp build/darwin/client $dir_name/binary/client
cp build/darwin/server $dir_name/binary/server
# package scripts
cp -R install/config/ $dir_name/config/
# create archive
tar cfz "build/${dir_name}_darwin_${version}.tar.gz" $dir_name/
rm -rf $dir_name/

echo "Packaging executable and configs for Windows"
mkdir -p $dir_name
cp version.txt $dir_name/version.txt
# package binary
mkdir -p $dir_name/binary
cp build/windows/client.exe $dir_name/binary/client.exe
cp build/windows/server.exe $dir_name/binary/server.exe
# package scripts
cp -R install/config/ $dir_name/config/
# create archive
zip -rq "build/${dir_name}_win64_${version}.zip" $dir_name/
rm -rf $dir_name/

echo "Device packages are available under folder $(dirname "$0")/build."

