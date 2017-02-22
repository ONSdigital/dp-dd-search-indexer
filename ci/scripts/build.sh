#!/bin/bash -eux

export BINPATH=$(pwd)/bin
export GOPATH=$(pwd)/go

pushd $GOPATH/src/github.com/ONSdigital/dp-dd-search-indexer
  go build -o $BINPATH/dp-dd-search-indexer
popd
