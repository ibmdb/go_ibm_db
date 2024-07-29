package main

import (
	"database/sql"
	"fmt"
	"time"
	"context"
	"strings"
	"os"
	"encoding/json"

	a "github.com/ibmdb/go_ibm_db"
)

var ctx = context.Background()

var host string
var port string
var database string
var uid string
var pwd string
var connStr string

//Read Config variable from json file
type Config  struct {
        Host string `json:"HOSTNAME"`
        Port string `json:"PORT"`
	Database string `json:"DATABASE"`
	Uid string `json:"UID"`
	Pwd string `json:"PWD"`
}

func LoadConfiguration(filename string) (Config, error) {
	fmt.Println("----LoadConfiguration() --")
        var config Config
        configFile, err := os.Open(filename)
        defer configFile.Close()
        if err != nil {
                return config, err
        }
        jsonParser := json.NewDecoder(configFile)
        err = jsonParser.Decode(&config)
        return config, err
}

//Get connection information from config.json
func GetConnectionInfoFromConfigFile() {
	fmt.Println("--GetConnectionInfoFromConfigFile()--")
       config, _:= LoadConfiguration("config.json")
       host = config.Host
       port = config.Port
       database = config.Database
       uid = config.Uid
       pwd =  config.Pwd
}


//Get connection information from environment variables 
func UpdateConnectionVariables() {
	var databaseFound bool
	var hostFound bool
	var portFound bool
	var uidFound bool
	var pwdFound bool
        fmt.Println("---UpdateConnectionVariables()--")
        config, _:= LoadConfiguration("./config.json")

	database, databaseFound = os.LookupEnv("DB2_DATABASE")
        if !databaseFound{
		database = config.Database
		if len(database) == 0 {
		    fmt.Println("Warning: Environment variable DB2_DATABASE is not set.")
	        }
	}else {
		fmt.Println("==Database: ", database)
		}
fmt.Println("==Database22: ", database)
	
	host, hostFound = os.LookupEnv("DB2_HOSTNAME")
        if !hostFound{
		host = config.Host
		if len(host)==0 {
		    fmt.Println("Warning: Environment variable DB2_HOSTNAME is not set.")
	        }
	}

	port, portFound = os.LookupEnv("DB2_PORT")
        if !portFound{
		port = config.Port
		if len(port)==0 {
		    fmt.Println("Warning: Environment variable DB2_PORT not set.")
	        }
	}

        uid, uidFound = os.LookupEnv("DB2_USER")
        if !uidFound{
		uid = config.Uid
		if len(uid)==0 {
		    fmt.Println("Warning: Environment variable DB2_USER is not set.")
	        }
	}

        pwd, pwdFound = os.LookupEnv("DB2_PASSWD")
        if !pwdFound{
               pwd = config.Pwd
               fmt.Println("Warning: Environment variable DB2_PASSWD is not set.")
               fmt.Println("Please set it before running test file and avoid")
               fmt.Println("hardcoded password in config.json file.")
	}
fmt.Println("database= " + database + " host= " + host)
}

//Createconnection will return the db instance
func Createconnection() (db *sql.DB) {
	fmt.Println("--Createconnection()--")
        UpdateConnectionVariables()
        //connStr = "PROTOCOL=tcpip;HOSTNAME=" + host + ";PORT=" + port + ";DATABASE=" + database + ";UID=" + uid + ";PWD=" + pwd
	connStr = "PROTOCOL=tcpip;HOSTNAME=waldevdbclnxtst06.dev.rocketsoftware.com;PORT=60000;DATABASE=sample;UID=zurbie;PWD=A2m8test"
	//connStr = "PROTOCOL=tcpip;HOSTNAME=" + host + ";PORT=" + port + ";DATABASE=" + database + ";UID=" + uid + ";PWD=" + pwd +";Security=ssl"
	fmt.Println("connStr: ", connStr)
	db, _ = sql.Open("go_ibm_db", connStr)
	return db
}

