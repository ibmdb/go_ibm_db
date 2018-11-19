@echo off

SETx ENV1 "%DB2HOME%\lib"
SETx ENV2 "%DB2HOME%\include"

go run setup.go