# go-ibm_db: Go Lang Driver for Db2 and Informix database servers.

The C API for DB2 is called CLI and is basically the same interface as
ODBC. However, it can be used without configuring ODBC which
simplifies usage for DB2-only projects.

This driver is based on code.google.com/p/odbc.

## Building

The following cgo environment variables must be set before building this
package:

* CGO_LDFLAGS
* CGO_CFLAGS

Here is a sample script which sets these variables:
```
    #!/bin/bash

    DB2HOME=$HOME/sqllib
    export CGO_LDFLAGS=-L$DB2HOME/lib
    export CGO_CFLAGS=-I$DB2HOME/include

    go build .
```
## Running

On Linux, `LD_LIBRARY_PATH` probably needs to be set to find the DB2
CLI libraries. Other platforms may have similar requirements.

This sample script demonstrates the bare minimum needed:
```
    #!/bin/bash

    DB2HOME=$HOME/sqllib
    export LD_LIBRARY_PATH=$DB2HOME/lib

    ./dbtest -conn 'DATABASE=db; HOSTNAME=dbhost; PORT=40000; PROTOCOL=TCPIP; UID=me; PWD=secret;'
```

### Note

On some DB2 instances (e.g. z/OS) you may have to connect to DB2-Connect which will forward connection requests to DB2.
In theses cases you may run into something like:

    SQLExecute: {42601} [IBM][CLI Driver][DB2] SQL0104N  An unexpected token " " was found following "". 
    Expected tokens may include:  ". <IDENTIFIER> JOIN INNER LEFT RIGHT FULL CROSS , HAVING GROUP".  SQLSTATE=42601
	
Although not really obvious, this means that a terminator is missing for SQL statements.
This may be due to a different parsing approach when DB2-Connect is involved.
If you terminate your SQL statements with ';' you should be fine.
Most of the times though you will connect directly to DB2 and SQL statements without ';' terminator work.

## Sample program
```
    package main

    import (
        _"go-ibm_db"
        "database/sql"
        "flag"
        "fmt"
        "os"
        "time"
    )

    var (
        connStr = flag.String("conn", "", "connection string to use")
        repeat  = flag.Uint("repeat", 1, "number of times to repeat query")
    )

    func usage() {
        fmt.Fprintf(os.Stderr, `usage: %s [options]
```
%s connects to DB2 and executes a simple SQL statement a configurable number of times.

Here is a sample connection string:

```
    DATABASE=MYDBNAME; HOSTNAME=localhost; PORT=60000; PROTOCOL=TCPIP; UID=username; PWD=password;
    , os.Args[0], os.Args[0])
        flag.PrintDefaults()
        os.Exit(1)
    }

    func execQuery(st *sql.Stmt) error {
        rows, err := st.Query()
        if err != nil {
            return err
        }
        defer rows.Close()
        for rows.Next() {
            var t time.Time
            err = rows.Scan(&t)
            if err != nil {
                return err
            }
            fmt.Printf("Time: %v\n", t)
        }
        return rows.Err()
    }

    func dbOperations() error {
        db, err := sql.Open("ibm_db", *connStr)
        if err != nil {
            return err
        }
        defer db.Close()
        // Attention: If you have to go through DB2-Connect you have to terminate SQL-statements with ';'
        st, err := db.Prepare("select current timestamp from sysibm.sysdummy1;")
        if err != nil {
            return err
        }
        defer st.Close()

        for i := 0; i < int(*repeat); i++ {
            err = execQuery(st)
            if err != nil {
                return err
            }
        }
        return nil
    }

    func main() {
        flag.Usage = usage
        flag.Parse()
        if *connStr == "" {
            fmt.Fprintln(os.Stderr, "-conn is required")
            flag.Usage()
        }

        if err := dbOperations(); err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    }
```

## CONTRIBUTORS

* http://tip.golang.org/AUTHORS
* [IBM](https://www.ibm.com)
* [Rocket Software](https://www.rocketsoftware.com/)

## Contributing to the go-ibm_db

[Contribution Guidelines](https://github.com/ibmdb/go-ibm_db/blob/master/CONTRIBUTING.md)

```
Contributor should add a reference to the DCO sign-off as comment in the pull request(example below):
DCO 1.1 Signed-off-by: Random J Developer <random@developer.org>
```

## License
[LICENSE](https://github.com/ibmdb/go-ibm_db/blob/master/LICENSE.md)

