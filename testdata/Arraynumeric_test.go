package main

import (
        "fmt"
        "strings"
        "testing"
)


func TestNumericArray(t *testing.T) {
	if NumericArray_1() != nil {
		t.Error("Error at NumericArray")
	}
}


//numeric(p,s) The decimal data type is an exact numeric data type defined by its precision (total number of digits)
//and scale (number of digits to the right of the decimal point).
//p Defines the precision. Minimum 1; maximum is 39
//s Defines the scale. The scale of a decimal value cannot exceed its precision. Scale can be 0 (no digits to the right of the decimal point).

func NumericArray_1() error {
        var tableOne string= "goarr"

        db := Createconnection()

        db.Query("DROP table " + tableOne)
        defer db.Close()

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 numeric(5,2))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3}
        b :=  []float32 { -152.3, 56.08, 100.238567}
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []numeric")
                return err
        }

        var errStr string
        substring := "SQLSTATE=22003"
        c :=  []int{4}
        d :=  []float32{1234.98}
         st, err = db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(c, d)
        if err != nil {
                errStr = fmt.Sprintf("%s", err)

                if !strings.Contains(errStr, substring) {
                        return err
                }
        }
        rows, err2 := db.Query("SELECT * from " + tableOne)
        if err2 != nil {
                return err
        }

        defer rows.Close()
        for rows.Next() {
              var c1, c2  string
              err = rows.Scan(&c1, &c2)
              if err != nil {
                      return err
              }

              //fmt.Printf("%v  %v \n", c1, c2)
        }

        db.Query("DROP table " + tableOne)

        return nil
}

