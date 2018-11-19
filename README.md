Package go_ibm_db
=================

The C API for DB2 is called CLI and is basically the same interface as
ODBC. However, it can be used without configuring ODBC which
simplifies usage for DB2-only projects.

This driver is based on code.google.com/p/odbc.

How to Download
=============

go get github.com/ibmdb/go_ibm_db


How to run sample program
==========================
```
package main

import (
    _ "github.com/ibmdb/go_ibm_db"
    "database/sql"
    "flag"
    "fmt"
    "os"
)

var (
    connStr = flag.String("conn", "<use your connection string here>", "connection string")
    repeat  = flag.Uint("repeat", 1, "1")
)



func usage() {
fmt.Println("function usage")
    fmt.Println(os.Stderr, `usage: %s [options]

%s connects to DB2 and executes a simple SQL statement a configurable
number of times.

Here is a sample connection string:

DATABASE=MYDBNAME; HOSTNAME=localhost; PORT=60000; PROTOCOL=TCPIP; UID=username; PWD=password;
`, os.Args[0], os.Args[0])
    flag.PrintDefaults()
    os.Exit(1)
}

func execQuery(st *sql.Stmt) error {
    fmt.Println("function execQuery", *st)
    rows, err := st.Query()
    if err != nil {
        return err
    }
    defer rows.Close()
    for rows.Next() {
        var t string
        err = rows.Scan(&t)
        if err != nil {
            return err
        }
        fmt.Printf("Row: %v\n", t)
    }
    return rows.Err()
}

func dbOperations() error {
    fmt.Println("function dbOperations")
    db, err := sql.Open("go_ibm_db", *connStr)
    if err != nil {
        return err
    }
    defer db.Close()
    // Attention: If you have to go through DB2-Connect you have to terminate SQL-statements with ';'
    st, err := db.Prepare("select * from myTable;")
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
	fmt.Println("connection string is : ", *connStr)
	
    if *connStr == "" {
	fmt.Println("inside if")
        fmt.Fprintln(os.Stderr, "-conn is required")
        flag.Usage()
    }

    if err := dbOperations(); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
}