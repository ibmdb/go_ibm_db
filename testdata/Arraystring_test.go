package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestStringArray(t *testing.T) {
	if StringArray() != nil {
		t.Error("Error at StringArray")
	}
}


//StringArray function performs inserting string array.
func StringArray() error {
        db := Createconnection()
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

