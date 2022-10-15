# Script file to set environment variables to use db2cli executable from 
# clidriver/bin folder
# This script is only for non-Windows platform.

if [ "$IBM_DB_HOME" == "" ]
then
    cd ../../
    IBM_DB_HOME=`pwd`/clidriver
fi

OS=`uname`

export CGO_CFLAGS=-I$IBM_DB_HOME/include
export CGO_LDFLAGS=-L$IBM_DB_HOME/lib

if [ "$OS" == "Darwin" ]
then
    export DYLD_LIBRARY_PATH=$IBM_DB_HOME/lib:$DYLD_LIBRARY_PATH
else
    export LD_LIBRARY_PATH=$IBM_DB_HOME/lib:$LD_LIBRARY_PATH
fi

export PATH=$PATH:$IBM_DB_HOME/bin

