package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestDecimalArray(t *testing.T) {
	if DecimalArray_1() != nil {
		t.Error("Error at DecimalArray")
	}
}

// decimal(p,s) The decimal data type is an exact numeric data type defined by its precision (total number of digits)
// and scale (number of digits to the right of the decimal point).
// p Defines the precision. Minimum 1; maximum is 39
// s Defines the scale. The scale of a decimal value cannot exceed its precision. Scale can be 0 (no digits to the right of the decimal point).
func DecimalArray_1() error {
	var tableOne string = "goarr"

	db := Createconnection()
	defer db.Close()

	db.Query("DROP table " + tableOne)

	_, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 decimal(5,2))")
	if err != nil {
		fmt.Println("Exec error: ", err)
		return err
	}

	a := []int{1, 2, 3}
	b := []float32{-152.3, 56.08, 100.238567}
	st, err := db.Prepare("Insert into " + tableOne + " values(?, ?)")
	if err != nil {
		fmt.Println("Prepare error: ", err)
		return err
	}
	defer st.Close()
	_, err = st.Query(a, b)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []decimal")
		return err
	}

	var errStr string
	substring := "SQLSTATE=22003"
	c := []int{4}
	d := []float32{1234.98}
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
