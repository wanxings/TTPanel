#!/bin/bash
set -e

rm -rf releases/amd64
rm -rf releases/arm64
mkdir -p releases/amd64
mkdir -p releases/arm64

#编译amd64包
GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o TTPanel_amd64
mv TTPanel_amd64 releases/amd64/TTPanel
echo "amd64架构编译完成"
#编译arm64包
GOARCH=arm64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o TTPanel_arm64
mv TTPanel_arm64 releases/arm64/TTPanel
echo "arm64架构编译完成"
