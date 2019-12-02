package main

import (
	"database/sql"
	"fmt"
	"time"

	a "github.com/ibmdb/go_ibm_db"
)

var con = "PROTOCOL=tcpip;HOSTNAME=localhost;PORT=50000;DATABASE=go;UID=uname;PWD=pwd"
var conDB = "PROTOCOL=tcpip;HOSTNAME=localhost;PORT=50000;UID=uname;PWD=pwd"

//Createconnection will return the db instance
func Createconnection() (db *sql.DB) {
	db, _ = sql.Open("go_ibm_db", con)
	return db
}

//Createtable will create the tables
func Createtable() error {
	db, err := sql.Open("go_ibm_db", con)
	db.Exec("DROP table rocket")
	_, err = db.Exec("create table rocket(a int)")
	_, err = db.Exec("create table rocket1(a int)")
	if err != nil {
		return err
	}
	return nil
}

//Insert will insert data in to the table
func Insert() error {
	db, err := sql.Open("go_ibm_db", con)
	_, err = db.Exec("insert into rocket values(1)")
	if err != nil {
		return err
	}
	return nil
}

//Drop will drop the table
func Drop() error {
	db, err := sql.Open("go_ibm_db", con)
	_, err = db.Exec("drop table rocket1")
	if err != nil {
		return err
	}
	return nil
}

//Prepare will prepare the statement
func Prepare() error {
	db, _ := sql.Open("go_ibm_db", con)
	_, err := db.Prepare("select * from rocket")
	if err != nil {
		return err
	}
	return nil
}

//Query will execute the prepared statement
func Query() error {
	db, _ := sql.Open("go_ibm_db", con)
	st, _ := db.Prepare("select * from rocket")
	_, err := st.Query()
	if err != nil {
		return err
	}
	return nil
}

//Scan will Scan the data in the rows
func Scan() error {
	db, _ := sql.Open("go_ibm_db", con)
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
	db, _ := sql.Open("go_ibm_db", con)
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
	db, _ := sql.Open("go_ibm_db", con)
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
	db, err := sql.Open("go_ibm_db", con)
	err = db.QueryRow("select a from rocket where a=?", a).Scan(&uname)
	if err != nil {
		return err
	}
	return nil
}

//Begin will start a transaction
func Begin() error {
	db, err := sql.Open("go_ibm_db", con)
	_, err = db.Begin()
	if err != nil {
		return err
	}
	return nil
}

//Commit will commit the uncommited transactions
func Commit() error {
	db, err := sql.Open("go_ibm_db", con)
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
	db, _ := sql.Open("go_ibm_db", con)
	err := db.Close()
	if err != nil {
		return err
	}
	return nil
}

//PoolOpen creates a pool and makes a connection.
func PoolOpen() int {
	pool := a.Pconnect("PoolSize=50")
	db := pool.Open(con)
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
	db, _ := sql.Open("go_ibm_db", con)
	_, err := db.Exec("call sysproc.get_dbsize_info(?, ?, ?,0)", sql.Out{Dest: &snapTime}, sql.Out{Dest: &dbsize}, sql.Out{Dest: &dbcapacity})
	if err != nil {
		return err
	}
	return nil
}

//StoredProcedureInOut function tests OUT Parameter by calling get_dbsize_info.
func StoredProcedureInOut() error {
	db, _ := sql.Open("go_ibm_db", con)
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
	db, _ := sql.Open("go_ibm_db", con)
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
	if err.Error() != "Stmt did not create a result set" {
		fmt.Println("Error while inserting []int")
		return err
	}
	_, err = st.Query(b)
	if err.Error() != "Stmt did not create a result set" {
		fmt.Println("Error while inserting []int8")
		return err
	}
	_, err = st.Query(c)
	if err.Error() != "Stmt did not create a result set" {
		fmt.Println("Error while inserting []int16")
		return err
	}
	_, err = st.Query(d)
	if err.Error() != "Stmt did not create a result set" {
		fmt.Println("Error while inserting []int32")
		return err
	}
	_, err = st.Query(e)
	if err.Error() != "Stmt did not create a result set" {
		fmt.Println("Error while inserting []int64")
		return err
	}
	return nil
}

//StringArray function performs inserting string array.
func StringArray() error {
	db, _ := sql.Open("go_ibm_db", con)
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
	if err.Error() != "Stmt did not create a result set" {
		fmt.Println("Error while inserting []string")
		return err
	}
	return nil
}

//BoolArray function performs inserting bool array.
func BoolArray() error {
	db, _ := sql.Open("go_ibm_db", con)
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
	if err.Error() != "Stmt did not create a result set" {
		fmt.Println("Error while inserting []bool")
		return err
	}
	return nil
}

//FloatArray function performs inserting float32,float64 datatypes.
func FloatArray() error {
	db, _ := sql.Open("go_ibm_db", con)
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
	if err.Error() != "Stmt did not create a result set" {
		fmt.Println("Error while inserting []float32")
		return err
	}
	_, err = st.Query(b)
	if err.Error() != "Stmt did not create a result set" {
		fmt.Println("Error while inserting []float64")
		return err
	}
	return nil
}

//CharArray function performs inserting float32,float64 datatypes.
func CharArray() error {
	db, _ := sql.Open("go_ibm_db", con)
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
	if err.Error() != "Stmt did not create a result set" {
		fmt.Println("Error while inserting []character")
		return err
	}
	return nil
}

//TimeStampArray function performs inserting float32,float64 datatypes.
func TimeStampArray() error {
	db, _ := sql.Open("go_ibm_db", con)
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
	if err.Error() != "Stmt did not create a result set" {
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
	db, _ := sql.Open("go_ibm_db", con)
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
	if err.Error() != "Stmt did not create a result set" {
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
	db, _ := sql.Open("go_ibm_db", con)
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
	if err.Error() != "Stmt did not create a result set" {
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
	db, _ := sql.Open("go_ibm_db", con)
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
	if err.Error() != "Stmt did not create a result set" {
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
	db, _ := sql.Open("go_ibm_db", con)
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
	if err.Error() != "Stmt did not create a result set" {
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
	db, _ := sql.Open("go_ibm_db", con)
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
	if err.Error() != "Stmt did not create a result set" {
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
	db, _ := sql.Open("go_ibm_db", con)
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
	if err.Error() != "Stmt did not create a result set" {
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
	res, err := a.CreateDb("Goo", conDB)
	if err != nil {
		return false
	}
	return res
}

//DropDB will drop database
func DropDB() bool {
	res, err := a.DropDb("Goo", conDB)
	if err != nil {
		return false
	}
	return res
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
}
