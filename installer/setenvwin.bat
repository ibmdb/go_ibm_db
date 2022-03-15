:: Script file to set environment variables to use db2cli executable from 
:: clidriver/bin folder
:: This script is only for Windows platform.


echo off
cd ../../clidriver/bin
set clidrvpath=%cd%
set PATH=%PATH%;%clidrvpath%

