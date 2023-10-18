package main

import (
        "fmt"
        "strings"
        "testing"
)

func TestDecimal2Array(t *testing.T) {
	if Decimal2Array_1() != nil {
		t.Error("Error at Decimal2Array")
	}
}

//decimal(p) The decimal data type is an exact numeric data type defined by its precision (total number of digits)
//p Defines the precision. It has at a total of p (<=32) significat digits

func Decimal2Array_1() error {
        var tableOne string= "goarr"

        db := Createconnection()
        defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 numeric(31))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3}
        b :=  []float32 { -10e30, 987654321.123456, 10e30}
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []decimal")
                return err
        }

        var errStr string
        substring := "SQLSTATE=22003"
        c :=  []int{4}
        d :=  []float32{10e31}
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

