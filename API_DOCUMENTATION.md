# go-ibm_db API Documentation

## Database APIs

**Database APIs**

1.	[.Open(drivername,ConnectionString)](#OpenApi)
2.	[.Prepare(sqlquery)](#PrepareApi)
3.	[.Query(sqlquery)](#QueryApi)
4.	[.Exec(sqlquery)](#ExecApi)
5.	[.Begin()](#BeginApi)
6.	[.Close()](#CloseApi)
7.	[.Commit()](#CommitApi)
8.	[.Rollback()](#RollbackApi)
9.	[.QueryRow(sqlquery)](#QueryRowApi)
10.	[.Columns()](#ColumnsApi)
11.	[.Next()](#NextApi)
12.	[.Scan(options)](#ScanApi)

### <a name="OpenApi"></a> 1) .Open(drivername,ConnectionString)

open a connection to database
* **connectionString** - The connection string for your database.
* For distributed platforms, the connection string is typically defined   as: `DATABASE=dbname;HOSTNAME=hostname;PORT=port;PROTOCOL=TCPIP;UID=username;PWD=passwd`

```go
var connStr = flag.String("conn", "HOSTNAME=hostname;PORT=port;DATABASE=dbname;UID=uid;PWD=Pass", "connection string")
func dboper() error {
    fmt.Println("connecting to driver");
    db, err := sql.Open("drivername",*connStr);
    if err != nil {
        return err;
    }

    defer db.Close()
}
```
### <a name="PrepareApi"></a> 2) .Prepare(sqlquery)

Prepare a statement for execution
* **sql** - SQL string to prepare

Returns a ‘statement’ object

```go
func oper() error {
    fmt.Println("connecting to go-ibm_db");
    db, err:=sql.Open("go-imb_db",*connStr);
    if err != nil {
        return err;
    }
    
    defer db.Close()

    st, err := db.Prepare("select * from ak")
    if err !=nil {
        return err
    }

    rows,err := st.Query()
    if err != nil {
        return err
    }

    defer rows.Close()
}
```

### <a name="QueryApi"></a> 3) .Query(sqlquery)

Issue a SQL query to the database

If the query is executed then it will return the rows or it will return error

```go

func oper() error {
    fmt.Println("connecting to go-ibm_db");
    db, err := sql.Open("go-ibm_db",*connStr);
    if err != nil {
        return err;
    }

    defer db.Close()
    
    rows,err := db.Query(“select * from ak”)
    if err != nil {
        return err
    }
    
    defer rows.Close()
}
```


### <a name="ExecApi"></a> 4) .Exec(sqlquery)

Execute a prepared statement.

Only DML commands are performed. No data is returned back.

```go

func oper() error {
    fmt.Println("connecting to go-ibm_db");
    db, err := sql.Open("go-ibm_db",*connStr);
    if err != nil {
        return err;
    }

    defer db.Close()

    _,err = db.Exec("create table ghh(a int, b float, c double,  d char, e varchar(30))")
    if err != nil {
        return err
    }
}
```

### <a name="BeginApi"></a> 5) .Begin()

Begin a transaction.

```go


func oper() error {
    fmt.Println("connecting to go-ibm_db");
    db, err := sql.Open("go-ibm_db",*connStr);
    if err != nil {
        return err;
    }

    defer db.Close()
    
    bg, err := db.Begin()
    if err !=nil {
        return err
    }

    return nil
}
```



### <a name="CloseApi"></a> 6) .Close()

Close the currently opened database.

```go

func dboper() error {
    fmt.Println("connecting to go-ibm_db");
    db, err := sql.Open("go-ibm_db",*connStr);
    if err != nil {
        return err;
    }
    
    defer db.Close()
}
```



### <a name="CommitApi"></a> 7) .Commit()

Commit a transaction.

```go

func oper() error{
    fmt.Println("connecting to go-ibm_db");
    db, err := sql.Open("go-ibm_db",*connStr);
    if err != nil{
        return err;
    }
    

    defer db.Close()

    bg, err := db.Begin()
    if err != nil {
        return err
    }
    
    _,err=bg.Exec("create table ghh(a int,b float,c double,d char,e varchar(30))")
    if err != nil {
        return err
    }
    
    err = bg.Commit()
    if err != nil {
        return err
    }

    return nil
}
```




### <a name="RollbackApi"></a> 8) .Rollback()

Rollback a transaction.

```go


func oper() error{
    fmt.Println("connecting to go-ibm_db");
    db, err := sql.Open("go-ibm_db",*connStr);
    if err != nil {
        return err;
    }

    defer db.Close()
    bg, err := db.Begin()
    if err !=nil {
        return err
    }

    _,err=bg.Exec("create table ghh(a int,b float,c double,d char,e varchar(30))")
    if err != nil {
        return err
    }

    err = bg.Rollback()
    if err != nil {
        return err
    }
    
    return nil
}
```

### <a name="QueryRowApi"></a> 9) .QueryRow(sqlquery)

QueryRow executes a query that is expected to return at most one row.
If there are more rows then it will scan first and discards the rest.
 
```go

func oper()  error {
    id := 123
    var username string
    err := db.QueryRow("SELECT name FROM ak WHERE id=?", id).Scan(&username)
    if err != nil {
        return err
    }

    fmt.Printf("Username is %s\n", username)
    return nil
}

```

### <a name="ColumnsApi"></a> 10) .Columns()

Returns the column names.

Returns error if the rows are closed.

```go

func oper() error {
    fmt.Println("connecting to databse");
    db, err := sql.Open("go-ibm_db",*connStr);
    if err != nil {
        return err;
    }

    defer db.Close()

    st, err := db.Prepare("select * from ak")
    if err !=nil {
    return err
    }
    
    rows,err :=st.Query()
    if err != nil {
        return err
    }
    
    defer rows.Close()
    name11 := make([]string,1)
    name11, err = rows.Columns()
    fmt.Printf("%v",name11);
    return nil
}
```

### <a name="NextApi"></a> 11) .Next()

Prepares the next result row for reading with the scan api.

```go

func oper() error {
    fmt.Println("connecting to database");
    db, err:=sql.Open("go-ibm_db",*connStr);
    if err != nil {
        return err;
    }

    defer db.Close()
    rows,err := db.Query()
    if err != nil {
        return err
    }
    
    defer rows.Close()
    for rows.Next() {
        var t string
        var x string
        err = rows.Scan(&t, &x)
        if err != nil {
            return err
        }

        fmt.Printf("%v %v\n",t,x)
    }
    
    return nil
}
```

### <a name="ScanApi"></a> 12) .Scan(options)

copies the columns in the current row into the values pointed.

```go

func oper() error {
    fmt.Println("connecting to database");
    db, err := sql.Open("go-ibm_db", *connStr);
    if err != nil {
        return err;
    }

    defer db.Close()

    rows,err := db.Query()
    if err != nil {
        return err
    }

    defer rows.Close()
    for rows.Next() {
        var t string
        var x string
        err = rows.Scan(&t, &x)
        if err != nil {
            return err
        }

        fmt.Printf("%v %v\n",t,x)
    }
    return nil
}
```

