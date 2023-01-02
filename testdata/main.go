package main

import (
	"database/sql"
	"fmt"
	"time"
	"context"
	"strings"

	a "github.com/ibmdb/go_ibm_db"
)

var ctx = context.Background()

var host = "<HOST>"
var port = "<PORT>"
var database = "<DATABASE>"
var uid = "<UID>"
var pwd = "<PWD>"

var connStr = "PROTOCOL=tcpip;HOSTNAME=" + host + ";PORT=" + port + ";DATABASE=" + database + ";UID=" + uid + ";PWD=" + pwd


//Createconnection will return the db instance
func Createconnection() (db *sql.DB) {
	db, _ = sql.Open("go_ibm_db", connStr)
	return db
}

//Createtable will create the tables
func Createtable() error {
	db, err := sql.Open("go_ibm_db", connStr)
	db.Exec("DROP table rocket")
	_, err = db.Exec("create table rocket(a int)")
	_, err = db.Exec("create table rocket1(a int)")
	if err != nil {
		return err
	}
	return nil
}

//Createtable will create the tables
func Createtable_ExecContext() error {
	db, err := sql.Open("go_ibm_db", connStr)
	db.ExecContext(ctx, "DROP table rocket2")
	_, err = db.ExecContext(ctx, "create table rocket2(a int)")
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
	db, err := sql.Open("go_ibm_db", connStr)
	_, err = db.Exec("insert into rocket values(1)")
	if err != nil {
		return err
	}
	return nil
}

//Drop will drop the table
func Drop() error {
	db, err := sql.Open("go_ibm_db", connStr)
	_, err = db.Exec("drop table rocket1")
	if err != nil {
		return err
	}
	return nil
}

//Prepare will prepare the statement
func Prepare() error {
	db, _ := sql.Open("go_ibm_db", connStr)
	_, err := db.Prepare("select * from rocket")
	if err != nil {
		return err
	}
	return nil
}

//PrepareContext will prepare the statement
func PrepareContext() error {
	db, _ := sql.Open("go_ibm_db", connStr)
	_, err := db.PrepareContext(ctx, "select * from rocket")
	if err != nil {
		return err
	}
	return nil
}

//Query will execute the prepared statement
func Query() error {
	db, _ := sql.Open("go_ibm_db", connStr)
	st, _ := db.Prepare("select * from rocket")
	_, err := st.Query()
	if err != nil {
		return err
	}
	return nil
}

//QueryContext will execute the prepared statement
func QueryContext() error {
	db, _ := sql.Open("go_ibm_db", connStr)
	st, _ := db.PrepareContext(ctx, "select * from rocket")
	_, err := st.QueryContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
//ExecDirect will execute the query without prepare
func ExecDirect() error {
	db, _ := sql.Open("go_ibm_db", connStr)
	_, err := db.Query("select * from rocket")
	if err != nil {
		return err
	}
	return nil
}

//Scan will Scan the data in the rows
func Scan() error {
	db, _ := sql.Open("go_ibm_db", connStr)
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
	db, _ := sql.Open("go_ibm_db", connStr)
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
	db, _ := sql.Open("go_ibm_db", connStr)
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
	db, err := sql.Open("go_ibm_db", connStr)
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
	db, err := sql.Open("go_ibm_db", connStr)
	_, err = db.Begin()
	if err != nil {
		return err
	}
	return nil
}

//Commit will commit the uncommited transactions
func Commit() error {
	db, err := sql.Open("go_ibm_db", connStr )
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
	db, _ := sql.Open("go_ibm_db", connStr)
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

//StoredProcedure function tests OUT Parameter by calling get_dbsize_info.
func StoredProcedure() error {
	var (
		snapTime   time.Time
		dbsize     int64
		dbcapacity int64
	)
	db, _ := sql.Open("go_ibm_db", connStr)
	_, err := db.Exec("call sysproc.get_dbsize_info(?, ?, ?,0)", sql.Out{Dest: &snapTime}, sql.Out{Dest: &dbsize}, sql.Out{Dest: &dbcapacity})
	if err != nil {
		return err
	}
	return nil
}

//StoredProcedureInOut function tests OUT Parameter by calling get_dbsize_info.
func StoredProcedureInOut() error {
	db, _ := sql.Open("go_ibm_db", connStr)
	in1 := 10
	inout1 := 2
	var out1, out2 int
	st, err := db.Prepare("create or replace procedure sp2(in var1 integer, inout var2 integer, out var3 integer, out var4 integer) LANGUAGE SQL BEGIN   SET var2 = var1 + var2; SET var3 = var1 - var2; SET var4 = var1 * var2; END")
	if err != nil {
		return err
	}
	st.Query()
	_, err = db.Exec("CALL sp2(?,?,?,?)", in1, sql.Out{Dest: &inout1, In: true}, sql.Out{Dest: &out1}, sql.Out{Dest: &out2})
	if err != nil {
		return err
	}
	if inout1 != 12 || out1 != -2 || out2 != 120 {
		return fmt.Errorf("Wrong data retrieved")
	}
	return nil
}

//IntArray function performs inserting int,int8,int16,int32,int64 datatypes.
func IntArray() error {
	db, _ := sql.Open("go_ibm_db", connStr)
	defer db.Close()
	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(var1 int)")
	if err != nil {
		return err
	}
	a := []int{2, 3}
	b := []int8{2, 3}
	c := []int16{2, 3}
	d := []int32{2, 3}
	e := []int64{2, 3}
	st, err := db.Prepare("Insert into arr values(?)")
	defer st.Close()
	if err != nil {
		return err
	}
	_, err = st.Query(a)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []int")
		return err
	}
	_, err = st.Query(b)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []int8")
		return err
	}
	_, err = st.Query(c)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []int16")
		return err
	}
	_, err = st.Query(d)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []int32")
		return err
	}
	_, err = st.Query(e)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []int64")
		return err
	}
	return nil
}

