#!/bin/bash

VERSION="v0.0.0"

set -e
GOPATH=`go env GOPATH`
ORIGINALPATH=`pwd`

rm -rf ./build

mkdir -p ./build/Linux64 ./build/Win64 ./build/Mac
mkdir -p ./preprocess/util

JQURL="https://github.com/stedolan/jq/releases/download/jq-1.6"
JQ="jq-linux64"
if [ ! -f ./preprocess/util/$JQ ]; then    
    curl -o $JQ -L $JQURL/$JQ && mv $JQ ./preprocess/util/ && chmod 777 ./preprocess/util/$JQ
fi
cp ./preprocess/util/$JQ ./build/Linux64/jq

JQ="jq-win64.exe"
if [ ! -f ./preprocess/util/$JQ ]; then    
    curl -o $JQ -L $JQURL/$JQ && mv $JQ ./preprocess/util/
fi
cp ./preprocess/util/$JQ ./build/Win64/jq.exe

JQ="jq-osx-amd64"
if [ ! -f ./preprocess/util/$JQ ]; then    
    curl -o $JQ -L $JQURL/$JQ && mv $JQ ./preprocess/util/ && chmod 777 ./preprocess/util/$JQ
fi
cp ./preprocess/util/$JQ ./build/Mac/jq

###

# go get

GOARCH=amd64
LDFLAGS="-s -w"

GOOS="linux" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o json-mask
mv json-mask ./build/Linux64/
GOOS="windows" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o json-mask.exe
mv json-mask.exe ./build/Win64/
GOOS="darwin" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o json-mask
mv json-mask ./build/Mac/