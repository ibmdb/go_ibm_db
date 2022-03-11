# Script file to set environment variables to use db2cli executable from 
# clidriver/bin folder
# This script is only for non-Windows platform.

if [ "$DB2HOME" == "" ]
then
    cd ../../
    DB2HOME=`pwd`/clidriver
fi

OS=`uname`

export CGO_CFLAGS=-I$DB2HOME/include
export CGO_LDFLAGS=-L$DB2HOME/lib

if [ "$OS" == "Darwin" ]
then
    export DYLD_LIBRARY_PATH=$DB2HOME/lib:$DYLD_LIBRARY_PATH
else
    export LD_LIBRARY_PATH=$DB2HOME/lib:$LD_LIBRARY_PATH
fi

export PATH=$PATH:$DB2HOME/bin