//A small integer is a binary integer with a precision of 15 bits. The range of small integers is -32768 to +32767.
func SmallintArray() error {
        var tableOne string= "goarr"
        var errStr string

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 smallint)")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []int{-32768, -123, 0, 100,  32767}
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []smallint")
                return err
        }

        substring := "SQLSTATE=22003"
        c :=  []int{6, 7}
        d :=  []int{-32769, 32768}
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


// A big integer is a binary integer with a precision of 63 bits. The range of big integers is -9223372036854775808 to +9223372036854775807.
func BigintArray() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 bigint)")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5, 6, 7}
        b :=  []int{-9223372036854775808, -2147483648, -32769, 0, 32768,  2147483647, 9223372036854775807}

        st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []bigint")
                return err
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

// A large integer is binary integer with a precision of 31 bits. The range is -2147483648 to +2147483647.
func IntegerArray() error {
        var tableOne string= "goarr"
        var errStr string

        db, _ := sql.Open("go_ibm_db", connStr )
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

// A large integer is binary integer with a precision of 31 bits. The range is -2147483648 to +2147483647.
func Int2Array() error {
        var tableOne string= "goarr"
        var errStr string

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 int)")
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
                fmt.Println("Error while inserting []int2")
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

//numeric(p,s) The decimal data type is an exact numeric data type defined by its precision (total number of digits)
//and scale (number of digits to the right of the decimal point).
//p Defines the precision. Minimum 1; maximum is 39
//s Defines the scale. The scale of a decimal value cannot exceed its precision. Scale can be 0 (no digits to the right of the decimal point).

func NumericArray() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )

        db.Query("DROP table " + tableOne)
	defer db.Close()

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 numeric(5,2))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3}
        b :=  []float32 { -152.3, 56.08, 100.238567}
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []numeric")
                return err
        }

        var errStr string
        substring := "SQLSTATE=22003"
        c :=  []int{4}
        d :=  []float32{1234.98}
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


//Numeric(p) The decimal data type is an exact numeric data type defined by its precision (total number of digits)
//p Defines the precision. It has at a total of p (<=32) significat digits

func Numeric2Array() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 numeric(31))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3}
        b :=  []float32 { -10e30, 987654321.123456, 10e30}
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []Numeric2")
                return err
        }


        var errStr string
        substring := "SQLSTATE=22003"
        c :=  []int{4}
        d :=  []float32{10e31}
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


//StringArray function performs inserting string array.
func StringArray() error {
	db, _ := sql.Open("go_ibm_db", connStr)
	defer db.Close()
	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(var1 varchar(10),var2 varchar(20))")
	if err != nil {
		return err
	}
	a := []string{"value1", "value"}
	b := []string{"value", "value22"}
	st, err := db.Prepare("Insert into arr values(?,?)")
	defer st.Close()
	if err != nil {
		return err
	}
	_, err = st.Query(a, b)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []string")
		return err
	}
	return nil
}


//CLOB(n) Varying-length character strings with a maximum of n characters. n cannot exceed 2,147,483,647. The default length is 1M.
func ClobArray() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 CLOB(5))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []string  { "a", "ab", "abc", "abcd", "abcde" }
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []clob")
                return err
        }

        var errStr string
        substring := "SQLSTATE=22001"
        c :=  []int{6}
        d :=  []string{"abcdef"}
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

//A CLOB (character large object) value can be up to 2,147,483,647 characters long.
//A CLOB is used to store unicode character-based data, such as large documents in any character set.
//The length is given in number characters for both CLOB, unless one of the suffixes K, M, or G is given,
//relating to the multiples of 1024, 1024*1024, 1024*1024*1024 respectively.
//{CLOB |CHARACTER LARGE OBJECT} [ ( length [{K |M |G}] ) ]
func Clob2Array() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 clob(64 K))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []string  { "Hello World!!", "!@#$%^&*()", "1234567890", "How are you?", "Happy Birthday!" }
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []clob2")
                return err
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

//VARCHAR(n):Varying-length character strings with a maximum length of n bytes.
//n must be greater than 0 and less than a number that depends on the page size of the table space. The maximum length is 32704.
func VarcharArray() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 varchar(5))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []string  { "a", "ab", "abc", "abcd", "abcde" }
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []varchar")
                return err
        }

        var errStr string
        substring := "SQLSTATE=22001"
        c :=  []int{6}
        d :=  []string{"abcdef"}
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

//GRAPHIC(n) Fixed-length graphic strings that contain n double-byte characters. n must be greater than 0 and less than 128. The default length is 1.
func GraphicArray() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 graphic(5))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []string  { "a", "ab", "abc", "abcd", "abcde" }
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []character")
                return err
        }

        var errStr string
        substring := "SQLSTATE=22001"
        c :=  []int{6}
        d :=  []string{"abcdef"}
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


