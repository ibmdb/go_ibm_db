# go_ibm_db

Interface for GoLang to `DB2 for z/OS`, `DB2 for LUW` and `DB2 for i` database servers.

## API Documentation

> For complete list of go_ibm_db APIs and examples please check [APIDocumentation.md](https://github.com/ibmdb/go_ibm_db/blob/master/API_DOCUMENTATION.md)

## Prerequisite

- Golang should be installed(Golang version should be >=1.12.x and <= 1.21.X)

- Git should be installed in your system.

- For non-windows users, GCC and tar should be present in your system.

- For Docker Linux Container(Ex: Amazon Linux2), use below commands:
```
yum install go git tar libpam
```
### Note:
* Environment variable `DB2HOME` is changed to `IBM_DB_HOME`.

* **SQL1598N Error** - It is expected in absence of valid db2connect license. Please click [here](#Licenserequirements) and read instructions about license requirement and how to apply the license.

## How to Install in Windows

- You may install go_ibm_db using either of below commands

```
go get -d github.com/ibmdb/go_ibm_db
go install github.com/ibmdb/go_ibm_db/installer@latest
go install github.com/ibmdb/go_ibm_db/installer@v0.5.0
```

- You can optionally specify a specific cli driver by setting the IBM_DB_DOWNLOAD_URL environment variable
to the full path of your desired cli driver.  For example, if you want to install the 64-bit macos v11.5.4
cli driver instead of the latest one, set the variable as below:
```
export IBM_DB_DOWNLOAD_URL=https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/v11.5.4/macos64_odbc_cli.tar.gz
```

- If you already have a clidriver available in your system, add the path of the same to your
PATH windows environment variable.
Example:
`set PATH = "C:\Program Files\IBM\IBM DATA SERVER DRIVER\bin";%PATH%`

- If you do not have a clidriver in your system, go to installer folder where `go_ibm_db`
is downloaded in your system, use below command:
(Example: C:\Users\uname\go\src\github.com\ibmdb\go_ibm_db\installer
 or C:\Users\uname\go\pkg\mod\github.com\ibmdb\go_ibm_db\installer
 where uname is the username ) and run setup.go file (`go run setup.go`).

```
Set IBM_DB_HOME to clidriver downloaded path and
set this path to your PATH windows environment variable
(Example: Path=C:\Users\uname\go\src\github.com\ibmdb\clidriver)
set IBM_DB_HOME=C:\Users\uname\go\src\github.com\ibmdb\clidriver
set PATH=%PATH%;C:\Users\uname\go\src\github.com\ibmdb\clidriver\bin
or 
set PATH=%PATH%;%IBM_DB_HOME%\bin
```

- Script file to set environment variable
```
cd .../go_ibm_db/installer
setenvwin.bat
```

## How to Install in Linux/Mac

- You may install go_ibm_db using either of below commands

```
go get -d github.com/ibmdb/go_ibm_db
go install github.com/ibmdb/go_ibm_db/installer@latest
go install github.com/ibmdb/go_ibm_db/installer@v0.5.0
```

- You can optionally specify a specific cli driver by setting the IBM_DB_DOWNLOAD_URL environment
variable to the full path of your desired driver. For example, if you want to install the 64-bit
macos v11.5.4 cli driver instead of the latest one, set the variable as below:
```
export IBM_DB_DOWNLOAD_URL=https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/v11.5.4/macos64_odbc_cli.tar.gz
```

- If you already have a clidriver available in your system, set the below environment variables with the clidriver path
```
export IBM_DB_HOME=/home/uname/clidriver
export CGO_CFLAGS=-I$IBM_DB_HOME/include
export CGO_LDFLAGS=-L$IBM_DB_HOME/lib 
Linux:
export LD_LIBRARY_PATH=/home/uname/clidriver/lib
or
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$IBM_DB_HOME/lib
Mac:
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:/Applications/clidriver/lib
```

- If you do not have a clidriver available in your system, use below commands.
(Example: /home/uname/go/src/github.com/ibmdb/go_ibm_db/installer
 or /home/uname/go/pkg/mod/github.com/ibmdb/go_ibm_db/installer
where uname is the username) and run setup.go file (`go run setup.go`)

```
export IBM_DB_HOME=/home/uname/go/src/github.com/ibmdb/clidriver
export CGO_CFLAGS=-I$IBM_DB_HOME/include
export CGO_LDFLAGS=-L$IBM_DB_HOME/lib
Linux:
export LD_LIBRARY_PATH=/home/uname/go/src/github.com/ibmdb/clidriver/lib
or
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$IBM_DB_HOME/lib
Mac:
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:/home/uname/go/src/github.com/ibmdb/clidriver/lib
or
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:$IBM_DB_HOME/lib
```

- Script file to set environment variables in Linux/Mac

```
cd .../go_ibm_db/installer
source setenv.sh
```

- For Docker Linux Container, use below commands:
```
yum install -y gcc git go wget tar xz make gcc-c++
cd /root
curl -OL https://golang.org/dl/go1.21.4.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz

rm /usr/bin/go
rm /usr/bin/gofmt
cp /usr/local/go/bin/go /usr/bin/
cp /usr/local/go/bin/gofmt /usr/bin/

go install github.com/ibmdb/go_ibm_db/installer@v0.5.0
or 
go install github.com/ibmdb/go_ibm_db/installer@latest
```

## How to Install in z/OS

- You may install go_ibm_db using the below command

```
go install github.com/ibmdb/go_ibm_db/installer@latest
```

### Configure ODBC driver on z/OS

> **Note**
> You must have the following environment variables set.
> * IBM_DB_HOME
> * STEPLIB
> * DSNAOINI


Please refer to the [ODBC Guide and References](https://www.ibm.com/support/knowledgecenter/SSEPEK/pdf/db2z_12_odbcbook.pdf) cookbook for how to configure your ODBC driver.   Specifically, you need to ensure you have:

1. Apply Db2 on z/OS PTF [UI60551](https://www-01.ibm.com/support/docview.wss?uid=swg1PH05953) to pick up new ODBC functionality to support Node.js applications.

2. Binded the ODBC packages.  A sample JCL is provided in the `SDSNSAMP` dataset in member `DSNTIJCL`.  Customize the JCL with specifics to your system.

3. Set the `IBM_DB_HOME` environment variable to the High Level Qualifier (HLQ) of your Db2 datasets.  For example, if your Db2 datasets are located as `DSNC10.SDSNC.H` and `DSNC10.SDSNMACS`, you need to set `IBM_DB_HOME` environment variable to `DSNC10` with the following statement (can be saved in `~/.profile`):

    ```sh
    # Set HLQ to Db2 datasets.
    export IBM_DB_HOME="DSNC10"
    ```

4. Update the `STEPLIB` environment variable to include the Db2 SDSNEXIT, SDSNLOAD and SDSNLOD2 data sets. You can set the `STEPLIB` environment variable with the following statement, after defining `IBM_DB_HOME` to the high level qualifier of your Db2 datasets as instructed above:

    ```sh
    # Assumes IBM_DB_HOME specifies the HLQ of the Db2 datasets.
    export STEPLIB=$STEPLIB:$IBM_DB_HOME.SDSNEXIT:$IBM_DB_HOME.SDSNLOAD:$IBM_DB_HOME.SDSNLOD2
    ```

5. Configured an appropriate _Db2 ODBC initialization file_ that can be read at application time. You can specify the file by using either a DSNAOINI data definition statement or by defining a `DSNAOINI` z/OS UNIX environment variable.  For compatibility with ibm_db, the following properties must be set:

    In COMMON section:

    ```
    MULTICONTEXT=2
    CURRENTAPPENSCH=ASCII
    FLOAT=IEEE
    ```

    In SUBSYSTEM section:

    ```
    MVSATTACHTYPE=RRSAF
    ```

    Here is a sample of a complete initialization file:

    ```
    ; This is a comment line...
    ; Example COMMON stanza
    [COMMON]
    MVSDEFAULTSSID=VC1A
    CONNECTTYPE=1
    MULTICONTEXT=2
    CURRENTAPPENSCH=ASCII
    FLOAT=IEEE
    ; Example SUBSYSTEM stanza for VC1A subsystem
    [VC1A]
    MVSATTACHTYPE=RRSAF
    PLANNAME=DSNACLI
    ; Example DATA SOURCE stanza for STLEC1 data source
    [STLEC1]
    AUTOCOMMIT=1
    CURSORHOLD=1
    ```

Here's a simple script you can run setup your ODBC INI file.
The following script exepcts that you have have defined your subsystem in following way:

```export SUBSYSTEM=<SSID>```

This script will setup your inistialization file. 

```shell
export DSNAOINI="$HOME/ODBC_${HOSTNAME}_${SUBSYSTEM}_CAF"
touch $DSNAOINI
/bin/cat /dev/null  > "$DSNAOINI"
/bin/chtag -t -c 1047 "$DSNAOINI"
_BPXK_AUTOCVT=ON /bin/cat <<EOF >"$DSNAOINI"
[COMMON]
MVSDEFAULTSSID=$SUBSYSTEM
CURRENTAPPENSCH=UNICODE
FLOAT=IEEE
[$SUBSYSTEM]
MVSATTACHTYPE=CAF
PLANNAME=DSNACLI
[$HOSTNAME$SUBSYSTEM]
AUTOCOMMIT=1
EOF
```

    Reference Chapter 3 in the [ODBC Guide and References](https://www.ibm.com/support/knowledgecenter/SSEPEK/pdf/db2z_12_odbcbook.pdf) for more instructions.


## <a name="Licenserequirements"></a>For z/OS and iSeries Connectivity and SQL1598N error

- Connection to `Db2 for z/OS` or `Db2 for i`(AS400) Server using `ibm_db` driver from distributed platforms (Linux, Unix, Windows and MacOS) is not free. It requires either client side or server side license.

- Connection to `Db2 for LUW` or `Informix` Server using `ibm_db` driver is free.

- `ibm_db` returns SQL1598N error in absence of a valid db2connect license. SQL1598N error is returned by the Db2 Server to client.
To suppress this error, Db2 server must be activated with db2connectactivate utility OR a client side db2connect license file must exist.

- Db2connect license can be applied on database server or client side. A **db2connect license of version 11.5** is required for ibm_db.

- For activating server side license, you can purchase either `Db2 Connect Unlimited Edition for System z®` or `Db2 Connect Unlimited Edition for System i®` license from IBM.

- Ask your DBA to run db2connectactivate utility on Server to activate db2connect license.

- If database Server is enabled for db2connect, no need to apply client side db2connect license.

- If Db2 Server is not db2connectactivated to accept unlimited number of client connection, you must need to apply client side db2connect license.

- db2connectactivate utility and client side db2connect license both comes together from IBM in a single zip file.

- Client side db2connect license is a `db2con*.lic` file that must be copied under `clidriver\license` directory.

- If you have a `db2jcc_license_cisuz.jar` file, it will not work for ibm_db. `db2jcc_license_cisuz.jar` is a db2connect license file for Java Driver. For non-Java Driver, client side db2connect license comes as a file name `db2con*.lic`.

- If environment variable `IBM_DB_HOME` or `IBM_DB_INSTALLER_URL` is not set, `ibm_db` automatically downloads [open source driver specific clidriver](https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/) and save as `github.com/ibmdb/clidriver`. Ignores any other installation.

- If `IBM_DB_HOME` or `IBM_DB_INSTALLER_URL` is set, you need to have same version of db2connect license as installed db2 client. Check db2 client version using `db2level` command to know version of db2connect license required. The license file should get copied under `$IBM_DB_HOME\license` directory.

- If you do not have db2connect license, contact [IBM Customer Support](https://www.ibm.com/mysupport/s/?language=en_US) to buy db2connect license. Find the `db2con*.lic` file in the db2connect license shared by IBM and copy it under `.../github.com/ibmdb/clidriver/license` directory or `$IBM_DB_HOME/license` directory if `IBM_DB_HOME` is set.

- To know more about license and purchasing cost, please contact [IBM Customer Support](https://www.ibm.com/mysupport/s/?language=en_US).

- To know more about server based licensing viz db2connectactivate, follow below links:
* [Activating the license certificate file for Db2 Connect Unlimited Edition](https://www.ibm.com/docs/en/db2/11.5?topic=li-activating-license-certificate-file-db2-connect-unlimited-edition).
* [Unlimited licensing using db2connectactivate utility](https://www.ibm.com/docs/en/db2/11.1?topic=edition-db2connectactivate-server-license-activation-utility).

## How to run sample program

### example1.go:-

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/ibmdb/go_ibm_db"
)

func main() {
	con := "HOSTNAME=host;DATABASE=name;PORT=number;UID=username;PWD=password"
	db, err := sql.Open("go_ibm_db", con)
	if err != nil {
		fmt.Println(err)
	}
	db.Close()
}
```
To run the sample:- go run example1.go

For complete list of connection parameters please check [this.](https://www.ibm.com/docs/en/db2/11.5?topic=file-data-server-driver-configuration-keywords)

### example2.go:-

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/ibmdb/go_ibm_db"
)

func Create_Con(con string) *sql.DB {
	db, err := sql.Open("go_ibm_db", con)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}

// Creating a table.

func create(db *sql.DB) error {
	_, err := db.Exec("DROP table SAMPLE")
	if err != nil {
		_, err := db.Exec("create table SAMPLE(ID varchar(20),NAME varchar(20),LOCATION varchar(20),POSITION varchar(20))")
		if err != nil {
			return err
		}
	} else {
		_, err := db.Exec("create table SAMPLE(ID varchar(20),NAME varchar(20),LOCATION varchar(20),POSITION varchar(20))")
		if err != nil {
			return err
		}
	}
	fmt.Println("TABLE CREATED")
	return nil
}

// Inserting row.

func insert(db *sql.DB) error {
	st, err := db.Prepare("Insert into SAMPLE(ID,NAME,LOCATION,POSITION) values('3242','Mike','Hyderabad','Manager')")
	if err != nil {
		return err
	}
	st.Query()
	return nil
}

// This API selects the data from the table and prints it.

func display(db *sql.DB) error {
	st, err := db.Prepare("select * from SAMPLE")
	if err != nil {
		return err
	}
	err = execquery(st)
	if err != nil {
		return err
	}
	return nil
}

func execquery(st *sql.Stmt) error {
	rows, err := st.Query()
	if err != nil {
		return err
	}
	cols, _ := rows.Columns()
	fmt.Printf("%s    %s   %s    %s\n", cols[0], cols[1], cols[2], cols[3])
	fmt.Println("-------------------------------------")
	defer rows.Close()
	for rows.Next() {
		var t, x, m, n string
		err = rows.Scan(&t, &x, &m, &n)
		if err != nil {
			return err
		}
		fmt.Printf("%v  %v   %v   %v\n", t, x, m, n)
	}
	return nil
}

func main() {
	con := "HOSTNAME=host;DATABASE=name;PORT=number;UID=username;PWD=password"
	type Db *sql.DB
	var re Db
	re = Create_Con(con)
	err := create(re)
	if err != nil {
		fmt.Println(err)
	}
	err = insert(re)
	if err != nil {
		fmt.Println(err)
	}
	err = display(re)
	if err != nil {
		fmt.Println(err)
	}
}
```
To run the sample:- go run example2.go

### example3.go:-(POOLING)

```go
package main

import (
	_ "database/sql"
	"fmt"
	a "github.com/ibmdb/go_ibm_db"
)

func main() {
	// Defining connection string
	// Depending on your connection type
	// you may wish to add: MULTICONTEXT=0
	con := "HOSTNAME=host;PORT=number;DATABASE=name;UID=username;PWD=password"
	pool := a.Pconnect("PoolSize=100")

	// SetConnMaxLifetime will take the value in SECONDS
	db := pool.Open(con, "SetConnMaxLifetime=30")
	st, err := db.Prepare("Insert into SAMPLE values('hi','hi','hi','hi')")
	if err != nil {
		fmt.Println(err)
	}
	st.Query()

	// Here the time out is default.
	db1 := pool.Open(con)
	st1, err := db1.Prepare("Insert into SAMPLE values('hi1','hi1','hi1','hi1')")
	if err != nil {
		fmt.Println(err)
	}
	st1.Query()

	db1.Close()
	db.Close()
	pool.Release()
	fmt.Println("success")
}
```
To run the sample:- go run example3.go

### example4.go:-(POOLING- Limit on the number of connections)

```go
package main

import (
           "database/sql"
           "fmt"
           "time"
          a "github.com/ibmdb/go_ibm_db"
       )

func ExecQuery(st *sql.Stmt) error {
        res, err := st.Query()
        if err != nil {
             fmt.Println(err)
        }
        cols, _ := res.Columns()

        fmt.Printf("%s    %s   %s     %s\n", cols[0], cols[1], cols[2], cols[3])
        defer res.Close()
        for res.Next() {
                    var t, x, m, n string
                    err = res.Scan(&t, &x, &m, &n)
                    fmt.Printf("%v  %v   %v     %v\n", t, x, m, n)
        }
        return nil
}

func main() {
	con := "HOSTNAME=host;PORT=number;DATABASE=name;UID=username;PWD=password"
        pool := a.Pconnect("PoolSize=5")

        ret := pool.Init(5, con)
        if ret != true {
                fmt.Println("Pool initializtion failed")
        }

        for i:=0; i<20; i++ {
                db1 := pool.Open(con, "SetConnMaxLifetime=10")
                if db1 != nil {
                        st1, err1 := db1.Prepare("select * from VMSAMPLE")
                        if err1 != nil {
                                fmt.Println("err1 : ", err1)
                        }else{
                                go func() {
                                        ExecQuery(st1)
                                        db1.Close()
                                }()
                        }
                }
        }
        time.Sleep(30*time.Second)
        pool.Release()
}
```
To run the sample:- go run example4.go

For Running the Tests:
======================

* Connection information must be specified in the environment variables
```
For example, by sourcing the following ENV variables:
export DB2_DATABASE=<Database Name>
export DB2_USER=<Username>
export DB2_PASSWD=<Password>
export DB2_HOSTNAME=<Hostname or IP>
export DB2_PORT=<Database Port>
```

* OR
```
If not using environment variables, update connection information in
go_ibm_db/testdata/config.json file.
```

* Now run go test command (use go test -v command for details) 

* To run a particular test case (use "go test sample_test.go main.go", example "go test Arraystring_test.go main.go")


Logging:
========

* You may enable logging in the go_ibm_db module to trace activities. Logging can be directed to be console or a specified file.

```
# Log on console:
go run sample.go -trace

# Log to a file (e.g., log.txt)
go run sample.go -trace log.txt

```
