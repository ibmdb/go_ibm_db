package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestFloat2Array(t *testing.T) {
	if Float2Array_1() != nil {
		t.Error("Error at Float2Array")
	}
}

func Float2Array_1() error {
	var tableOne string = "goarr"

	db := Createconnection()
	defer db.Close()

	db.Query("DROP table " + tableOne)

	_, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 float)")
	if err != nil {
		fmt.Println("Exec error: ", err)
		return err
	}

	a := []int{1, 2, 3, 4, 5}
	b := []float64{-1.79769e308, 1.79769e308, 9876543210.123456789, 2.225e-307, -2.225e-307}
	st, err := db.Prepare("Insert into " + tableOne + " values(?, ?)")
	if err != nil {
		fmt.Println("Prepare error: ", err)
		return err
	}
	defer st.Close()
	_, err = st.Query(a, b)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []float")
		return err
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
