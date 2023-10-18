package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestBoolArray(t *testing.T) {
	if BoolArray() != nil {
		t.Error("Error at BoolArray")
	}
}


//BoolArray function performs inserting bool array.
func BoolArray() error {
        db := Createconnection()
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