//VARGRAPHIC(n) Varying-length graphic strings. The maximum length, n, must be greater than 0
//and less than a number that depends on the page size of the table space. The maximum length is 16352.
func VargraphicArray() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 vargraphic(20))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []string  { "Hello World!!", "!@#$%^&*()", "1234567890", "How are you?", "Happy Birthday!" }
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []vargraphic")
                return err
        }
        var errStr string
        substring := "SQLSTATE=22001"
        c :=  []int{6}
        //d :=  []string{"abcdefghijklmnopqurstuvwxyz"}
        d :=  []string{"123456789012345678901"}
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


//BoolArray function performs inserting bool array.
func BoolArray() error {
	db, _ := sql.Open("go_ibm_db", connStr)
	defer db.Close()
	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(var1 boolean,var2 boolean)")
	if err != nil {
		return err
	}
	a := []bool{false, false, false, false, false}
	b := []bool{true, true, true, true, true}
	st, err := db.Prepare("Insert into arr values(?,?)")
	defer st.Close()
	if err != nil {
		return err
	}
	_, err = st.Query(a, b)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []bool")
		return err
	}
	return nil
}

//FloatArray function performs inserting float32,float64 datatypes.
func FloatArray() error {
	db, _ := sql.Open("go_ibm_db", connStr)
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

func Float2Array() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 float)")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []float64 { -1.79769E308, 1.79769E308, 9876543210.123456789, 2.225E-307, -2.225E-307}
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []float")
                return err
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


//DOUBLE value ranges:
//Smallest DOUBLE value: -1.79769E+308
//Largest DOUBLE value: 1.79769E+308
//Smallest positive DOUBLE value: 2.225E-307
//Largest negative DOUBLE value: -2.225E-307
func DoubleArray() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 double)")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []float64 { -1.79769E308, 1.79769E308, 9876543210.123456789, 2.225E-307, -2.225E-307}
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []double")
                return err
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

//DOUBLE PRECISION value ranges:
//Smallest DOUBLE value: -1.79769E+308
//Largest DOUBLE value: 1.79769E+308
//Smallest positive DOUBLE value: 2.225E-307
//Largest negative DOUBLE value: -2.225E-307
func DoubleprecisionArray() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 double precision)")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4}
        b :=  []float64 { -1.79769E+308, 1.79769E+308, 2.225E-307,  -2.225E-307 }
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []double precision")
                return err
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


//REAL value ranges:
//Smallest REAL value: -3.402E+38
//Largest REAL value: 3.402E+38
//Smallest positive REAL value: 1.175E-37
//Largest negative REAL value: -1.175E-37
func RealArray() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 real)")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []float32 { -3.402e38, 987654321.123456, 3.402e38, 1.175e-37, -1.175e-37}
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []real")
                return err
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

//decimal(p,s) The decimal data type is an exact numeric data type defined by its precision (total number of digits)
//and scale (number of digits to the right of the decimal point).
//p Defines the precision. Minimum 1; maximum is 39
//s Defines the scale. The scale of a decimal value cannot exceed its precision. Scale can be 0 (no digits to the right of the decimal point).
func DecimalArray() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 decimal(5,2))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3}
        b :=  []float32 { -152.3, 56.08, 100.238567}
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []decimal")
                return err
        }

        var errStr string
        substring := "SQLSTATE=22003"
        c :=  []int{4}
        d :=  []float32{1234.98}
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

//decimal(p) The decimal data type is an exact numeric data type defined by its precision (total number of digits)
//p Defines the precision. It has at a total of p (<=32) significat digits

func Decimal2Array() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 numeric(31))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3}
        b :=  []float32 { -10e30, 987654321.123456, 10e30}
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []decimal")
                return err
        }

        var errStr string
        substring := "SQLSTATE=22003"
        c :=  []int{4}
        d :=  []float32{10e31}
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

//dec(p) The decimal data type is an exact numeric data type defined by its precision (total number of digits)
//p Defines the precision. It has at a total of p (<=32) significat digits

func Dec2Array() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 dec(31))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3}
        b :=  []float32 { -10e30, 987654321.123456, 10e30}
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []dec")
                return err
        }

        var errStr string
        substring := "SQLSTATE=22003"
        c :=  []int{4}
        d :=  []float32{10e31}
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









//CharArray function performs inserting float32,float64 datatypes.
func CharArray() error {
	db, _ := sql.Open("go_ibm_db", connStr)
	defer db.Close()
	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(var1 character, var2 character)")
	if err != nil {
		return err
	}
	a := []string{"a", "b", "c"}
	b := []string{"z", "y", "x"}
	st, err := db.Prepare("Insert into arr values(?,?)")
	defer st.Close()
	if err != nil {
		return err
	}
	_, err = st.Query(a, b)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []character")
		return err
	}
	return nil
}

//CHARACTER(n)Fixed-length character strings with a length of n bytes. n must be greater than 0 and not greater than 255. The default length is 1.
func CharacterArray() error {
        var tableOne string= "goarr"

        db, _ := sql.Open("go_ibm_db", connStr )
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 character(5))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []string  { "a", "ab", "abc", "abcd", "abcde" }
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []character")
                return err
        }

        var errStr string
        substring := "SQLSTATE=22001"
        c :=  []int{6}
        d :=  []string{"abcdef"}
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



