#!/bin/bash

export GOPATH=`pwd`
export GOBIN=`pwd`/bin

go get -d gochoose-example
go install gochoose-example

