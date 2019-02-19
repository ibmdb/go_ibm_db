cd %GOPATH%
cd src\github.com\ibmdb\go_ibm_db\installer
go run setup.go
setx /M Path "%PATH%%GOPATH%\src\github.com\ibmdb\go_ibm_db\installer\clidriver\bin"
pause