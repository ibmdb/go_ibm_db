:: Script file to set environment variables to use db2cli executable from 
:: clidriver/bin folder
:: This script is only for Windows platform.


echo off
cd ../../clidriver
set clidrvpath=%cd%
set IBM_DB_HOME=%clidrvpath%
set PATH=%IBM_DB_HOME%\bin;%PATH%
set LIB=%IBM_DB_HOME%\lib;%LIB%
set CGO_CFLAGS=%IBM_DB_HOME%\include
set CGO_LDFLAGS=%IBM_DB_HOME%\lib

