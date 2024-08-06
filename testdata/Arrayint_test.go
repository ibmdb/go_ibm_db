package main

import (
	"fmt"
	"strings"
        "testing"
)

func TestIntArray(t *testing.T) {
	if IntArray() != nil {
		t.Error("Error at IntArray")
	}
}

//IntArray function performs inserting int,int8,int16,int32,int64 datatypes.
func IntArray() error {
        db := Createconnection()
        defer db.Close()

        db.Exec("Drop table arr")
        _, err := db.Exec("create table arr(var1 int)")
        if err != nil {
                fmt.Println("Exec error: ", err)
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
                fmt.Println("Prepare error: ", err)
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

