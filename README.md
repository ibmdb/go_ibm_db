Package go_ibm_db
=================

The C API for DB2 is called CLI and is basically the same interface as
ODBC. However, it can be used without configuring ODBC which
simplifies usage for DB2-only projects.

This driver is based on code.google.com/p/odbc.

How to clone
=============

go get -d github.com/ibmdb/go_ibm_db


How to build in Windows
=======================

1)To Download Clidriver

go to installer folder in go_ibm_db and run setup.go file (go run setup.go) which will download the clidriver.

2) Now set below env variables:

cond. 1: If you use the clidriver downloaded by the godriver
{
path=C:\Users\uname\go\src\github.com\ibmdb\go_ibm_db\installer\clidriver\bin
}

Else (If you use your cli driver){
path=\Path\To\Clidriver\bin
}



How to build in Linux
======================

1)To Download Clidriver

go to installer folder in go_ibm_db and run setup.go file (go run setup.go) which will download the clidriver.


2) Now set below env variables:

cond. 1: If you use the clidriver downloaded by the godriver
{
export DB2HOME=/home/Users/go/src/github.com/imdb/go_ibm_db/installer/clidriver
export CGO_CFLAGS=-I$DB2HOME/include
export CGO_LDFLAGS=-L$DB2HOME/lib
export LD_LIBRARY_PATH=/home/Users/go/src/github.com/ibmdb/go_ibm_db/installer/clidriver/lib

}

Else (If you use your cli driver){
export DB2HOME=/Path/To/clidriver
export CGO_CFLAGS=-I$DB2HOME/include
export CGO_LDFLAGS=-L$DB2HOME/lib 
export LD_LIBRARY_PATH=/Path/To/clidriver/lib
}


How to build in MacOS
======================

1)To Download Clidriver

go to installer folder in go_ibm_db and run setup.go file (go run setup.go) which will download the clidriver.


2) Now set below env variables:

cond. 1: If you use the clidriver downloaded by the godriver
{
export DB2HOME=/home/Users/go/src/github.com/imdb/go_ibm_db/installer/clidriver
export CGO_CFLAGS=-I$DB2HOME/include
export CGO_LDFLAGS=-L$DB2HOME/lib
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:/home/Users/go/src/github.com/ibmdb/go_ibm_db/installer/clidriver/lib

}

Else(If you use your cli driver) {
export DB2HOME=/Path/To/clidriver
export CGO_CFLAGS=-I$DB2HOME/include
export CGO_LDFLAGS=-L$DB2HOME/lib 
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:/Path/To/clidriver/lib
}


Testing the driver
==================

1) Put your connection string in the main.go file in testdata folder

2) Now run go test command (use go test -v command for details) 


Note
====

On some DB2 instances (e.g. z/OS) you may have to connect to DB2-Connect which will forward connection requests to DB2.
In theses cases you may run into something like:

    SQLExecute: {42601} [IBM][CLI Driver][DB2] SQL0104N  An unexpected token " " was found following "". 
    Expected tokens may include:  ". <IDENTIFIER> JOIN INNER LEFT RIGHT FULL CROSS , HAVING GROUP".  SQLSTATE=42601
	
Although not really obvious, this means that a terminator is missing for SQL statements.
This may be due to a different parsing approach when DB2-Connect is involved.
If you terminate your SQL statements with ';' you should be fine.
Most of the times though you will connect directly to DB2 and SQL statements without ';' terminator work.




How to run sample program
==========================
```
package main

import (
    _ "github.com/ibmdb/go_ibm_db"
    "database/sql"
    "fmt"
)

func Create_Con(con string) *sql.DB{
 db, err:=sql.Open("go_ibm_db", con)
    if err != nil{
        
		fmt.Println(err)
		return nil
	}
	return db
}

//creating a table

func create(db *sql.DB) error{
    _,err:=db.Exec("DROP table ROCKET")
	if(err!=nil){
    _,err:=db.Exec("create table ROCKET(ID varchar(20),NAME varchar(20),LOCATION varchar(20),POSITION varchar(20))")
    if err != nil{
        return err
    }
	} else {
    _,err:=db.Exec("create table ROCKET(ID varchar(20),NAME varchar(20),LOCATION varchar(20),POSITION varchar(20))")
    if err != nil{
        return err
    }
	}
	fmt.Println("TABLE CREATED")
    return nil
}

//inserting row

func insert(db *sql.DB) error{
    st, err:=db.Prepare("Insert into ROCKET(ID,NAME,LOCATION,POSITION) values('3242','mike','hyd','manager')")
    if err != nil{
        return err
    }
    st.Query()
    return nil
}

//this api selects the data from the table and prints it

func display(db *sql.DB) error{
    st, err:=db.Prepare("select * from ROCKET")
    if err !=nil{
        return err
    }
    err=execquery(st)
    if err!=nil{
        return err
    }
    return nil
}


func execquery(st *sql.Stmt) error{
    rows,err :=st.Query()
    if err != nil{
        return err
    }
	cols, _ := rows.Columns()
    fmt.Printf("%s    %s   %s    %s\n",cols[0],cols[1],cols[2],cols[3])
    fmt.Println("-------------------------------------")
    defer rows.Close()
    for rows.Next(){
        var t,x,m,n string
        err = rows.Scan(&t,&x,&m,&n)
        if err != nil{
            return err
        }
        fmt.Printf("%v  %v   %v         %v\n",t,x,m,n)
    }
    return nil
}


func main(){
    con:="HOSTNAME=host;DATABASE=name;PORT=number;UID=username;PWD=password"
	type Db *sql.DB
	var re Db
	re=Create_Con(con)
    err:=create(re)
	if err != nil{
        fmt.Println(err)
    }
    err=insert(re)
    if err != nil{
        fmt.Println(err)
    }
    err=display(re)
    if err != nil{
        fmt.Println(err)
    }
}


How to run sample program for Pool
==================================
package main

import (
    a "github.com/ibmdb/go_ibm_db"
	_"database/sql"
    "fmt"
	"time"
)

func main(){
    con:="HOSTNAME=host;PORT=number;DATABASE=name;UID=username;PWD=password";
	pool:=a.Pconnect()
	
	//SetConnMaxLifetime will atake the value in MINUTES
	db:=pool.Open(con,"SetConnMaxLifetime=1","SetMaxOpenConns=3","SetMaxIdleConns=4")
    st, err:=db.Prepare("Insert into ROCKET values('hi')")
    if err != nil{
        fmt.Println(err)
    }
	st.Query()
	
	db1:=pool.Open(con)
    st1, err:=db1.Prepare("Insert into ROCKET values('hi1')")
    if err != nil{
        fmt.Println(err)
    }
	st1.Query()
	
	db1.Close()
	db.Close()
	pool.Display()
	pool.Release()
	pool.Display()
}


