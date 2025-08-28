package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestGraphicArray(t *testing.T) {
	if GraphicArray_1() != nil {
		t.Error("Error at GraphicArray")
	}
}

// GRAPHIC(n) Fixed-length graphic strings that contain n double-byte characters. n must be greater than 0 and less than 128. The default length is 1.
func GraphicArray_1() error {
	var tableOne string = "goarr"

	db := Createconnection()
	defer db.Close()

	db.Query("DROP table " + tableOne)

	_, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 graphic(5))")
	if err != nil {
		fmt.Println("Exec error: ", err)
		return err
	}

	a := []int{1, 2, 3, 4, 5}
	b := []string{"a", "ab", "abc", "abcd", "abcde"}
	st, err := db.Prepare("Insert into " + tableOne + " values(?, ?)")
	if err != nil {
		fmt.Println("Prepare error: ", err)
		return err
	}
	defer st.Close()
	_, err = st.Query(a, b)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []character")
		return err
	}

	var errStr string
	substring := "SQLSTATE=22001"
	c := []int{6}
	d := []string{"abcdef"}
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
