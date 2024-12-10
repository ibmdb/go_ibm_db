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





## <a name="inslnx"></a> 2. Go_ibm_db Installation on Linux.

### 2.1 Install GoLang for Linux

Download the
[GoLang Linux binaries](https://golang.org/dl) or [Go Latest binaries](https://go.dev/dl) and
extract the file, for example into `/mygo`:

```
cd /mygo
wget -c https://golang.org/dl/go1.20.5.linux-amd64.tar.gz
tar -xzf go1.20.5.linux-amd64.tar.gz
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
   go install github.com/ibmdb/go_ibm_db/installer@v0.4.3
```

It's Done.

#### 2.2.2 Manual Installation by using git clone.

```
1. mkdir goapp
2. cd goapp
3. git clone https://github.com/ibmdb/go_ibm_db/
```

### 2.3 Download clidriver

Download clidriver in your system, use below command:
go to installer folder where go_ibm_db is downloaded in your system 
(Example: /home/uname/go/src/github.com/ibmdb/go_ibm_db/installer or /home/uname/goapp/go_ibm_db/installer 
where uname is the username) and run setup.go file (go run setup.go)


### 2.4 Set environment variables to clidriver directory path

#### 2.4.1 Manual
```
export IBM_DB_HOME=/home/uname/clidriver
export CGO_CFLAGS=-I$IBM_DB_HOME/include
export CGO_LDFLAGS=-L$IBM_DB_HOME/lib 
export LD_LIBRARY_PATH=/home/uname/clidriver/lib
or
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$IBM_DB_HOME/lib
```

#### 2.4.2 Script file
```
cd .../go_ibm_db/installer
source setenv.sh
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
   go install github.com/ibmdb/go_ibm_db/installer@v0.4.3
```

It's Done.

#### 3.2.2 Manual Installation by using git clone.
```
1. mkdir goapp
2. cd goapp
3. git clone https://github.com/ibmdb/go_ibm_db/
4. cd go_ibm_db/installer
5. go run setup.go
```

### 3.3 Download clidriver

To download clidriver in your system, use below command:
Cd to installer folder where go_ibm_db is downloaded in your system 
(Example: /home/uname/go/src/github.com/ibmdb/go_ibm_db/installer or /home/uname/goapp/go_ibm_db/installer 
where uname is the username) and run setup.go file (`go run setup.go`)

#### 3.3.1 downloaded driver version

* Latest version of clidriver available for MacOS x64 system is: v11.5.9
* By default on Intel Chip Macos, clidriver of v11.5.9 will get downloaded.
* First version of clidriver supported on MacOS ARM64 system is: v12.1.0
* On MacOS M1/M2/M3 Chip system, by default clidriver of v12.1.0 will get downloaded.


### 3.4 Set environment variables to clidriver directory path

#### 3.4.1 Manual
```
export IBM_DB_HOME=/home/uname/clidriver
export CGO_CFLAGS=-I$IBM_DB_HOME/include
export CGO_LDFLAGS=-L$IBM_DB_HOME/lib

export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:/home/uname/go/src/github.com/ibmdb/clidriver/lib
or
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:$IBM_DB_HOME/lib
```

#### 3.4.2 Script file
```
cd .../go_ibm_db/installer
source setenv.sh
```

#### 3.4.3 Disable SIP or create symlink of libdb2 on MacARM64 sysem

* New MacOS systems comes with System Integrity Protection(SIP) enabled which discards setting of DYLD_LIBRARY_PATH env variable
* Disable SIP if your Go app gives error that: file `libdb2.dylib` not found.
* If you can not disable SIP, create softlink of `.../clidriver/lib/libdb2.dylib` file under your applications home directory.
```
    ln -s .../clidriver/lib/libdb2.dylib libdb2.dylib
```
* If you see `libdb2.dylib` file under `go_ibm_db` directory, you can copy it too in your app root directory.

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
   go install github.com/ibmdb/go_ibm_db/installer@v0.4.3
```

#### 4.2.2 Manual Installation by using git clone.
```
1. mkdir goapp
2. cd goapp
3. git clone https://github.com/ibmdb/go_ibm_db/
```

### 4.3 Download clidriver

Download clidriver in your system, go to installer folder where go_ibm_db is downloaded in your system, use below command: 
(Example: C:\Users\uname\go\src\github.com\ibmdb\go_ibm_db\installer or C:\goapp\go_ibm_db\installer 
 where uname is the username ) and run setup.go file (go run setup.go).


### 4.4 Set environment variables to clidriver directory path

#### 4.4.1 Manual
```
set IBM_DB_HOME=C:\Users\uname\go\src\github.com\ibmdb\clidriver
set PATH=%PATH%;C:\Users\uname\go\src\github.com\ibmdb\clidriver\bin
or 
set PATH=%PATH%;%IBM_DB_HOME%\bin
```

### 4.4.2 Script file 
```
cd .../go_ibm_db/installer
Run setenvwin.bat
```
It's Done.

4. Download platform specific clidriver from https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/ , untar/unzip it and set `IBM_DB_HOME` environmental variable to full path of extracted 'clidriver' directory, for example if clidriver is extracted as: `/home/mysystem/clidriver`, then set system level environment variable `IBM_DB_HOME=/home/mysystem/clidriver`.

