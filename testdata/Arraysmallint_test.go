package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestSmallintArray(t *testing.T) {
	if SmallintArray() != nil {
		t.Error("Error at SmallintArray")
	}
}

// A small integer is a binary integer with a precision of 15 bits. The range of small integers is -32768 to +32767.
func SmallintArray() error {
	var tableOne string = "goarr"
	var errStr string

	db := Createconnection()
	defer db.Close()

	db.Query("DROP table " + tableOne)

	_, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 smallint)")
	if err != nil {
		fmt.Println("Exec error: ", err)
		return err
	}

	a := []int{1, 2, 3, 4, 5}
	b := []int{-32768, -123, 0, 100, 32767}
	st, err := db.Prepare("Insert into " + tableOne + " values(?, ?)")
	if err != nil {
		fmt.Println("Prepare error: ", err)
		return err
	}
	defer st.Close()
	_, err = st.Query(a, b)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []smallint")
		return err
	}

	substring := "SQLSTATE=22003"
	c := []int{6, 7}
	d := []int{-32769, 32768}
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
