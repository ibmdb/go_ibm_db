@echo off

SETx ENV1 "%DB2HOME%\lib"
SETx ENV2 "%DB2HOME%\include"

echo "%ENV1%"
go run setup.go