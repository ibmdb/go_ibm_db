package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestDec2Array(t *testing.T) {
	if Dec2Array_1() != nil {
		t.Error("Error at Dec2Array")
	}
}

//dec(p) The decimal data type is an exact numeric data type defined by its precision (total number of digits)
//p Defines the precision. It has at a total of p (<=32) significat digits

func Dec2Array_1() error {
	var tableOne string = "goarr"

	db := Createconnection()
	defer db.Close()

	db.Query("DROP table " + tableOne)

	_, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 dec(31))")
	if err != nil {
		fmt.Println("Exec error: ", err)
		return err
	}

	a := []int{1, 2, 3}
	b := []float32{-10e30, 987654321.123456, 10e30}
	st, err := db.Prepare("Insert into " + tableOne + " values(?, ?)")
	if err != nil {
		fmt.Println("Prepare error: ", err)
		return err
	}
	defer st.Close()
	_, err = st.Query(a, b)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []dec")
		return err
	}

	var errStr string
	substring := "SQLSTATE=22003"
	c := []int{4}
	d := []float32{10e31}
	st, err = db.Prepare("Insert into " + tableOne + " values(?, ?)")
	if err != nil {
		fmt.Println("Prepare error: ", err)
		return err
	}
	defer st.Close()
	_, err = st.Query(c, d)
	if err != nil {
		errStr = fmt.Sprintf("%s", err)

		if !strings.Contains(errStr, substring) {
			fmt.Println("Query error: ", err)
			return err
		}

	}

	rows, err2 := db.Query("SELECT * from " + tableOne)
	if err2 != nil {
		fmt.Println("Query error: ", err2)
		return err2
	}

	defer rows.Close()
	for rows.Next() {
		var c1, c2 string
		err = rows.Scan(&c1, &c2)
		if err != nil {
			fmt.Println("Scan error: ", err)
			return err
		}

		//fmt.Printf("%v  %v \n", c1, c2)
	}

	db.Query("DROP table " + tableOne)

	return nil
}
