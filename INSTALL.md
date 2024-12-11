# Installing go_ibm_db

*Copyright (c) 2014 IBM Corporation and/or its affiliates. All rights reserved.*

Permission is hereby granted, free of charge, to any person obtaining a copy of this
software and associated documentation files (the "Software"), to deal in the Software
without restriction, including without limitation the rights to use, copy, modify, 
merge, publish, distribute, sublicense, and/or sell copies of the Software, 
and to permit persons to whom the Software is furnished to do so, subject to the 
following conditions:

The above copyright notice and this permission notice shall be included in all copies 
or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR 
PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE 
FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR 
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER 
DEALINGS IN THE SOFTWARE.

## Contents

1. [Overview](#Installation)
2. [ibm_db Installation on Linux](#inslnx)
3. [ibm_db Installation on MacOS](#insmac)
4. [ibm_db Installation on Windows](#inswin)

## <a name="overview"></a> 1. Overview

The [*go_ibm_db*](https://github.com/ibmdb/go_ibm_db) is an asynchronous/synchronous interface for GoLang to IBM DB2.

Following are the steps to installation in your system.

This go_ibm_db driver has been tested on 64-bit/32-bit IBM Linux, MacOS and Windows.

### 1.1 clidriver info for MacOS

* Latest version of clidriver available for MacOS x64 system is: v11.5.9
* By default on Intel Chip Macos, clidriver of v11.5.9 will get downloaded.
* First version of clidriver supported on MacOS ARM64 system is: v12.1.0
* On MacOS M1/M2/M3 Chip system, by default clidriver of v12.1.0 will get downloaded.

### 1.2 License requirement to connect to Db2 for z/OS and Db2 for iSeries servers

* Please read [this doc](https://github.com/ibmdb/go_ibm_db/blob/master/README.md#for-zos-and-iseries-connectivity-and-sql1598n-error) for detail info about license requiremnet and resolving SQL1598N error during connection.
* clidriver v12.1.0 requires db2connect v12.1 license to connect z/OS or iSeries severs.
* MacOS Silicon Chip (arm64 processor) is supported using v12.1 clidriver only and  hence require db2connect v12.1 license.
* You can force go_ibm_db driver to use older version of clidirver by setting system level environment varialbe CLIDRIVER_DOWNLOAD_VERSION or explicitly setting IBM_DB_DOWNLOAD_URL to point path of clidriver.tar.gz file.
```
    export CLIDRIVER_DOWNLOAD_VERSION=v11.5.9
    export IBM_DB_DOWNLOAD_URL=https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/v11.5.9/linuxx64_odbc_cli.tar.gz
    go run setup.go
```

## <a name="inslnx"></a> 2. Go_ibm_db Installation on Linux.

### 2.1 Install GoLang for Linux

Download the
[GoLang Linux binaries](https://golang.org/dl) or [Go Latest binaries](https://go.dev/dl) and
extract the file, for example into `$HOME/mygo`:

```
cd $HOME/mygo
wget -c https://golang.org/dl/go1.22.1.linux-amd64.tar.gz
tar -xzf go1.22.1.linux-amd64.tar.gz
export GOROOT=$HOME/mygo/go
export GOPATH=$HOME/mygo
```

Set PATH to include Go:

```
export PATH=/mygo/go/bin:$PATH
```

### 2.2 Install go_ibm_db

Following are the steps to install [*go_ibm_db*](https://github.com/ibmdb/go_ibm_db) from github.
using directory `/goapp` for example.

#### 2.2.1 Direct Installation.
```
1. mkdir goapp
2. cd goapp
3. go install github.com/ibmdb/go_ibm_db/installer@latest
   or
   go install github.com/ibmdb/go_ibm_db/installer@v0.5.2
4. ls $GOPATH/pkg/mod/github.com/ibmdb
   go_ibm_db@v0.5.2
5. cd $GOPATH/pkg/mod/github.com/ibmdb/go_ibm_db@v0.5.2/installer
6. go run ./setup.go
7. export IBM_DB_HOME=$GOPATH/pkg/mod/github.com/ibmdb/clidriver
8. source ./setenv.sh
```

It's Done.

#### 2.2.2 Manual Installation by using git clone.

```
1. mkdir $HOME/goapp
2. cd goapp
3. git clone https://github.com/ibmdb/go_ibm_db/
4. go env GOPATH
5. cd go_ibm_db/installer
6. go run ./setup.go
7. export IBM_DB_HOME=$HOME/goapp/clidriver
8. source ./setenv.sh
```

If IBM_DB_HOME is already set or, sourcing setenv.sh fails, create below environment variables:
```
export CGO_CFLAGS=-I$IBM_DB_HOME/include
export CGO_LDFLAGS=-L$IBM_DB_HOME/lib 
export LD_LIBRARY_PATH=$IBM_DB_HOME/lib:$LD_LIBRARY_PATH
```

## <a name="insmac"></a> 3. Go_ibm_db Installation on MacOS x64 and arm64 Systems

### 3.1 Install GoLang for Mac

Download the
[GoLang MacOS binaries](https://golang.org/dl) or [GoLang Latest binaries](https://go.dev/dl) and
extract the file.

### 3.2 Install go_ibm_db

#### 3.2.1 Direct Installation.
```
1. mkdir goapp
2. cd goapp
3. go install github.com/ibmdb/go_ibm_db/installer@latest
   or
   go install github.com/ibmdb/go_ibm_db/installer@v0.5.2
4. go env GOPATH
5. cd $GOPATH/pkg/mod/github.com/ibmdb/go_ibm_db@v0.5.2/installer
6. go run setup.go
7. export IBM_DB_HOME=$GOPATH/pkg/mod/github.com/ibmdb/clidriver
8. source ./setenv.sh
```

It's Done.

#### 3.2.2 Manual Installation by using git clone.
```
1. mkdir goapp
2. cd goapp
3. git clone https://github.com/ibmdb/go_ibm_db/
4. cd go_ibm_db/installer
5. go run setup.go
6. source ./setenv.sh
7. cd ../testdata
8. Edit config.json file and update database connection info, save it.
9. go mod init testdata
10. go mod tidy
11. go run main.go
```

### 3.3 Set environment variables to clidriver directory path

#### 3.3.1 Manual
```
export IBM_DB_HOME=/home/uname/clidriver
export CGO_CFLAGS=-I$IBM_DB_HOME/include
export CGO_LDFLAGS=-L$IBM_DB_HOME/lib
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:$IBM_DB_HOME/lib
```

#### 3.3.2 Script file
```
cd .../go_ibm_db/installer
source setenv.sh
```

#### 3.3.3 Disable SIP or create symlink of libdb2 on MacARM64 sysem

* New MacOS systems comes with System Integrity Protection(SIP) enabled which discards setting of DYLD_LIBRARY_PATH env variable
* Disable SIP if your Go app gives error that: file `libdb2.dylib` not found.
* OR, if you can not disable SIP, create softlink of `.../clidriver/lib/libdb2.dylib` file under your applications home directory.
```
    ln -s .../clidriver/lib/libdb2.dylib libdb2.dylib
```
* OR, if you see `libdb2.dylib` file under `go_ibm_db` directory, you can copy it too in your app root directory.

## <a name="inswin"></a> 4. Go_ibm_db Installation on Windows.

### 4.1 Install GoLang for Windows

Download the [Go Windows binary/installer](https://golang.org/dl) or [Go Latest binaries](https://go.dev/dl/) and
install it.

### 4.2 Install go_ibm_db

Following are the steps to install [*go_ibm_db*](https://github.com/ibmdb/go_ibm_db) from github.
using directory `/goapp` for example.

#### 4.2.1 Direct Installation.
```
1. mkdir goapp
2. cd gopapp
3. go install github.com/ibmdb/go_ibm_db/installer@latest
   or
   go install github.com/ibmdb/go_ibm_db/installer@v0.5.2
4. go env GOPATH
5. cd %GOPATH%\pkg\mod\github.com\ibmdb\go_ibm_db@v0.5.2\installer
6. go run setup.go
7.
    set IBM_DB_HOME=%GOPATH%\pkg\mod\github.com\ibmdb\clidriver
    set CGO_CFLAGS=-I%IBM_DB_HOME%\include
    set CGO_LDFLAGS=-L%IBM_DB_HOME%\lib
    set LIB=%IBM_DB_HOME%\lib;%LIB%
```

#### 4.2.2 Manual Installation by using git clone.
```
1. mkdir %HOME%\goapp
2. cd goapp
3. git clone https://github.com/ibmdb/go_ibm_db/
4. go env GOPATH
5. cd go_ibm_db\installer
6. go run setup.go
7.
    set IBM_DB_HOME=%HOME%\goapp\clidriver
    set CGO_CFLAGS=-I%IBM_DB_HOME%\include
    set CGO_LDFLAGS=-L%IBM_DB_HOME%\lib
    set LIB=%IBM_DB_HOME%\lib;%LIB%
```

### 4.3 Script file
```
cd .../go_ibm_db/installer
Run setenvwin.bat
```
It's Done.

4. Download platform specific clidriver from https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/ , untar/unzip it and set `IBM_DB_HOME` environmental variable to full path of extracted 'clidriver' directory, for example if clidriver is extracted as: `/home/mysystem/clidriver`, then set system level environment variable `IBM_DB_HOME=/home/mysystem/clidriver`.

