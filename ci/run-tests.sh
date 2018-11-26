#! /bin/bash

set -e
set -x

export GOPATH=$PWD
mkdir -p src/github.com/evcraddock/
cp -R ./source src/github.com/evcraddock/goarticles

BUILD_ROOT=${PWD}/built

cd src/github.com/evcraddock/goarticles

go test -short -v ./...