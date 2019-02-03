#!/bin/sh

set -xe

export GOPATH=$PWD
mkdir -p src/github.com/evcraddock/
cp -R ./source src/github.com/evcraddock/goarticles

BUILD_ROOT=${PWD}/built

go build -o ${BUILD_ROOT}/app/goarticles-api ./src/github.com/evcraddock/goarticles/cmd/goarticles-api

cp ./source/ci/Dockerfile ${BUILD_ROOT}

ls -alR ${BUILD_ROOT}