//Createtable will create the tables
func Createtable() error {
	db := Createconnection()
	defer db.Close()
	db.Exec("DROP table rocket")
	db.Exec("DROP table rocket1")
	_, err1 := db.Exec("create table rocket(a int)")
	if err1 != nil {
		return err1
	}

	_, err2 := db.Exec("create table rocket1(a int)")
	if err2 != nil {
		return err2
	}

	return nil
}

//Createtable will create the tables
func Createtable_ExecContext() error {
	db := Createconnection()
	defer db.Close()
	db.ExecContext(ctx, "DROP table rocket2")
	_, err := db.ExecContext(ctx, "create table rocket2(a int)")
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, "drop table rocket2")
	if err != nil {
		return err
	}
	return nil
}

//Insert will insert data in to the table
func Insert() error {
	db := Createconnection()
	defer db.Close()
	_, err := db.Exec("insert into rocket values(1)")
	if err != nil {
		return err
	}
	return nil
}

//Drop will drop the table
func Drop() error {
	db := Createconnection()
	defer db.Close()
	_, err := db.Exec("drop table rocket1")
	if err != nil {
		return err
	}
	return nil
}

//Prepare will prepare the statement
func Prepare() error {
	db := Createconnection()
	defer db.Close()
	_, err := db.Prepare("select * from rocket")
	if err != nil {
		return err
	}
	return nil
}

//Query will execute the prepared statement
func Query() error {
	db := Createconnection()
	defer db.Close()
	st, _ := db.Prepare("select * from rocket")
	_, err := st.Query()
	if err != nil {
		return err
	}
	return nil
}

//Scan will Scan the data in the rows
func Scan() error {
	db := Createconnection()
	defer db.Close()
	st, _ := db.Prepare("select * from rocket")
	rows, err := st.Query()
	for rows.Next() {
		var a string
		err = rows.Scan(&a)
		if err != nil {
			return err
		}
	}
	return nil
}

//Next will fetch the data from the result set
func Next() error {
	db := Createconnection()
	defer db.Close()
	st, _ := db.Prepare("select * from rocket")
	rows, err := st.Query()
	for rows.Next() {
		var a string
		err = rows.Scan(&a)
		if err != nil {
			return err
		}
	}
	return nil
}

//Columns will return the names of the cols
func Columns() error {
	db := Createconnection()
	defer db.Close()
	st, _ := db.Prepare("select * from rocket")
	rows, _ := st.Query()
	_, err := rows.Columns()
	if err != nil {
		return err
	}
	for rows.Next() {
		var a string
		_ = rows.Scan(&a)
	}
	return nil
}

//Queryrow will return the frirst row it matches
func Queryrow() error {
	a := 1
	var uname int
	db := Createconnection()
	defer db.Close()
	st, err := db.Prepare("select a from rocket where a=?")
	if err != nil {
		return err
	}
	err = st.QueryRow(a).Scan(&uname)
	if err != nil {
		return err
	}
	return nil
}

//Begin will start a transaction
func Begin() error {
	db := Createconnection()
	defer db.Close()
	_, err := db.Begin()
	if err != nil {
		return err
	}
	return nil
}

//Commit will commit the uncommited transactions
func Commit() error {
	db := Createconnection()
	defer db.Close()
	bg, err := db.Begin()
	db.Exec("DROP table u")
	_, err = bg.Exec("create table u(id int)")
	err = bg.Commit()
	if err != nil {
		return err
	}
	return nil
}

//Close will close the active connection
func Close() error {
	db := Createconnection()
	defer db.Close()
	err := db.Close()
	if err != nil {
		return err
	}
	return nil
}

//PoolOpen creates a pool and makes a connection.
func PoolOpen() int {
	pool := a.Pconnect("PoolSize=50")
	db := pool.Open(connStr)
	if db == nil {
		return 0
	}
	return 1
}

