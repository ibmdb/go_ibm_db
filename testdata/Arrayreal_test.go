package main

import (
        "fmt"
        "strings"
        "testing"
)

func TestRealArray(t *testing.T) {
	if RealArray_1() != nil {
		t.Error("Error at RealArray")
	}
}



//REAL value ranges:
//Smallest REAL value: -3.402E+38
//Largest REAL value: 3.402E+38
//Smallest positive REAL value: 1.175E-37
//Largest negative REAL value: -1.175E-37
func RealArray_1() error {
        var tableOne string= "goarr"

        db := Createconnection()
        defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 real)")
        if err != nil {
                fmt.Println("Exec error: ", err)
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []float32 { -3.402e38, 987654321.123456, 3.402e38, 1.175e-37, -1.175e-37}
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                fmt.Println("Prepare error: ", err)
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []real")
                return err
        }
        rows, err2 := db.Query("SELECT * from " + tableOne)
        if err2 != nil {
                fmt.Println("Query error: ", err)
                return err
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

