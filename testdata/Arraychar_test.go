package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestCharArray(t *testing.T) {
	if CharArray() != nil {
		t.Error("Error at CharArray")
	}
}

// CharArray function performs inserting float32,float64 datatypes.
func CharArray() error {
	db := Createconnection()
	defer db.Close()

	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(var1 character, var2 character)")
	if err != nil {
		fmt.Println("Exec error: ", err)
		return err
	}
	a := []string{"a", "b", "c"}
	b := []string{"z", "y", "x"}
	st, err := db.Prepare("Insert into arr values(?,?)")
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
	return nil
}
