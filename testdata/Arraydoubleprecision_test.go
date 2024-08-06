package main

import (
        "fmt"
        "strings"
        "testing"
)


func TestDoublePrecisionArray(t *testing.T) {
	if DoublePrecisionArray() != nil {
		t.Error("Error at DoublePrecisionArray")
	}
}

//DOUBLE PRECISION value ranges:
//Smallest DOUBLE value: -1.79769E+308
//Largest DOUBLE value: 1.79769E+308
//Smallest positive DOUBLE value: 2.225E-307
//Largest negative DOUBLE value: -2.225E-307
func DoublePrecisionArray() error {
        var tableOne string= "goarr"

        db := Createconnection()
        defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 double precision)")
        if err != nil {
                fmt.Println("Exec error: ", err)
                return err
        }

        a :=  []int{1, 2, 3, 4}
        b :=  []float64 { -1.79769E+308, 1.79769E+308, 2.225E-307,  -2.225E-307 }
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                fmt.Println("Prepare error: ", err)
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []double precision")
                return err
        }
        rows, err2 := db.Query("SELECT * from " + tableOne)
        if err2 != nil {
                fmt.Println("Query error: ", err2)
                return err2
        }

        defer rows.Close()
        for rows.Next() {
              var c1, c2  string
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

