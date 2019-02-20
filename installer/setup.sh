#!/bin/bash

go run setup.go
a=`pwd`
export DB2HOME=$a/clidriver
export CGO_CFLAGS=-I$DB2HOME/include
export CGO_LDFLAGS=-L$DB2HOME/lib
export LD_LIBRARY_PATH=$a/clidriver/lib
