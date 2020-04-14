# go_ibm_db

Interface for GoLang to DB2 for z/OS, DB2 for LUW, DB2 for i.

## API Documentation

> For complete list of go_ibm_db APIs and examples please check [APIDocumentation.md](https://github.com/ibmdb/go_ibm_db/blob/master/API_DOCUMENTATION.md)

## Prerequisite

Golang should be installed(Golang version should be >=1.12.x and <= 1.14.X)

Git should be installed in your system.

For non-windows users, GCC and tar should be present in your system.

```
For Docker Linux Container(Ex: Amazon Linux2), use below commands:
yum install go git tar libpam
```

## How to Install in Windows
```
go get -d github.com/ibmdb/go_ibm_db

If you already have a cli driver available in your system, add the path of the same to your Path windows environment variable
Example: Path = C:\Program Files\IBM\IBM DATA SERVER DRIVER\bin


If you do not have a clidriver in your system, go to installer folder where go_ibm_db is downloaded in your system (Example: C:\Users\uname\go\src\github.com\ibmdb\go_ibm_db\installer) and run setup.go file (go run setup.go).

where uname is the username

Above command will download clidriver.

Add the path of the clidriver downloaded to your Path windows environment variable
(Example: Path=C:\Users\rakhil\go\src\github.com\ibmdb\go_ibm_db\installer\clidriver\bin)


```

## How to Install in Linux/Mac
```
go get -d github.com/ibmdb/go_ibm_db

If you already have a cli driver available in your system, set the below environment variables with the clidriver path

export DB2HOME=/home/rakhil/dsdriver
export CGO_CFLAGS=-I$DB2HOME/include
export CGO_LDFLAGS=-L$DB2HOME/lib 
Linux:
export LD_LIBRARY_PATH=/home/rakhil/dsdriver/lib
Mac:
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:/Applications/dsdriver/lib

If you do not have a clidriver available in your system
go to installer folder where go_ibm_db is downloaded in your system (Example: /home/uname/go/src/github.com/ibmdb/go_ibm_db/installer) and run setup.go file (go run setup.go)
where uname is the username

Above command will download clidriver.

Set the below envronment variables with the path of the clidriver downloaded

export DB2HOME=/home/uname/go/src/github.com/ibmdb/go_ibm_db/installer/clidriver
export CGO_CFLAGS=-I$DB2HOME/include
export CGO_LDFLAGS=-L$DB2HOME/lib
Linux:
export LD_LIBRARY_PATH=/home/uname/go/src/github.com/ibmdb/go_ibm_db/installer/clidriver/lib
Mac:
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:/home/uname/go/src/github.com/ibmdb/go_ibm_db/installer/clidriver/lib


```

### <a name="Licenserequirements"></a> License requirements for connecting to databases

go_ibm_db driver can connect to DB2 on Linux Unix and Windows without any additional license/s, however, connecting to databases on DB2 for z/OS or DB2 for i(AS400) Servers require either client side or server side license/s. The client side license would need to be copied under `license` folder of your `cidriver` installation directory and for activating server side license, you would need to purchase DB2 Connect Unlimited for System z® and DB2 Connect Unlimited Edition for System i®.

To know more about license and purchasing cost, please contact [IBM Customer Support](http://www-05.ibm.com/support/operations/zz/en/selectcountrylang.html).

To know more about server based licensing viz db2connectactivate, follow below links:
* [Activating the license certificate file for DB2 Connect Unlimited Edition](https://www.ibm.com/developerworks/community/blogs/96960515-2ea1-4391-8170-b0515d08e4da/entry/unlimited_licensing_in_non_java_drivers_using_db2connectactivate_utlility1?lang=en).
* [Unlimited licensing using db2connectactivate utility](https://www.ibm.com/developerworks/community/blogs/96960515-2ea1-4391-8170-b0515d08e4da/entry/unlimited_licensing_in_non_java_drivers_using_db2connectactivate_utlility1?lang=en.)

## How to run sample program

### example1.go:-

```
package main

import (
    _ "github.com/ibmdb/go_ibm_db"
    "database/sql"
    "fmt"
)

func main(){
    con:="HOSTNAME=host;DATABASE=name;PORT=number;UID=username;PWD=password"
 db, err:=sql.Open("go_ibm_db", con)
    if err != nil{
        
		fmt.Println(err)
	}
	db.Close()
}
To run the sample:- go run example1.go
```

For complete list of connection parameters please check [this.](https://www.ibm.com/support/knowledgecenter/SSEPGG_11.1.0/com.ibm.swg.im.dbclient.config.doc/doc/c0054698.html)

### example2.go:-

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
    _,err:=db.Exec("DROP table SAMPLE")
	if(err!=nil){
    _,err:=db.Exec("create table SAMPLE(ID varchar(20),NAME varchar(20),LOCATION varchar(20),POSITION varchar(20))")
    if err != nil{
        return err
    }
	} else {
    _,err:=db.Exec("create table SAMPLE(ID varchar(20),NAME varchar(20),LOCATION varchar(20),POSITION varchar(20))")
    if err != nil{
        return err
    }
	}
	fmt.Println("TABLE CREATED")
    return nil
}

//inserting row

func insert(db *sql.DB) error{
    st, err:=db.Prepare("Insert into SAMPLE(ID,NAME,LOCATION,POSITION) values('3242','mike','hyd','manager')")
    if err != nil{
        return err
    }
    st.Query()
    return nil
}

//this api selects the data from the table and prints it

func display(db *sql.DB) error{
    st, err:=db.Prepare("select * from SAMPLE")
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
To run the sample:- go run example2.go
```

### example3.go:-(POOLING)

```
package main

import (
    a "github.com/ibmdb/go_ibm_db"
	_"database/sql"
    "fmt"
)

func main(){
    con:="HOSTNAME=host;PORT=number;DATABASE=name;UID=username;PWD=password";
	pool:=a.Pconnect("PoolSize=100")
	
	//SetConnMaxLifetime will atake the value in SECONDS
	db:=pool.Open(con,"SetConnMaxLifetime=30")
    st, err:=db.Prepare("Insert into SAMPLE values('hi','hi','hi','hi')")
    if err != nil{
        fmt.Println(err)
    }
	st.Query()
	
	//Here the the time out is default.
	db1:=pool.Open(con)
    st1, err:=db1.Prepare("Insert into SAMPLE values('hi1','hi1','hi1','hi1')")
    if err != nil{
        fmt.Println(err)
    }
	st1.Query()
	
	db1.Close()
	db.Close()
	pool.Release()
	fmt.println("success")
}
To run the sample:- go run example3.go
```
For Running the Tests:
======================

1) Put your connection string in the main.go file in testdata folder

2) Now run go test command (use go test -v command for details) 