//TimeStampArray function performs inserting float32,float64 datatypes.
func TimeStampArray() error {
	db, _ := sql.Open("go_ibm_db", connStr)
	defer db.Close()
	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(var1 timestamp, var2 time, var3 date)")
	if err != nil {
		return err
	}
	a := []time.Time{}
	for i := 0; i < 5; i++ {
		a = append(a, time.Now())
		time.Sleep(1 * time.Second)
	}
	st, err := db.Prepare("Insert into arr values(?,?,?)")
	defer st.Close()
	if err != nil {
		return err
	}
	_, err = st.Query(a, a, a)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []timestamp")
		return err
	}
	return nil
}

//NullValueCharacter function performs
func NullValueCharacter() error {
	var out1, out2 sql.NullString
	var out3 sql.NullInt64
	var out4 sql.NullBool
	var out5 sql.NullFloat64
	var out6 sql.NullTime
	db, _ := sql.Open("go_ibm_db", connStr)
	defer db.Close()
	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(var1 character, var2 varchar(30), var3 integer, var4 boolean, var5 double, var6 timestamp)")
	if err != nil {
		return err
	}
	c2 := "test"
	c3 := int64(10)
	c4 := true
	c5 := 1.234
	c6 := time.Now()
	st, err := db.Prepare("Insert into arr(var2,var3,var4,var5,var6) values(?,?,?,?,?)")
	defer st.Close()
	if err != nil {
		return err
	}
	_, err = st.Query(c2, c3, c4, c5, c6)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting NullValueCharacter")
		return err
	}
	st1, err := db.Prepare("select * from arr")
	defer st1.Close()
	if err != nil {
		return err
	}
	rows, err := st1.Query()
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rows.Scan(&out1, &out2, &out3, &out4, &out5, &out6)
		if err != nil {
			fmt.Println("Error while retrieving NullValueCharacter")
			return err
		}
		if out2.String != c2 {
			fmt.Println("Data is mismatched at NullValueCharacter - String")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out3.Int64 != c3 {
			fmt.Println("Data is mismatched at NullValueCharacter - Int64")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out4.Bool != c4 {
			fmt.Println("Data is mismatched at NullValueCharacter - Bool")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out5.Float64 != c5 {
			fmt.Println("Data is mismatched at NullValueCharacter - Float")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out6.Time.Format("2006-01-02 15:04:05 PM -07:00 Jan Mon MST") != c6.Format("2006-01-02 15:04:05 PM -07:00 Jan Mon MST") {
			fmt.Println("Data is mismatched at NullValueCharacter - Time")
			return fmt.Errorf("Wrong data retrieved")
		}
	}
	return nil
}

//NullValueString function performs
func NullValueString() error {
	var out1, out2 sql.NullString
	var out3 sql.NullInt64
	var out4 sql.NullBool
	var out5 sql.NullFloat64
	var out6 sql.NullTime
	db, _ := sql.Open("go_ibm_db", connStr)
	defer db.Close()
	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(var1 character, var2 varchar(30), var3 integer, var4 boolean, var5 double, var6 timestamp)")
	if err != nil {
		return err
	}
	c1 := "a"
	c3 := int64(10)
	c4 := true
	c5 := 1.234
	c6 := time.Now()
	st, err := db.Prepare("Insert into arr(var1,var3,var4,var5,var6) values(?,?,?,?,?)")
	defer st.Close()
	if err != nil {
		return err
	}
	_, err = st.Query(c1, c3, c4, c5, c6)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting NullValueString")
		return err
	}
	st1, err := db.Prepare("select * from arr")
	defer st1.Close()
	if err != nil {
		return err
	}
	rows, err := st1.Query()
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rows.Scan(&out1, &out2, &out3, &out4, &out5, &out6)
		if err != nil {
			fmt.Println("Error while retrieving NullValueString")
			return err
		}
		if out1.String != c1 {
			fmt.Println("Data is mismatched at NullValueString - Character")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out3.Int64 != c3 {
			fmt.Println("Data is mismatched at NullValueString - Int64")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out4.Bool != c4 {
			fmt.Println("Data is mismatched at NullValueString - Bool")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out5.Float64 != c5 {
			fmt.Println("Data is mismatched at NullValueString - Float")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out6.Time.Format("2006-01-02 15:04:05 PM -07:00 Jan Mon MST") != c6.Format("2006-01-02 15:04:05 PM -07:00 Jan Mon MST") {
			fmt.Println("Data is mismatched at NullValueString - Time")
			return fmt.Errorf("Wrong data retrieved")
		}
	}
	return nil
}

//NullValueInteger function performs
func NullValueInteger() error {
	var out1, out2 sql.NullString
	var out3 sql.NullInt64
	var out4 sql.NullBool
	var out5 sql.NullFloat64
	var out6 sql.NullTime
	db, _ := sql.Open("go_ibm_db", connStr)
	defer db.Close()
	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(var1 character, var2 varchar(30), var3 integer, var4 boolean, var5 double, var6 timestamp)")
	if err != nil {
		return err
	}
	c1 := "a"
	c2 := "test"
	c4 := true
	c5 := 1.234
	c6 := time.Now()
	st, err := db.Prepare("Insert into arr(var1,var2,var4,var5,var6) values(?,?,?,?,?)")
	defer st.Close()
	if err != nil {
		return err
	}
	_, err = st.Query(c1, c2, c4, c5, c6)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting NullValueInteger")
		return err
	}
	st1, err := db.Prepare("select * from arr")
	defer st1.Close()
	if err != nil {
		return err
	}
	rows, err := st1.Query()
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rows.Scan(&out1, &out2, &out3, &out4, &out5, &out6)
		if err != nil {
			fmt.Println("Error while retrieving NullValueInteger")
			return err
		}
		if out1.String != c1 {
			fmt.Println("Data is mismatched at NullValueInteger - character")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out2.String != c2 {
			fmt.Println("Data is mismatched at NullValueInteger - String")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out4.Bool != c4 {
			fmt.Println("Data is mismatched at NullValueInteger - Bool")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out5.Float64 != c5 {
			fmt.Println("Data is mismatched at NullValueInteger - Float")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out6.Time.Format("2006-01-02 15:04:05 PM -07:00 Jan Mon MST") != c6.Format("2006-01-02 15:04:05 PM -07:00 Jan Mon MST") {
			fmt.Println("Data is mismatched at NullValueInteger - Time")
			return fmt.Errorf("Wrong data retrieved")
		}
	}
	return nil
}

//NullValueBool function performs
func NullValueBool() error {
	var out1, out2 sql.NullString
	var out3 sql.NullInt64
	var out4 sql.NullBool
	var out5 sql.NullFloat64
	var out6 sql.NullTime
	db, _ := sql.Open("go_ibm_db", connStr)
	defer db.Close()
	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(var1 character, var2 varchar(30), var3 integer, var4 boolean, var5 double, var6 timestamp)")
	if err != nil {
		return err
	}
	c1 := "a"
	c2 := "test"
	c3 := int64(10)
	c5 := 1.234
	c6 := time.Now()
	st, err := db.Prepare("Insert into arr(var1,var2,var3,var5,var6) values(?,?,?,?,?)")
	defer st.Close()
	if err != nil {
		return err
	}
	_, err = st.Query(c1, c2, c3, c5, c6)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting NullValueInteger")
		return err
	}
	st1, err := db.Prepare("select * from arr")
	defer st1.Close()
	if err != nil {
		return err
	}
	rows, err := st1.Query()
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rows.Scan(&out1, &out2, &out3, &out4, &out5, &out6)
		if err != nil {
			fmt.Println("Error while retrieving NullValueBool")
			return err
		}
		if out1.String != c1 {
			fmt.Println("Data is mismatched at NullValueBool - Character")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out2.String != c2 {
			fmt.Println("Data is mismatched at NullValueBool - string")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out3.Int64 != c3 {
			fmt.Println("Data is mismatched at NullValueBool - Int64")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out5.Float64 != c5 {
			fmt.Println("Data is mismatched at NullValueBool - Float")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out6.Time.Format("2006-01-02 15:04:05 PM -07:00 Jan Mon MST") != c6.Format("2006-01-02 15:04:05 PM -07:00 Jan Mon MST") {
			fmt.Println("Data is mismatched at NullValueBool - Time")
			return fmt.Errorf("Wrong data retrieved")
		}
	}
	return nil
}

//NullValueFloat function performs
func NullValueFloat() error {
	var out1, out2 sql.NullString
	var out3 sql.NullInt64
	var out4 sql.NullBool
	var out5 sql.NullFloat64
	var out6 sql.NullTime
	db, _ := sql.Open("go_ibm_db", connStr)
	defer db.Close()
	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(var1 character, var2 varchar(30), var3 integer, var4 boolean, var5 double, var6 timestamp)")
	if err != nil {
		return err
	}
	c1 := "a"
	c2 := "test"
	c3 := int64(10)
	c4 := true
	c6 := time.Now()
	st, err := db.Prepare("Insert into arr(var1,var2,var3,var4,var6) values(?,?,?,?,?)")
	defer st.Close()
	if err != nil {
		return err
	}
	_, err = st.Query(c1, c2, c3, c4, c6)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting NullValueFloat")
		return err
	}
	st1, err := db.Prepare("select * from arr")
	defer st1.Close()
	if err != nil {
		return err
	}
	rows, err := st1.Query()
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rows.Scan(&out1, &out2, &out3, &out4, &out5, &out6)
		if err != nil {
			fmt.Println("Error while retrieving NullValueFloat")
			return err
		}
		if out1.String != c1 {
			fmt.Println("Data is mismatched at NullValueFloat - Character")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out2.String != c2 {
			fmt.Println("Data is mismatched at NullValueFloat - string")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out3.Int64 != c3 {
			fmt.Println("Data is mismatched at NullValueFloat - Int64")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out4.Bool != c4 {
			fmt.Println("Data is mismatched at NullValueFloat - Bool")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out6.Time.Format("2006-01-02 15:04:05 PM -07:00 Jan Mon MST") != c6.Format("2006-01-02 15:04:05 PM -07:00 Jan Mon MST") {
			fmt.Println("Data is mismatched at NullValueFloat - Time")
			return fmt.Errorf("Wrong data retrieved")
		}
	}
	return nil
}

//NullValueTime function performs
func NullValueTime() error {
	var out1, out2 sql.NullString
	var out3 sql.NullInt64
	var out4 sql.NullBool
	var out5 sql.NullFloat64
	var out6 sql.NullTime
	db, _ := sql.Open("go_ibm_db", connStr)
	defer db.Close()
	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(var1 character, var2 varchar(30), var3 integer, var4 boolean, var5 double, var6 timestamp)")
	if err != nil {
		return err
	}
	c1 := "a"
	c2 := "test"
	c3 := int64(10)
	c4 := true
	c5 := 1.234
	st, err := db.Prepare("Insert into arr(var1,var2,var3,var4,var5) values(?,?,?,?,?)")
	defer st.Close()
	if err != nil {
		return err
	}
	_, err = st.Query(c1, c2, c3, c4, c5)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting NullValueString")
		return err
	}
	st1, err := db.Prepare("select * from arr")
	defer st1.Close()
	if err != nil {
		return err
	}
	rows, err := st1.Query()
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rows.Scan(&out1, &out2, &out3, &out4, &out5, &out6)
		if err != nil {
			fmt.Println("Error while retrieving NullValueString")
			return err
		}
		if out1.String != c1 {
			fmt.Println("Data is mismatched at NullValueFloat - Character")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out2.String != c2 {
			fmt.Println("Data is mismatched at NullValueFloat - string")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out3.Int64 != c3 {
			fmt.Println("Data is mismatched at NullValueFloat - Int64")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out4.Bool != c4 {
			fmt.Println("Data is mismatched at NullValueFloat - Bool")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out5.Float64 != c5 {
			fmt.Println("Data is mismatched at NullValueFloat - Float64")
			return fmt.Errorf("Wrong data retrieved")
		}
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

func ChineseChar() error {
        db, _ := sql.Open("go_ibm_db", connStr)
        defer db.Close()
        db.Exec("Drop table arr")
        _, err := db.Exec("create table arr(ID bigint, var2 varchar(30))")
        if err != nil {
                return err
        }
        //st, err := db.Prepare("Insert into arr values('101','2019年2▒~V~RH')")
        st, err := db.Prepare("Insert into arr values('101',x'32303139E5B9B431E69C88E4')")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query()
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []string")
                return err
        }

	return nil
}


//InserttArray function performs inserting int,float, boolean, character and string datatypes.
func InsertArray() error {
        db, _ := sql.Open("go_ibm_db", connStr)
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
		fmt.Println("Fail")
	}

	result1 := Createtable()
	if result1 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result2 := Insert()
	if result2 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result3 := Drop()
	if result3 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result4 := Prepare()
	if result4 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result5 := Query()
	if result5 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result6 := Scan()
	if result6 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result7 := Next()
	if result7 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result8 := Columns()
	if result8 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result9 := Queryrow()
	if result9 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result10 := Begin()
	if result10 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result11 := Commit()
	if result11 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result12 := Close()
	if result12 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result13 := PoolOpen()
	if result13 == 1 {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result14 := StoredProcedure()
	if result14 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result15 := IntArray()
	if result15 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result16 := StringArray()
	if result16 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result17 := BoolArray()
	if result17 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result18 := CharArray()
	if result18 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result19 := TimeStampArray()
	if result19 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result20 := NullValueCharacter()
	if result20 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result21 := NullValueString()
	if result21 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result22 := NullValueInteger()
	if result22 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result23 := NullValueBool()
	if result23 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result24 := NullValueFloat()
	if result24 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}
	result25 := NullValueTime()
	if result25 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result26 := StoredProcedureInOut()
	if result26 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result27 := CreateDB()
	if result27 == true {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result28 := DropDB()
	if result28 == true {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result29 := ExecDirect()
	if result29 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result30 := PrepareContext()
	if result30 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result31 := QueryContext()
	if result31 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result32 := Createtable_ExecContext()
	if result32 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result33 := ConnectionPool()
        if result33 == 1 {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail")
        }

	result34 := ConnectionPoolWithTimeout()
        if result34 == 1 {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail")
        }

	result35 := ChineseChar()
	if result35 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result36 := InsertArray()
	if result36 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result37 := AllDataTypes()
	if result37 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result38 := MultipleQuery()
	if result38 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result39 := BadConnectionString()
        if result39 == 1 {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail")
        }

	result40 := CreateDropTable()
        if result40 == 1  {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail")
        }

	result41 := QueryCreateDropTable()
        if result41 == 1  {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail")
        }

	result42 := QueryInsertSelect()
        if result42 == 1  {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail")
        }

	result43 := ConnectionInvalidUserPassword()
        if result43 == 1  {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail")
        }

	result44 := ConnectionInvalidUserID()
        if result44 == 1  {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail")
        }

	result45 := ConnectionInvalidPortNumber()
        if result45 == 1  {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail")
        }

	result46 := ConnectionInvalidDatabaseName()
        if result46 == 1  {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail")
        }

	result47 := HugeQuery()
        if result47 == 1  {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail")
        }

	result48 := DecimalColumn()
        if result48 == 1  {
                fmt.Println("Pass")
        } else {
                fmt.Println("Fail")
        }

	result49 := RollbackTransaction()
	if result49 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result50 := Decfloat()
	if result50 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result51 := BigintArray()
	if result51 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result52 := CharacterArray()
	if result52 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result53 := SmallintArray()
	if result53 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result54 := IntegerArray()
	if result54 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result55 := Int2Array()
	if result55 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result56 := NumericArray()
	if result56 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result57 := Numeric2Array()
	if result57 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result58 := ClobArray()
	if result58 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result59 := Clob2Array()
	if result59 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result60 := VarcharArray()
	if result60 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result61 := GraphicArray()
	if result61 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result62 := VargraphicArray()
	if result62 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result63 := Float2Array()
	if result63 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result64 := DoubleArray()
	if result64 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result65 := RealArray()
	if result65 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result66 := DecimalArray()
	if result66 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result67 := Decimal2Array()
	if result67 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result68 := Dec2Array()
	if result68 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}

	result69 := DoubleprecisionArray()
	if result69 == nil {
		fmt.Println("Pass")
	} else {
		fmt.Println("Fail")
	}
}

func AllDataTypes() error {
        db, _ := sql.Open("go_ibm_db", connStr)
        defer db.Close()
        db.Exec("Drop table arr")
	 _, err := db.Exec("create table arr (c1 int, c2 SMALLINT, c3 BIGINT, c4 INTEGER, c5 DECIMAL(4,2), c6 NUMERIC, c7 float, c8 double, c9 decfloat, c10 char(10), c11 varchar(10), c12 char for bit data, c13 clob(10),c14 dbclob(100), c15 date, c16 time, c17 timestamp, c18 blob(10), c19 boolean) ccsid unicode")
        if err != nil {
                return err
        }
	st, err := db.Prepare("insert into arr values (1, 2, 9007199254741997, 1234, 67.98, 5689, 56.2390, 34567890, 45.234, 'Vijay', 'Raj', '\x50', 'test123456','▒~P~@▒~P~A▒~P~B▒~P~C▒~P~D▒~P~E▒~P~F','2015-09-10', '10:16:33', '2015-09-10 10:16:33.770139', BLOB(x'616263'), true)");
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query()
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []string")
                return err
        }

	return nil
}

func MultipleQuery() error {
        db, _ := sql.Open("go_ibm_db", connStr)
        defer db.Close()
        db.Exec("Drop table arr")
        _, err := db.Exec("create table arr(PID bigint, C1 varchar(255), C2 varchar(255), C3 varchar(255))")
        if err != nil {
                return err
        }
        _, err = db.Query("Insert into arr values('1', 'PersonA', 'LastNameA', 'QA')")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []string")
                return err
        }

        _, err = db.Query("Insert into arr values('2', 'PersonB', 'LastNameB', 'Dev')")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []string")
                return err
        }

        _, err = db.Query("Insert into arr values('3', 'PersonC', 'LastNameC', 'Dev')")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []string")
                return err
        }

        _, err = db.Query("Insert into arr values('4', 'PersonD', 'LastNameD', 'Dev')")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []string")
                return err
        }

        _, err = db.Query("Insert into arr values('5', 'PersonE', 'LastNameE', 'Dev')")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []string")
                return err
        }

        _, err = db.Query("UPDATE arr SET C3 = 'QA Intern' where C2 = 'LastNameD'")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while updating []string")
                return err
        } else {
                fmt.Println("Update statement successful")
        }

        _, err = db.Query("SELECT count(*) from arr where PID = 7")
        if err != nil {
                return err
        } else {
                fmt.Println("Select statement successful")
        }

        _, err = db.Query("SELECT * from arr where C3 = 'QA Intern'")
        if err != nil {
                return err
        } else {
                fmt.Println("Select statement successful")
        }


        _, err = db.Query("DELETE from arr where PID = 5")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while deleting []string")
                return err
        } else {
                fmt.Println("Delete statement successful")
        }

        _, err = db.Query("INSERT into arr values('6', 'PersonF', 'LastNameF', 'QA Lead')")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []string")
                return err
        }else {
                fmt.Println("Insert statement successful")
        }

        return nil
}

func BadConnectionString() int {
        var errStr string
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

func CreateDropTable() int {

        db, err1 := sql.Open("go_ibm_db", connStr)
        if err1 != nil {
                return 0
        }

        _, err2 := db.Exec("DROP table VMSAMPLE")
        if err2 != nil {
               _, err3 := db.Exec("create table VMSAMPLE(ID varchar(20),NAME varchar(20),LOCATION varchar(20),POSITION varchar(20))")
               if err3 != nil {
                         return 0
               }
        } else {
              _, err4 := db.Exec("create table VMSAMPLE(ID varchar(20),NAME varchar(20),LOCATION varchar(20),POSITION varchar(20))")
              if err4 != nil {
			return 0
             }
       }
       fmt.Println("TABLE CREATED Successfully")

       st, err5 := db.Prepare("Insert into VMSAMPLE(ID,NAME,LOCATION,POSITION) values('3242','mike','hyd','manager')")
       if err5 != nil {
           return 0
       }
       st.Query()

        _, err6 := db.Exec("DROP table VMSAMPLE")
        if err6 != nil {
                return 0
        }
        fmt.Println("TABLE DROP Successfully")

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

func QueryCreateDropTable() int {
        type Db *sql.DB
        var re Db
        re = Createconnection()
        if re == nil {
                return 0
        }

        err := QueryCreateTable(re)
        if err != nil {
                return 0
        }

        err = QueryDropTable(re)
        if err != nil {
                return 0
        }

        return 1
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
//              fmt.Printf("%v  %v   %v         %v\n", t, x, m, n)
        }
        return nil
}


func QueryInsertSelect() int {

        db, _ := sql.Open("go_ibm_db", connStr)
        err := QueryCreateTable(db)
        if err != nil {
                return 0
        }

        err = QueryInsertRow(db)
        if err != nil {
              return 0
        }

        err = QueryDisplayTable(db)
        if err != nil {
               return 0
        }

        err = QueryDropTable(db)
        if err != nil {
                return 0
        }

        return 1
}

func ConnectionInvalidUserPassword() int {
        var errStr string
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

func HugeQuery() int {
        var insertCount string  = "1"
        var tableOne string= "goleaktable1"
        var tableTwo string= "goleaktable2"

        var maxVarChar string = "LEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTES"

        db, _ := sql.Open("go_ibm_db", connStr )

        db.Query("DROP table " + tableOne)
        db.Query("DROP table " + tableTwo)

        _, err := db.Query("CREATE table " + tableOne + "(PID VARCHAR(10), C1 VARCHAR(255), C2 VARCHAR(255), C3 VARCHAR(255))")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                return 0
        }

        _, err = db.Query("CREATE table " + tableTwo + "(PID VARCHAR(10), C1 VARCHAR(255), C2 VARCHAR(255), C3 VARCHAR(255))")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                return 0
        }
        query := "values('" + insertCount + "', '" + maxVarChar + "', '" + maxVarChar + "', '" + maxVarChar + "')"

        for i := 1; i <= 5; i++ {
                 _, err = db.Query("INSERT into " + tableOne + "(PID, C1, C2, C3) " + query)
                if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                        return 0
                }

                _, err = db.Query("INSERT into " + tableTwo + "(PID, C1, C2, C3) " + query)
                if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                        return 0
                }
        }

        rows, err2 := db.Query("SELECT * from " + tableOne)
        if err2 != nil {
                return 0
        }

        defer rows.Close()
        for rows.Next() {
              var t, x, m, n string
              err = rows.Scan(&t, &x, &m, &n)
              if err != nil {
                return 0
              }
              //fmt.Printf("%v  %v   %v  %v\n", t, x, m, n)
        }

        db.Query("DROP table " + tableOne)
        db.Query("DROP table " + tableTwo)

        return 1
}

func DecimalColumn() int {
        var tableOne string= "godecmaltable"

        db, _ := sql.Open("go_ibm_db", connStr )

        db.Query("DROP table " + tableOne)

        _, err := db.Query("CREATE table " + tableOne + "(col1 DECIMAL(30, 2))")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                return 0
        }

        _, err = db.Query("INSERT into " + tableOne + "(col1) values(9999999999999999999999999999.99)")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                 return 0
        }

        _, err = db.Query("INSERT into " + tableOne + "(col1) values(99999999999999999999)")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                 return 0
        }

        rows, err2 := db.Query("SELECT * from " + tableOne)
        if err2 != nil {
                return 0
        }

        defer rows.Close()
        for rows.Next() {
              var f  string
              err = rows.Scan(&f)
              if err != nil {
                return 0
              }
              //fmt.Printf("%v \n", f)
        }

        db.Query("DROP table " + tableOne)

        return 1
}