// A large integer is binary integer with a precision of 31 bits. The range is -2147483648 to +2147483647.
func IntegerArray() error {
        var tableOne string= "goarr"
        var errStr string

	db := Createconnection()
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 integer)")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []int{-2147483648, -32769, 0, 32768,  2147483647}
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []integer")
                return err
        }

        substring := "SQLSTATE=22003"
        c :=  []int{6}
        d :=  []int{-2147483649}
         st, err = db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(c, d)
        if err != nil {
                errStr = fmt.Sprintf("%s", err)
                if !strings.Contains(errStr, substring) {
                        return err
                }
        }


        e :=  []int{7}
        f :=  []int{2147483648}
         st, err = db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(e, f)
        if err != nil {
                errStr = fmt.Sprintf("%s", err)
                if !strings.Contains(errStr, substring) {
                        return err
                }
        }


        rows, err2 := db.Query("SELECT * from " + tableOne)
        if err2 != nil {
                return err
        }

        defer rows.Close()
        for rows.Next() {
              var c1, c2  string
              err = rows.Scan(&c1, &c2)
              if err != nil {
                      return err
              }

              //fmt.Printf("%v  %v \n", c1, c2)
        }

        db.Query("DROP table " + tableOne)

        return nil
}

//FloatArray function performs inserting float32,float64 datatypes.
func FloatArray() error {
	db := Createconnection()
	defer db.Close()
	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(var1 double)")
	if err != nil {
		return err
	}
	a := []float32{1.232, 2.34245, 3}
	b := []float64{3.43214321, 4.3243214645763235, 0}
	st, err := db.Prepare("Insert into arr values(?)")
	defer st.Close()
	if err != nil {
		return err
	}
	_, err = st.Query(a)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []float32")
		return err
	}
	_, err = st.Query(b)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []float64")
		return err
	}
	return nil
}


//CreateDB create database
func CreateDB() bool {
	res, err := a.CreateDb("Goo", connStr)
	if err != nil {
		return false
	}
	return res
}

//DropDB will drop database
func DropDB() bool {
	res, err := a.DropDb("Goo", connStr)
	if err != nil {
		return false
	}
	return res
}
//Execute Query for Connection pool
func ExecQuery(st *sql.Stmt) error {
        res, err := st.Query()
        if err != nil {
                return err
        }
        defer res.Close()
        for res.Next() {
                    var a string
                    err = res.Scan(&a)
                    if err != nil {
                            return err
                    }
        }
        return nil
}
//Connection pool
func ConnectionPool() int {
        var flag int
        flag = 0
        pool := a.Pconnect("PoolSize=5")

        ret := pool.Init(5, connStr)
        if ret != true {
		return 0
        }

        for i:=0; i<20; i++ {
                db := pool.Open(connStr, "SetConnMaxLifetime=10")
                if db != nil {
                        st, err := db.Prepare("select * from rocket")
                        if err != nil {
                                return 0
                        } else {
                                go func() {
                                        err := ExecQuery(st)
                                        if  err != nil && flag == 0{
                                             flag = 1
                                        }
                                        db.Close()
                                }()
                        }
                }
        }
        time.Sleep(10*time.Second)
        pool.Release()

        if flag == 1 {
                return 0
        }
        return 1
}


//Connection pool with Timeout
func ConnectionPoolWithTimeout() int {
        var flag int
        flag = 0
        pool := a.Pconnect("PoolSize=3")

	pool.SetConnMaxLifetime(10)
        ret := pool.Init(3, connStr)
        if ret != true {
		return 0
        }

        for i:=0; i<20; i++ {
                db := pool.Open(connStr, "SetConnMaxLifetime=10")
                if db != nil {
                        st, err := db.Prepare("select * from rocket")
                        if err != nil {
                                return 0
                        } else {
                                go func() {
                                        err := ExecQuery(st)
                                        if  err != nil && flag == 0{
                                             flag = 1
                                        }
                                        db.Close()
                                }()
                        }
                }
        }
        time.Sleep(30*time.Second)
	pool.SetConnMaxLifetime(0)

        pool.Release()

        if flag == 1 {
                return 0
        }
        return 1
}


