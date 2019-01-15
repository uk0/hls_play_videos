#!/usr/bin/env bash
DIR_PATH=$(cd `dirname $0`; pwd)
export DIR_PATH=$(cd `dirname $0`; pwd)

cd  $DIR_PATH



rm -rf release/*

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/clouddisk_linux_amd64  ./


#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./release/clouddisk_windows_amd64.exe  ./

cp -r app/ release/app/

cp ./dev-compose.yaml release/dev-compose.yaml

chmod +x release/clouddisk_linux_amd64