func RollbackTransaction() error {
	db, err := sql.Open("go_ibm_db", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	bg, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = bg.Exec("CREATE table gorollback(C1 int, C2 float, C3 double, C4 char, C5 varchar(30))")
	if err != nil {
		return err
	}

	err = bg.Rollback()
	if err != nil {
		return err
	}

	return nil
}

func Decfloat() error {
	var tableOne string= "goarr"

	db, err := sql.Open("go_ibm_db", connStr)
        if err != nil {
                return err
        }
	defer db.Close()

	db.Query("DROP table " + tableOne)

	_, err = db.Exec("CREATE table " + tableOne + "(col1 int, col2 decfloat)")
        if err != nil {
                return err
        }

	st1, err1 := db.Prepare("INSERT into " + tableOne + " values(1, 45.678)")
        defer st1.Close()
        if err1 != nil {
                return err1
        }
        _, err1 = st1.Query()
        if !strings.Contains(fmt.Sprint(err1), "did not create a result set") {
                return err1
        }

	st2, err2 := db.Prepare("INSERT into " + tableOne + " values(2, 0.2345600)")
        defer st2.Close()
        if err2 != nil {
                return err2
        }
        _, err2 = st2.Query()
        if !strings.Contains(fmt.Sprint(err2), "did not create a result set") {
                return err2
        }

	st3, err3 := db.Prepare("INSERT into " + tableOne + " values(3, 111e99)")
        defer st3.Close()
        if err3 != nil {
                return err3
        }
        _, err3 = st3.Query()
        if !strings.Contains(fmt.Sprint(err3), "did not create a result set") {
                return err3
        }


	st4, err4 := db.Prepare("INSERT into " + tableOne + " values(4, 111e-99)")
        defer st4.Close()
        if err4 != nil {
                return err4
        }
        _, err4 = st4.Query()
        if !strings.Contains(fmt.Sprint(err4), "did not create a result set") {
                return err4
        }

	st5, err5 := db.Prepare("INSERT into " + tableOne + " values(5, 100.2001234)")
        defer st5.Close()
        if err5 != nil {
                return err5
        }
        _, err5 = st5.Query()
        if !strings.Contains(fmt.Sprint(err5), "did not create a result set") {
                return err5
        }

	st6, err6 := db.Prepare("INSERT into " + tableOne + " values(6, -1000)")
        defer st6.Close()
        if err6 != nil {
                return err6
        }
        _, err6 = st6.Query()
        if !strings.Contains(fmt.Sprint(err6), "did not create a result set") {
                return err6
        }

	st7, err7 := db.Prepare("INSERT into " + tableOne + " values(7, -Inf)")
        defer st7.Close()
        if err7 != nil {
                return err7
        }
        _, err7 = st7.Query()
        if !strings.Contains(fmt.Sprint(err7), "did not create a result set") {
                return err7
        }

        rows, err8 := db.Query("SELECT * from " + tableOne)
        if err8 != nil {
                return err8
        }

        defer rows.Close()
        for rows.Next() {
              var c1, c2  string
	      err9 := rows.Scan(&c1, &c2)
              if err9 != nil {
                      return err9
              }

              //fmt.Printf("%v  %v \n", c1, c2)
        }

        db.Query("DROP table " + tableOne)
        return nil
}