//InserttArray function performs inserting int,float, boolean, character and string datatypes.
func InsertArray() error {
	db := Createconnection()
        defer db.Close()
        db.Exec("Drop table arr")

        _, err := db.Exec("create table arr(c1 int, c2 float, c3 boolean, c4 character, c5 varchar(20))")
        if err != nil {
                return err
        }
        a := []int{2, 3}
        b := []float32{20.45, 32.89}
        c := []bool{true, false}
        d := []string{"A", "B"}
        e := []string{"Hello!", "World"}

        st, err := db.Prepare("Insert into arr values(?, ?, ?, ?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }

        _, err = st.Query(a, b, c, d, e)

        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []int")
                return err
	}

	return nil
}




func main() {
	result := Createconnection()
	if result != nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-0")
	}

	result1 := Createtable()
	if result1 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-1")
	}

	result2 := Insert()
	if result2 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-2")
	}

	result3 := Drop()
	if result3 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-3")
	}

	result4 := Prepare()
	if result4 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-4")
	}

	result5 := Query()
	if result5 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-5")
	}

	result6 := Scan()
	if result6 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-6")
	}

	result7 := Next()
	if result7 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-7")
	}

	result8 := Columns()
	if result8 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-8")
	}

	result9 := Queryrow()
	if result9 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-9")
	}

	result10 := Begin()
	if result10 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-10")
	}

	result11 := Commit()
	if result11 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-11")
	}

	result12 := Close()
	if result12 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-12")
	}

	result13 := PoolOpen()
	if result13 == 1 {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-13")
	}
	result27 := CreateDB()
	if result27 == true {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-27")
	}

	result28 := DropDB()
	if result28 == true {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-28")
	}
	result32 := Createtable_ExecContext()
	if result32 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-32")
	}

	result33 := ConnectionPool()
        if result33 == 1 {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail-33")
        }

	result34 := ConnectionPoolWithTimeout()
        if result34 == 1 {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail-34")
        }
	result36 := InsertArray()
	if result36 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-36")
	}
	result39 := BadConnectionString()
        if result39 == 1 {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail-39")
        }
	result43 := ConnectionInvalidUserPassword()
        if result43 == 1  {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail-43")
        }

	result44 := ConnectionInvalidUserID()
        if result44 == 1  {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail-44")
        }

	result45 := ConnectionInvalidPortNumber()
        if result45 == 1  {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail-45")
        }

	result46 := ConnectionInvalidDatabaseName()
        if result46 == 1  {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail-46")
        }
	result54 := IntegerArray()
	if result54 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail-54")
	}
}

func BadConnectionString() int {
        var errStr string
	UpdateConnectionVariables()
        badConnStr := "HOSTNAME=hostname1;PORT1234=;DATABASE=sample;UID=uid;PWD=pwd"
        db, _ := sql.Open("go_ibm_db", badConnStr )
        _, err := db.Prepare("select * from arr")
        if err != nil {
                errStr = fmt.Sprintf("%s", err)
        }

        substring1 := "SQLSTATE=08001"
        substring2 := "SQLSTATE=08004"
        if strings.Contains(errStr, substring1) || strings.Contains(errStr, substring2) {
                return 1
        }  else {
                return 0
        }

        return 1
}

// Creating a table.
func QueryCreateTable(db *sql.DB) error {
        _, err := db.Query("DROP table VMSAMPLE")
        if err != nil {
               _, err := db.Query("CREATE table VMSAMPLE(ID varchar(20),NAME varchar(20),LOCATION varchar(20),POSITION varchar(20))")
               if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                        return err
                }
        } else {
              _, err := db.Query("CREATE table VMSAMPLE(ID varchar(20),NAME varchar(20),LOCATION varchar(20),POSITION varchar(20))")
               if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                        return err
                }
       }
       fmt.Println("TABLE CREATED Successfully")
       return nil
}

