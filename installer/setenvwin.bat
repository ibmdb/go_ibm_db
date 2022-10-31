:: Script file to set environment variables to use db2cli executable from 
:: clidriver/bin folder
:: This script is only for Windows platform.


echo off
cd ../../clidriver
set clidrvpath=%cd%
set IBM_DB_HOME=%clidrvpath%
set PATH=%PATH%;%IBM_DB_HOME%\bin