// Drop a table.
func QueryDropTable(db *sql.DB) error {
        _, err := db.Query("DROP table VMSAMPLE")
               if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                        return err
                }
        fmt.Println("TABLE DROP Successfully")
        return nil
}
func QueryInsertRow(db *sql.DB) error {
      _, err := db.Query("INSERT into VMSAMPLE(ID,NAME,LOCATION,POSITION) values('3242','Vikas','Blr','Developer')")
               if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                        return err
                }
        return nil
}

// This api selects the data from the table and prints it.
func QueryDisplayTable(db *sql.DB) error {
        rows, err := db.Query("SELECT * from VMSAMPLE")
        if err != nil {
                return err
        }

        defer rows.Close()
        for rows.Next() {
              var t, x, m, n string
              err = rows.Scan(&t, &x, &m, &n)
              if err != nil {
                       return err
              }
              fmt.Printf("%v  %v   %v         %v\n", t, x, m, n)
        }
        return nil
}


func ConnectionInvalidUserPassword() int {
        var errStr string
	UpdateConnectionVariables()
        //badConnStr := "HOSTNAME=hostname1;PORT1234=;DATABASE=sample;UID=uid;PWD=pwd"
	badConnStr := "PROTOCOL=tcpip;HOSTNAME=" + host + ";PORT=" + port + ";DATABASE=" + database + ";UID=" + uid + ";PWD=abcd"

        db, _ := sql.Open("go_ibm_db", badConnStr )
        _, err := db.Prepare("select * from arr")
        if err != nil {
                errStr = fmt.Sprintf("%s", err)
        }

        substring1 := "SQLSTATE=08001"
        substring2 := "SQLSTATE=08004"
        if strings.Contains(errStr, substring1) || strings.Contains(errStr, substring2) {
                return 1
        }  else {
                return 0
        }

        return 1
}

func ConnectionInvalidUserID() int {
        var errStr string
	UpdateConnectionVariables()
        badConnStr := "PROTOCOL=tcpip;HOSTNAME=" + host + ";PORT=" + port + ";DATABASE=" + database + ";UID=uid" + ";PWD=" + pwd
        db, _ := sql.Open("go_ibm_db", badConnStr )
        _, err := db.Prepare("select * from arr")
        if err != nil {
                errStr = fmt.Sprintf("%s", err)
        }

        substring1 := "SQLSTATE=08001"
        substring2 := "SQLSTATE=08004"
        if strings.Contains(errStr, substring1) || strings.Contains(errStr, substring2) {
                return 1
        }  else {
                return 0
        }

        return 1
}

func ConnectionInvalidPortNumber() int {
        var errStr string
	UpdateConnectionVariables()
        badConnStr := "PROTOCOL=tcpip;HOSTNAME=" + host + ";PORT=0" + ";DATABASE=" + database + ";UID=" + uid + ";PWD=" + pwd
        db, _ := sql.Open("go_ibm_db", badConnStr )
        _, err := db.Prepare("select * from arr")
        if err != nil {
                errStr = fmt.Sprintf("%s", err)
        }

        substring1 := "SQLSTATE=08001"
        substring2 := "SQLSTATE=08004"
        if strings.Contains(errStr, substring1) || strings.Contains(errStr, substring2) {
                return 1
        }  else {
                return 0
        }

        return 1
}

func ConnectionInvalidDatabaseName() int {
        var errStr string
	UpdateConnectionVariables()
        badConnStr := "PROTOCOL=tcpip;HOSTNAME=" + host + ";PORT=" + port + ";DATABASE=database"  + ";UID=" + uid + ";PWD=" + pwd
        db, _ := sql.Open("go_ibm_db", badConnStr )
        _, err := db.Prepare("select * from arr")
        if err != nil {
                errStr = fmt.Sprintf("%s", err)
        }

        substring1 := "SQLSTATE=08001"
        substring2 := "SQLSTATE=08004"
        if strings.Contains(errStr, substring1) || strings.Contains(errStr, substring2) {
                return 1
        }  else {
                return 0
        }

        return 1
}

